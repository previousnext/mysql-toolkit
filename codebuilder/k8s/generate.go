package k8s

import (
	"fmt"

	"github.com/pkg/errors"
	batchv1 "k8s.io/api/batch/v1"
	batchv1beta1 "k8s.io/api/batch/v1beta1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	cmdenv "github.com/previousnext/mysql-toolkit/cmd/env"
)

// Helper function to convert a PersistentVolumeClaim into a backup CronJob task.
func generateCronJob(configmap corev1.ConfigMap, params SyncParams) (*batchv1beta1.CronJob, error) {
	cronjob := &batchv1beta1.CronJob{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: configmap.ObjectMeta.Namespace,
			Name:      fmt.Sprintf("mysql-toolkit-codebuild-%s", configmap.ObjectMeta.Name),
		},
	}

	var (
		// Backoff determines how many times the build fails before it does not get recreated.
		// This is set to 2 for:
		//  * Generally a fail will be OOMKiller not happy with how much memory awscli is using,
		//    we shouldn't run builds over and over again, they will keep failing.
		//  * This amount allows for any "transient" issues that could be fixed with a rerun.
		backoff int32 = 2
		// CronJobs will have to start within 30 min.
		// https://kubernetes.io/docs/concepts/workloads/controllers/cron-jobs/#starting-deadline-seconds
		deadline int64 = 1800
	)

	hostname, username, password, database, image, err := getBuildInfo(configmap, params.Keys)
	if err != nil {
		return cronjob, errors.Wrap(err, "failed to get MySQL connection info")
	}

	resources, err := resourceRequirements(params.Resources)
	if err != nil {
		return cronjob, errors.Wrap(err, "failed to build resource requirements")
	}

	cronjob.Spec = batchv1beta1.CronJobSpec{
		Schedule:                params.Frequency,
		ConcurrencyPolicy:       batchv1beta1.ForbidConcurrent,
		StartingDeadlineSeconds: &deadline,
		JobTemplate: batchv1beta1.JobTemplateSpec{
			ObjectMeta: metav1.ObjectMeta{
				Namespace: configmap.ObjectMeta.Namespace,
			},
			Spec: batchv1.JobSpec{
				BackoffLimit: &backoff,
				Template: corev1.PodTemplateSpec{
					ObjectMeta: metav1.ObjectMeta{
						Namespace: configmap.ObjectMeta.Namespace,
					},
					Spec: corev1.PodSpec{
						RestartPolicy: "Never",
						InitContainers: []corev1.Container{
							{
								Name:  "version",
								Image: params.Image,
								Command: []string{
									"/bin/sh", "-c",
								},
								Args: []string{
									fmt.Sprintf("mysql-toolkit version"),
								},
								Resources:       resources,
								ImagePullPolicy: "Always",
							},
							{
								Name:  "dump",
								Image: params.Image,
								Env: []corev1.EnvVar{
									{
										Name:  cmdenv.MySQLFile,
										Value: "/tmp/db.sql",
									},
									{
										Name:  cmdenv.MySQLHostname,
										Value: hostname,
									},
									{
										Name:  cmdenv.MySQLUsername,
										Value: username,
									},
									{
										Name:  cmdenv.MySQLPassword,
										Value: password,
									},
									{
										Name:  cmdenv.MySQLDatabase,
										Value: database,
									},
								},
								Command: []string{
									"/bin/sh", "-c",
								},
								Args: []string{
									fmt.Sprintf("mysql-toolkit dump"),
								},
								Resources:       resources,
								ImagePullPolicy: "Always",
								VolumeMounts: []corev1.VolumeMount{
									{
										Name:      "tmp",
										MountPath: "/tmp",
									},
								},
							},
						},
						Containers: []corev1.Container{
							{
								Name:  "codebuild",
								Image: params.Image,
								Env: []corev1.EnvVar{
									{
										Name:  cmdenv.AWSCodeBuildProject,
										Value: fmt.Sprintf("mysql-toolkit-%s-%s", configmap.ObjectMeta.Namespace, configmap.ObjectMeta.Name),
									},
									{
										Name:  cmdenv.AWSIAMRole,
										Value: params.Role,
									},
									{
										Name:  cmdenv.AWSS3Bucket,
										Value: params.Bucket,
									},
									{
										Name:  cmdenv.AWSAccessKeyID,
										Value: params.AWS.KeyID,
									},
									{
										Name:  cmdenv.AWSSecretAccessKey,
										Value: params.AWS.AccessKey,
									},
									{
										Name:  cmdenv.DockerUsername,
										Value: params.Docker.Username,
									},
									{
										Name:  cmdenv.DockerPassword,
										Value: params.Docker.Password,
									},
									{
										Name:  cmdenv.DockerImage,
										Value: image,
									},
									{
										Name:  cmdenv.MySQLFile,
										Value: "/tmp/db.sql",
									},
								},
								Command: []string{
									"/bin/sh", "-c",
								},
								Args: []string{
									fmt.Sprintf("mysql-toolkit codebuild"),
								},
								Resources:       resources,
								ImagePullPolicy: "Always",
								VolumeMounts: []corev1.VolumeMount{
									{
										Name:      "tmp",
										MountPath: "/tmp",
									},
								},
							},
						},
						Volumes: []corev1.Volume{
							{
								Name: "tmp",
								VolumeSource: corev1.VolumeSource{
									EmptyDir: &corev1.EmptyDirVolumeSource{
										Medium: corev1.StorageMediumDefault,
									},
								},
							},
						},
					},
				},
			},
		},
	}

	return cronjob, nil
}
