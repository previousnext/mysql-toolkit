package snapshot

import (
	"fmt"
	"time"

	sdkaction "github.com/operator-framework/operator-sdk/pkg/sdk/action"
	sdkhandler "github.com/operator-framework/operator-sdk/pkg/sdk/handler"
	sdkquery "github.com/operator-framework/operator-sdk/pkg/sdk/query"
	sdktypes "github.com/operator-framework/operator-sdk/pkg/sdk/types"
	"github.com/pkg/errors"
	"github.com/prometheus/common/log"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"

	"github.com/previousnext/mysql-toolkit/internal/operator/apis/mtk/v1alpha1"
	"github.com/previousnext/mysql-toolkit/internal/operator/handler/acquia/job"
	"github.com/previousnext/mysql-toolkit/internal/operator/handler/acquia/secrets"
)

const (
	prefix          = "mtk-acquia"
	errUpdateStatus = "failed to update status"
)

// NewHandler for responding to Acquia Snapshot events.
func NewHandler(namespace, secret, image, cpu, memory string) sdkhandler.Handler {
	return &Handler{
		Namespace: namespace,
		Secret:    secret,
		Image:     image,
		CPU:       cpu,
		Memory:    memory,
	}
}

// Handler for responding to Acquia Snapshot events.
type Handler struct {
	Namespace string
	Secret    string
	Image     string
	CPU       string
	Memory    string
}

// Handle Acquia Snapshot events.
func (h *Handler) Handle(ctx sdktypes.Context, event sdktypes.Event) error {
	switch cr := event.Object.(type) {
	case *v1alpha1.AcquiaSnapshot:
		logger := log.With("namespace", cr.ObjectMeta.Namespace).With("name", cr.ObjectMeta.Name)

		err := reconcile(logger, h.Namespace, h.Secret, h.Image, h.CPU, h.Memory, cr)
		if err != nil {
			return err
		}
	}

	return nil
}

// Helper function to reconcile status updates and execution.
func reconcile(logger log.Logger, namespace, secret, image, cpu, memory string, cr *v1alpha1.AcquiaSnapshot) error {
	if cr.Status.Phase != "" {
		return nil
	}

	err := updateStatus(logger, cr, v1alpha1.PhaseRunning, "Loading Configuration Job")
	if err != nil {
		return errors.Wrap(err, errUpdateStatus)
	}

	compiledSecrets, err := secrets.Load(namespace, secret, cr.ObjectMeta.Namespace, cr.Spec.Credentials)
	if err != nil {
		return errors.Wrap(err, "failed to compile configuration")
	}

	err = updateStatus(logger, cr, v1alpha1.PhaseRunning, "Generating Job")
	if err != nil {
		return errors.Wrap(err, errUpdateStatus)
	}

	generateParams := job.Params{
		Namespace: namespace,
		Name:      fmt.Sprintf("%s-%s-%s", prefix, cr.ObjectMeta.Namespace, cr.ObjectMeta.Name),
		Database:  cr.Spec.Database,
		Docker:    cr.Spec.Docker,
		Secrets:   compiledSecrets,
		Config:    cr.Spec.Config,
		Image:     image,
		CPU:       cpu,
		Memory:    memory,
	}

	genJob, genConfigmap, genSecret, err := job.Generate(generateParams)
	if err != nil {
		return errors.Wrap(err, "failed to generate Job")
	}

	ref := []metav1.OwnerReference{
		*metav1.NewControllerRef(cr, schema.GroupVersionKind{
			Group:   v1alpha1.SchemeGroupVersion.Group,
			Version: v1alpha1.SchemeGroupVersion.Version,
			Kind:    cr.Kind,
		}),
	}

	genJob.ObjectMeta.OwnerReferences = ref
	genConfigmap.ObjectMeta.OwnerReferences = ref
	genSecret.ObjectMeta.OwnerReferences = ref

	err = updateStatus(logger, cr, v1alpha1.PhaseRunning, "Creating Job")
	if err != nil {
		return errors.Wrap(err, errUpdateStatus)
	}

	err = sdkaction.Create(genConfigmap)
	if err != nil {
		return errors.Wrap(err, "failed to create ConfigMap")
	}

	err = sdkaction.Create(genSecret)
	if err != nil {
		return errors.Wrap(err, "failed to create Secret")
	}

	err = sdkaction.Create(genJob)
	if err != nil {
		return errors.Wrap(err, "failed to create Job")
	}

	err = updateStatus(logger, cr, v1alpha1.PhaseRunning, "Waiting for Job to finish")
	if err != nil {
		return errors.Wrap(err, errUpdateStatus)
	}

	err = wait(genJob)
	if err != nil {
		return errors.Wrap(err, "job failed")
	}

	err = updateStatus(logger, cr, v1alpha1.PhaseComplete, "Image has been built")
	if err != nil {
		return errors.Wrap(err, errUpdateStatus)
	}

	return nil
}

// Helper function to update the CustomResource status.
func updateStatus(logger log.Logger, cr *v1alpha1.AcquiaSnapshot, status v1alpha1.Phase, message string) error {
	logger.Info(status, message)

	cr.Status = v1alpha1.AcquiaStatus{
		Phase:   status,
		Updated: time.Now(),
		Message: message,
	}

	return sdkaction.Update(cr)
}

// Helper function to wait for Job to finish.
func wait(job *batchv1.Job) error {
	limiter := time.Tick(time.Second * 5)

	for {
		<-limiter

		err := sdkquery.Get(job)
		if err != nil {
			return errors.Wrap(err, "failed to load Job")
		}

		// Set status phase to failed if job containers are in failed state.
		if job.Status.Failed > int32(0) {
			return errors.Wrap(err, "job failed")
		}

		// Set status phase to successful if the job indicates as complete.
		if finished(job) {
			break
		}
	}

	return nil
}

// Helper function to check if a job has finished.
func finished(job *batchv1.Job) bool {
	for _, condition := range job.Status.Conditions {
		if condition.Type == batchv1.JobComplete && condition.Status == corev1.ConditionTrue {
			return true
		}
	}

	return false
}
