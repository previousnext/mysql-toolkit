package scheduled

import (
	"fmt"
	"time"

	sdkaction "github.com/operator-framework/operator-sdk/pkg/sdk/action"
	sdkhandler "github.com/operator-framework/operator-sdk/pkg/sdk/handler"
	sdktypes "github.com/operator-framework/operator-sdk/pkg/sdk/types"
	"github.com/pkg/errors"
	"github.com/prometheus/common/log"
	batchv1beta1 "k8s.io/api/batch/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"

	"github.com/previousnext/mysql-toolkit/internal/operator/apis/mtk/v1alpha1"
	"github.com/previousnext/mysql-toolkit/internal/operator/handler/acquia/job"
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
		logger := log.With("namespace", cr.ObjectMeta.Namespace).With("name", cr.ObjectMeta.Name)

		err := reconcile(logger, h.Namespace, h.Secret, h.Image, h.CPU, h.Memory, cr)
		if err != nil {
			logger.Errorln(err)
			return err
		}
	}

	return nil
}

// Helper function to reconcile status updates and execution.
func reconcile(logger log.Logger, namespace, secret, image, cpu, memory string, cr *v1alpha1.AcquiaSnapshotScheduled) error {
	if cr.Status.Phase != "" {
		return nil
	}

	err := updateStatus(logger, cr, v1alpha1.PhaseRunning, "Loading Configuration")
	if err != nil {
		return errors.Wrap(err, errUpdateStatus)
	}

	compiledSecrets, err := secrets.Load(namespace, secret, cr.ObjectMeta.Namespace, cr.Spec.Credentials)
	if err != nil {
		return errors.Wrap(err, "failed to compile configuration")
	}

	err = updateStatus(logger, cr, v1alpha1.PhaseRunning, "Generating CronJob")
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

	err = updateStatus(logger, cr, v1alpha1.PhaseRunning, "Creating CronJob")
	if err != nil {
		return errors.Wrap(err, errUpdateStatus)
	}

	ref := []metav1.OwnerReference{
		*metav1.NewControllerRef(cr, schema.GroupVersionKind{
			Group:   v1alpha1.SchemeGroupVersion.Group,
			Version: v1alpha1.SchemeGroupVersion.Version,
			Kind:    cr.Kind,
		}),
	}

	genConfigmap.ObjectMeta.OwnerReferences = ref
	genSecret.ObjectMeta.OwnerReferences = ref

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
			Namespace: namespace,
			Name:      fmt.Sprintf("%s-%s-%s", prefix, cr.ObjectMeta.Namespace, cr.ObjectMeta.Name),
			OwnerReferences: []metav1.OwnerReference{
				*metav1.NewControllerRef(cr, schema.GroupVersionKind{
					Group:   v1alpha1.SchemeGroupVersion.Group,
					Version: v1alpha1.SchemeGroupVersion.Version,
					Kind:    cr.Kind,
				}),
			},
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
				Spec: genJob.Spec,
			},
		},
	}

	err = updateStatus(logger, cr, v1alpha1.PhaseRunning, "Creating CronJob")
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

	err = sdkaction.Create(cronjob)
	if err != nil {
		return errors.Wrap(err, "failed to create CronJob")
	}

	err = updateStatus(logger, cr, v1alpha1.PhaseComplete, fmt.Sprintf("CronJob set to run: %s", cronjob.Spec.Schedule))
	if err != nil {
		return errors.Wrap(err, errUpdateStatus)
	}

	return nil
}

// Helper function to update the CustomResource status.
func updateStatus(logger log.Logger, cr *v1alpha1.AcquiaSnapshotScheduled, status v1alpha1.Phase, message string) error {
	logger.Info(status, message)

	cr.Status = v1alpha1.AcquiaStatus{
		Phase:   status,
		Updated: time.Now(),
		Message: message,
	}

	return sdkaction.Update(cr)
}
