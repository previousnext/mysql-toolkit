package oneshot

import (
	"fmt"
	"os"
	"time"

	sdkaction "github.com/operator-framework/operator-sdk/pkg/sdk/action"
	sdkhandler "github.com/operator-framework/operator-sdk/pkg/sdk/handler"
	sdkquery "github.com/operator-framework/operator-sdk/pkg/sdk/query"
	sdktypes "github.com/operator-framework/operator-sdk/pkg/sdk/types"
	"github.com/pkg/errors"
	"github.com/prometheus/common/log"
	"k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"

	"github.com/previousnext/mysql-toolkit/internal/envar"
	"github.com/previousnext/mysql-toolkit/internal/operator/apis/mtk/v1alpha1"
	"github.com/previousnext/mysql-toolkit/internal/operator/handler/helpers"
)

// NewHandler for responding to Acquia OneShot events.
func NewHandler() sdkhandler.Handler {
	var (
		cpu = os.Getenv(envar.K8sJobCPU)
		mem = os.Getenv(envar.K8sJobMemory)
	)

	return &Handler{
		Image: os.Getenv(envar.K8sJobImage),
		Resources: corev1.ResourceRequirements{
			Requests: corev1.ResourceList{
				corev1.ResourceCPU:    resource.MustParse(cpu),
				corev1.ResourceMemory: resource.MustParse(mem),
			},
			Limits: corev1.ResourceList{
				corev1.ResourceCPU:    resource.MustParse(cpu),
				corev1.ResourceMemory: resource.MustParse(mem),
			},
		},
		Secret: os.Getenv(envar.K8sJobSecret),
	}
}

// Handler for responding to Acquia OneShot events.
type Handler struct {
	Image     string
	Resources corev1.ResourceRequirements
	Secret    string
}

// Handle Acquia OneShot events.
func (h *Handler) Handle(ctx sdktypes.Context, event sdktypes.Event) error {
	switch cr := event.Object.(type) {
	case *v1alpha1.AcquiaOneShot:
		logger := helpers.NewLogger(cr.ObjectMeta.Namespace, cr.ObjectMeta.Name)

		err := reconcile(logger, h.Image, h.Secret, h.Resources, cr)
		if err != nil && !apierrors.IsAlreadyExists(err) {
			cr.Status.Phase = v1alpha1.PhaseFailed
			cr.Status.Error = err.Error()

			err = sdkaction.Update(cr)
			if err != nil {
				logger.Errorln(err)
				return errors.Wrap(err, "failed to save AcquiaOneShot status update")
			}

			logger.Errorln(err)
			return err
		}

		cr.Status.Phase = v1alpha1.PhaseComplete
		err = sdkaction.Update(cr)
		if err != nil {
			logger.Errorln(err)
			return errors.Wrap(err, "failed to save AcquiaOneShot status update")
		}

		logger.Infoln("Finished!")
	}

	return nil
}

// Helper function to reconcile Acquia OneShot events.
func reconcile(logger log.Logger, image, secret string, resources corev1.ResourceRequirements, cr *v1alpha1.AcquiaOneShot) error {
	if cr.Status.Job != "" {
		return apierrors.NewAlreadyExists(schema.GroupResource{}, cr.Status.Job)
	}

	logger.Infoln("Validating object")

	if cr.Spec.Site == "" {
		return errors.New("Not found: spec.Site")
	}

	if cr.Spec.Environment == "" {
		return errors.New("Not found: spec.Environment")
	}

	if cr.Spec.Database == "" {
		return errors.New("Not found: spec.Database")
	}

	if cr.Spec.Credentials.Acquia == "" {
		return errors.New("Not found: spec.Credentials.Acquia")
	}

	if cr.Spec.Credentials.Docker == "" {
		return errors.New("Not found: spec.Credentials.Docker")
	}

	logger.Infoln("Generating Job Spec")

	spec, err := GenerateJobSpec(GenerateJobSpecInput{
		ObjectMeta:  cr.ObjectMeta,
		Site:        cr.Spec.Site,
		Environment: cr.Spec.Environment,
		Database:    cr.Spec.Database,
		Image:       cr.Spec.Image,
		Credentials: GenerateJobSpecInputCredentials{
			Acquia: cr.Spec.Credentials.Acquia,
			Docker: cr.Spec.Credentials.Docker,
			AWS:    secret,
		},
		Job: GenerateJobSpecInputJob{
			Image:     image,
			Resources: resources,
		},
	})
	if err != nil {
		return errors.Wrap(err, "failed to generate Job")
	}

	job := &v1.Job{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Job",
			APIVersion: "batch/v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			GenerateName: fmt.Sprintf("mtk-%s-%s", cr.ObjectMeta.Namespace, cr.ObjectMeta.Name),
			Namespace:    cr.ObjectMeta.Namespace,
			OwnerReferences: []metav1.OwnerReference{
				*metav1.NewControllerRef(cr, schema.GroupVersionKind{
					Group:   v1alpha1.SchemeGroupVersion.Group,
					Version: v1alpha1.SchemeGroupVersion.Version,
					Kind:    cr.Kind,
				}),
			},
		},
		Spec: spec,
	}

	logger.Infoln("Creating Job")

	err = sdkaction.Create(job)
	if err != nil {
		return errors.Wrap(err, "failed to create Job")
	}

	logger.Infoln("Updating Job status")

	cr.Status.Job = job.Name
	cr.Status.Phase = v1alpha1.PhaseRunning
	err = sdkaction.Update(cr)
	if err != nil {
		return errors.Wrap(err, "failed to save AcquiaOneShot status update")
	}

	logger.With("job", job.Name).Infoln("Waiting for Job to finish")

	limiter := time.Tick(time.Second * 5)

Watcher:
	for {
		<-limiter

		err := sdkquery.Get(job)
		if err != nil {
			logger.With("job", job.Name).Errorln(err)
			break Watcher
		}

		// Set status phase to failed if job containers are in failed state.
		if job.Status.Failed > int32(0) {
			logger.With("job", job.Name).Errorln(err)
			break Watcher
		}

		// Set status phase to successful if the job indicates as complete.
		for _, condition := range job.Status.Conditions {
			if condition.Type == "Complete" {
				switch condition.Status {
				case corev1.ConditionTrue:
					break Watcher

				case corev1.ConditionFalse:
					if cr.Status.Phase != v1alpha1.PhaseRunning {
						cr.Status.Phase = v1alpha1.PhaseRunning
						err = sdkaction.Update(cr)
						if err != nil {
							logger.With("job", job.Name).Errorln(err)
							return errors.Wrap(err, "failed to save AcquiaOneShot status update")
						}
					}

					continue Watcher
				}

			}
		}
	}

	return nil
}
