package scheduled

import (
	"fmt"

	sdkstatus "github.com/nickschuch/operator-sdk-status"
	sdkaction "github.com/operator-framework/operator-sdk/pkg/sdk/action"
	sdkhandler "github.com/operator-framework/operator-sdk/pkg/sdk/handler"
	sdktypes "github.com/operator-framework/operator-sdk/pkg/sdk/types"
	"github.com/pkg/errors"
	"github.com/prometheus/common/log"
	batchv1beta1 "k8s.io/api/batch/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"

	"github.com/previousnext/mysql-toolkit/internal/operator/apis/mtk/v1alpha1"
	"github.com/previousnext/mysql-toolkit/internal/operator/handler/acquia/generate"
	"github.com/previousnext/mysql-toolkit/internal/operator/handler/acquia/secrets"
)

const (
	prefix          = "mtk-acquia-scheduled"
	errUpdateStatus = "failed to update status"
)

// NewHandler for responding to Acquia Scheduled events.
func NewHandler(namespace, secret, image, cpu, memory string) sdkhandler.Handler {
	return &Handler{
		Namespace: namespace,
		Secret:    secret,
		Image:     image,
		CPU:       cpu,
		Memory:    memory,
	}
}

// Handler for responding to Acquia Scheduled events.
type Handler struct {
	Namespace string
	Secret    string
	Image     string
	CPU       string
	Memory    string
}

// Handle Acquia Scheduled events.
func (h *Handler) Handle(ctx sdktypes.Context, event sdktypes.Event) error {
	switch cr := event.Object.(type) {
	case *v1alpha1.AcquiaSnapshotScheduled:
		err := reconcile(h.Namespace, h.Secret, h.Image, h.CPU, h.Memory, cr)
		if err != nil {
			return err
		}
	}

	return nil
}

// Helper function to reconcile status updates and execution.
func reconcile(namespace, secret, image, cpu, memory string, cr *v1alpha1.AcquiaSnapshotScheduled) error {
	var p sdkstatus.Pipeline

	// Common identifier for all Job/ConfigMap/Secret objects within the operators namespace.
	name := fmt.Sprintf("%s-%s-%s", prefix, cr.ObjectMeta.Namespace, cr.ObjectMeta.Name)

	ref := []metav1.OwnerReference{
		*metav1.NewControllerRef(cr, schema.GroupVersionKind{
			Group:   v1alpha1.SchemeGroupVersion.Group,
			Version: v1alpha1.SchemeGroupVersion.Version,
			Kind:    cr.Kind,
		}),
	}

	logger := log.With("namespace", cr.ObjectMeta.Namespace).With("name", cr.ObjectMeta.Name)

	logger.Infoln("Starting reconciliation loop")

	p.Add("Create ConfigMap", func() (sdkstatus.Status, error) {
		logger.Infoln("Generating ConfigMap")

		obj, err := generate.ConfigMap(namespace, name, cr.Spec.Config)
		if err != nil {
			return sdkstatus.StatusFailed, errors.Wrap(err, "failed to generate: ConfigMap")
		}

		obj.ObjectMeta.OwnerReferences = ref

		logger.Infoln("Creating ConfigMap")

		err = sdkaction.Create(obj)
		if err != nil {
			return sdkstatus.StatusFailed, errors.Wrap(err, "failed to create: ConfigMap")
		}

		return sdkstatus.StatusFinished, nil
	})

	p.Add("Create Secret", func() (sdkstatus.Status, error) {
		logger.Infoln("Loading Secrets")

		values, err := secrets.Load(namespace, secret, cr.ObjectMeta.Namespace, cr.Spec.Credentials)
		if err != nil {
			return sdkstatus.StatusFailed, errors.Wrap(err, "failed to get values: Secret")
		}

		logger.Infoln("Generating Secret")

		obj, err := generate.Secret(namespace, name, values)
		if err != nil {
			return sdkstatus.StatusFailed, errors.Wrap(err, "failed to generate: Secret")
		}

		obj.ObjectMeta.OwnerReferences = ref

		logger.Infoln("Creating Secret")

		err = sdkaction.Create(obj)
		if err != nil {
			return sdkstatus.StatusFailed, errors.Wrap(err, "failed to create: Secret")
		}

		return sdkstatus.StatusFinished, nil
	})

	p.Add("Create CronJob", func() (sdkstatus.Status, error) {
		logger.Infoln("Generating CronJob")

		obj, err := generate.Job(namespace, name, image, cpu, memory, cr.Spec.Database, cr.Spec.Docker)
		if err != nil {
			return sdkstatus.StatusFailed, errors.Wrap(err, "failed to generate: Job")
		}

		var (
			deadline  int64 = 800
			successes int32 = 5
			failures  int32 = 5
		)

		cronjob := &batchv1beta1.CronJob{
			TypeMeta: metav1.TypeMeta{
				Kind:       "CronJob",
				APIVersion: "batch/v1beta1",
			},
			ObjectMeta: metav1.ObjectMeta{
				Namespace:       namespace,
				Name:            name,
				OwnerReferences: ref,
			},
			Spec: batchv1beta1.CronJobSpec{
				Schedule:                   cr.Spec.Schedule,
				StartingDeadlineSeconds:    &deadline,
				ConcurrencyPolicy:          batchv1beta1.ForbidConcurrent,
				SuccessfulJobsHistoryLimit: &successes,
				FailedJobsHistoryLimit:     &failures,
				JobTemplate: batchv1beta1.JobTemplateSpec{
					ObjectMeta: metav1.ObjectMeta{
						Namespace: namespace,
					},
					Spec: obj.Spec,
				},
			},
		}

		logger.Infoln("Creating CronJob")

		err = sdkaction.Create(cronjob)
		if err != nil {
			return sdkstatus.StatusFailed, errors.Wrap(err, "failed to create: CronJob")
		}

		return sdkstatus.StatusFinished, nil
	})

	result, err := p.Run(cr.Status.Steps)

	logger.Infoln("Updating CustomResource status")

	// Save the object back so we can pickup the pipeline from where we left off.
	cr.Status.Steps = result
	err = sdkaction.Update(cr)

	logger.Infoln("Reconciliation loop finished")

	return err
}
