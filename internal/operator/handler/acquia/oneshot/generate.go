package oneshot

import (
	"fmt"

	"github.com/pkg/errors"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/previousnext/mysql-toolkit/internal/envar"
)

const (
	// KeyAcquiaUsername for loading Acquia Cloud credentials.
	KeyAcquiaUsername = "acquia.username"
	// KeyAcquiaPassword for loading Acquia Cloud credentials.
	KeyAcquiaPassword = "acquia.password"
	// KeyDockerUsername for loading Docker Registry credentials.
	KeyDockerUsername = "docker.username"
	// KeyDockerPassword for loading Docker Registry credentials.
	KeyDockerPassword = "docker.password"
	// KeyAWSRole for building a MySQL image.
	KeyAWSRole = "aws.role"
	// KeyAWSBucket for building a MySQL image.
	KeyAWSBucket = "aws.bucket"
	// KeyAWSKeyID for building a MySQL image.
	KeyAWSKeyID = "aws.key.id"
	// KeyAWSKeyAccess for building a MySQL image.
	KeyAWSKeyAccess = "aws.key.access"
)

// GenerateJobSpecInput passed to GenerateJobSpec.
type GenerateJobSpecInput struct {
	ObjectMeta  metav1.ObjectMeta
	Site        string
	Environment string
	Database    string
	Image       string
	Credentials GenerateJobSpecInputCredentials
	Job         GenerateJobSpecInputJob
}

// GenerateJobSpecInputCredentials for dumping and building a mysql image.
type GenerateJobSpecInputCredentials struct {
	Docker string
	Acquia string
	AWS    string
}

// GenerateJobSpecInputJob configuration for running the image through a pipeline.
type GenerateJobSpecInputJob struct {
	Image     string
	Resources corev1.ResourceRequirements
}

// GenerateJobSpec which will return a job to dump, sanitize and build a Acquia database image.
func GenerateJobSpec(input GenerateJobSpecInput) (batchv1.JobSpec, error) {
	if input.Site == "" {
		return batchv1.JobSpec{}, errors.New("not found: site")
	}

	if input.Environment == "" {
		return batchv1.JobSpec{}, errors.New("not found: environment")
	}

	if input.Database == "" {
		return batchv1.JobSpec{}, errors.New("not found: database")
	}

	if input.Image == "" {
		return batchv1.JobSpec{}, errors.New("not found: image")
	}

	if input.Job.Image == "" {
		return batchv1.JobSpec{}, errors.New("not found: job: image")
	}

	if input.Credentials.Docker == "" {
		return batchv1.JobSpec{}, errors.New("not found: credentials: docker")
	}

	if input.Credentials.Acquia == "" {
		return batchv1.JobSpec{}, errors.New("not found: credentials: acquia")
	}

	if input.Credentials.AWS == "" {
		return batchv1.JobSpec{}, errors.New("not found: credentials: aws")
	}

	// Backoff determines how many times the build fails before it does not get recreated.
	// This is set to 2 for:
	//  * Generally a fail will be OOMKiller not happy with how much memory awscli is using,
	//    we shouldn't run builds over and over again, they will keep failing.
	//  * This amount allows for any "transient" issues that could be fixed with a rerun.
	backoff := int32(2)

	// Provide a common mount so our build pipeline can transform the database.
	// eg. dump, sanitize and build in separate containers.
	mounts := []corev1.VolumeMount{
		{
			Name:      "tmp",
			MountPath: "/tmp",
		},
	}

	return batchv1.JobSpec{
		BackoffLimit: &backoff,
		Template: corev1.PodTemplateSpec{
			ObjectMeta: metav1.ObjectMeta{
				Namespace: input.ObjectMeta.Namespace,
			},
			Spec: corev1.PodSpec{
				RestartPolicy: corev1.RestartPolicyNever,
				InitContainers: []corev1.Container{
					{
						Name:  "version",
						Image: input.Job.Image,
						Command: []string{
							"/bin/sh", "-c",
						},
						Args: []string{
							"mtk version",
						},
						ImagePullPolicy: corev1.PullAlways,
						Resources:       input.Job.Resources,
					},
					{
						Name:  "dump",
						Image: input.Job.Image,
						Env: []corev1.EnvVar{
							{
								Name: envar.AcquiaUsername,
								ValueFrom: &corev1.EnvVarSource{
									SecretKeyRef: &corev1.SecretKeySelector{
										LocalObjectReference: corev1.LocalObjectReference{
											Name: input.Credentials.Acquia,
										},
										Key: KeyAcquiaUsername,
									},
								},
							},
							{
								Name: envar.AcquiaPassword,
								ValueFrom: &corev1.EnvVarSource{
									SecretKeyRef: &corev1.SecretKeySelector{
										LocalObjectReference: corev1.LocalObjectReference{
											Name: input.Credentials.Acquia,
										},
										Key: KeyAcquiaPassword,
									},
								},
							},
							{
								Name:  envar.AcquiaSite,
								Value: input.Site,
							},
							{
								Name:  envar.AcquiaEnvironment,
								Value: input.Environment,
							},
							{
								Name:  envar.AcquiaDatabase,
								Value: input.Database,
							},
						},
						Command: []string{
							"/bin/sh", "-c",
						},
						Args: []string{
							"mtk acquia dump --file=/tmp/unsanitized.sql",
						},
						ImagePullPolicy: corev1.PullIfNotPresent,
						Resources:       input.Job.Resources,
						VolumeMounts:    mounts,
					},
					{
						Name:  "sanitize",
						Image: input.Job.Image,
						Command: []string{
							"/bin/sh", "-c",
						},
						Args: []string{
							"database-sanitize /tmp/unsanitized.sql /tmp/db.sql",
						},
						ImagePullPolicy: corev1.PullIfNotPresent,
						Resources:       input.Job.Resources,
						VolumeMounts:    mounts,
					},
				},
				Containers: []corev1.Container{
					corev1.Container{
						Name:            "codebuild",
						Image:           input.Job.Image,
						ImagePullPolicy: "Always",
						Env: []corev1.EnvVar{
							{
								Name:  envar.AWSCodeBuildProject,
								Value: fmt.Sprintf("mtk-%s-%s", input.ObjectMeta.Namespace, input.ObjectMeta.Name),
							},
							{
								Name: envar.AWSIAMRole,
								ValueFrom: &corev1.EnvVarSource{
									SecretKeyRef: &corev1.SecretKeySelector{
										LocalObjectReference: corev1.LocalObjectReference{
											Name: input.Credentials.AWS,
										},
										Key: input.Credentials.AWS,
									},
								},
							},
							{
								Name: envar.AWSS3Bucket,
								ValueFrom: &corev1.EnvVarSource{
									SecretKeyRef: &corev1.SecretKeySelector{
										LocalObjectReference: corev1.LocalObjectReference{
											Name: input.Credentials.AWS,
										},
										Key: input.Credentials.AWS,
									},
								},
							},
							{
								Name: envar.AWSAccessKeyID,
								ValueFrom: &corev1.EnvVarSource{
									SecretKeyRef: &corev1.SecretKeySelector{
										LocalObjectReference: corev1.LocalObjectReference{
											Name: input.Credentials.AWS,
										},
										Key: input.Credentials.AWS,
									},
								},
							},
							{
								Name: envar.AWSSecretAccessKey,
								ValueFrom: &corev1.EnvVarSource{
									SecretKeyRef: &corev1.SecretKeySelector{
										LocalObjectReference: corev1.LocalObjectReference{
											Name: input.Credentials.AWS,
										},
										Key: input.Credentials.AWS,
									},
								},
							},
							{
								Name: envar.DockerUsername,
								ValueFrom: &corev1.EnvVarSource{
									SecretKeyRef: &corev1.SecretKeySelector{
										LocalObjectReference: corev1.LocalObjectReference{
											Name: input.Credentials.Docker,
										},
										Key: KeyDockerUsername,
									},
								},
							},
							{
								Name: envar.DockerPassword,
								ValueFrom: &corev1.EnvVarSource{
									SecretKeyRef: &corev1.SecretKeySelector{
										LocalObjectReference: corev1.LocalObjectReference{
											Name: input.Credentials.Docker,
										},
										Key: KeyDockerPassword,
									},
								},
							},
							{
								Name:  envar.DockerImage,
								Value: input.Image,
							},
						},
						Command: []string{
							"/bin/sh", "-c",
						},
						Args: []string{
							"mtk build aws --file=/tmp/db.sql",
						},
						Resources:    input.Job.Resources,
						VolumeMounts: mounts,
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
	}, nil
}
