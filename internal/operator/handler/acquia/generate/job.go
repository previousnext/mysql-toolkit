package generate

import (
	"fmt"

	"github.com/pkg/errors"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/previousnext/mysql-toolkit/internal/envar"
	"github.com/previousnext/mysql-toolkit/internal/operator/apis/mtk/v1alpha1"
)

// Job for executing an Acquia database snapshot.
func Job(namespace, name, image, cpu, mem string, source v1alpha1.AcquiaDatabase, target v1alpha1.Docker) (*batchv1.Job, error) {
	// Backoff determines how many times the build fails before it does not get recreated.
	// This amount allows for any "transient" issues that could be fixed with a rerun.
	backoff := int32(2)

	resCPU, err := resource.ParseQuantity(cpu)
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse: cpu")
	}

	resMem, err := resource.ParseQuantity(mem)
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse: memory")
	}

	resources := corev1.ResourceRequirements{
		Requests: corev1.ResourceList{
			corev1.ResourceCPU:    resCPU,
			corev1.ResourceMemory: resMem,
		},
		Limits: corev1.ResourceList{
			corev1.ResourceCPU:    resCPU,
			corev1.ResourceMemory: resMem,
		},
	}

	return &batchv1.Job{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Job",
			APIVersion: "batch/v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Namespace: namespace,
			Name:      name,
		},
		Spec: batchv1.JobSpec{
			BackoffLimit: &backoff,
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Namespace: namespace,
				},
				Spec: corev1.PodSpec{
					RestartPolicy: corev1.RestartPolicyNever,
					InitContainers: []corev1.Container{
						{
							Name:  "version",
							Image: image,
							Command: []string{
								"/bin/sh", "-c",
							},
							Args: []string{
								"mtk version",
							},
							ImagePullPolicy: corev1.PullAlways,
							Resources:       resources,
						},
						{
							Name:  "dump",
							Image: image,
							Env: []corev1.EnvVar{
								{
									Name: envar.AcquiaUsername,
									ValueFrom: &corev1.EnvVarSource{
										SecretKeyRef: &corev1.SecretKeySelector{
											Key: KeyAcquiaUsername,
											LocalObjectReference: corev1.LocalObjectReference{
												Name: name,
											},
										},
									},
								},
								{
									Name: envar.AcquiaPassword,
									ValueFrom: &corev1.EnvVarSource{
										SecretKeyRef: &corev1.SecretKeySelector{
											Key: KeyAcquiaPassword,
											LocalObjectReference: corev1.LocalObjectReference{
												Name: name,
											},
										},
									},
								},
								{
									Name:  envar.AcquiaSite,
									Value: source.Site,
								},
								{
									Name:  envar.AcquiaEnvironment,
									Value: source.Environment,
								},
								{
									Name:  envar.AcquiaDatabase,
									Value: source.Name,
								},
							},
							Command: []string{
								"/bin/sh", "-c",
							},
							Args: []string{
								"mtk acquia dump --file=/tmp/unsanitized.sql",
							},
							ImagePullPolicy: corev1.PullIfNotPresent,
							Resources:       resources,
							VolumeMounts: []corev1.VolumeMount{
								{
									Name:      "tmp",
									MountPath: "/tmp",
								},
							},
						},
						{
							Name:  "sanitize",
							Image: image,
							Env: []corev1.EnvVar{
								{
									Name:  envar.MySQLConfig,
									Value: fmt.Sprintf("/config/%s", KeyMtkConfig),
								},
							},
							Command: []string{
								"/bin/sh", "-c",
							},
							Args: []string{
								"database-sanitize /tmp/unsanitized.sql /tmp/db.sql",
							},
							ImagePullPolicy: corev1.PullIfNotPresent,
							Resources:       resources,
							VolumeMounts: []corev1.VolumeMount{
								{
									Name:      "config",
									MountPath: "/config",
								},
								{
									Name:      "tmp",
									MountPath: "/tmp",
								},
							},
						},
					},
					Containers: []corev1.Container{
						corev1.Container{
							Name:            "codebuild",
							Image:           image,
							ImagePullPolicy: "Always",
							Env: []corev1.EnvVar{
								{
									Name:  envar.AWSCodeBuildProject,
									Value: name,
								},
								{
									Name: envar.AWSIAMRole,
									ValueFrom: &corev1.EnvVarSource{
										SecretKeyRef: &corev1.SecretKeySelector{
											Key: KeyAWSRole,
											LocalObjectReference: corev1.LocalObjectReference{
												Name: name,
											},
										},
									},
								},
								{
									Name: envar.AWSS3Bucket,
									ValueFrom: &corev1.EnvVarSource{
										SecretKeyRef: &corev1.SecretKeySelector{
											Key: KeyAWSBucket,
											LocalObjectReference: corev1.LocalObjectReference{
												Name: name,
											},
										},
									},
								},
								{
									Name: envar.AWSAccessKeyID,
									ValueFrom: &corev1.EnvVarSource{
										SecretKeyRef: &corev1.SecretKeySelector{
											Key: KeyAWSKey,
											LocalObjectReference: corev1.LocalObjectReference{
												Name: name,
											},
										},
									},
								},
								{
									Name: envar.AWSSecretAccessKey,
									ValueFrom: &corev1.EnvVarSource{
										SecretKeyRef: &corev1.SecretKeySelector{
											Key: KeyAWSAccess,
											LocalObjectReference: corev1.LocalObjectReference{
												Name: name,
											},
										},
									},
								},
								{
									Name: envar.DockerUsername,
									ValueFrom: &corev1.EnvVarSource{
										SecretKeyRef: &corev1.SecretKeySelector{
											Key: KeyDockerUsername,
											LocalObjectReference: corev1.LocalObjectReference{
												Name: name,
											},
										},
									},
								},
								{
									Name: envar.DockerPassword,
									ValueFrom: &corev1.EnvVarSource{
										SecretKeyRef: &corev1.SecretKeySelector{
											Key: KeyDockerPassword,
											LocalObjectReference: corev1.LocalObjectReference{
												Name: name,
											},
										},
									},
								},
								{
									Name:  envar.DockerImage,
									Value: target.Image,
								},
							},
							Command: []string{
								"/bin/sh", "-c",
							},
							Args: []string{
								"mtk build aws --file=/tmp/db.sql",
							},
							Resources: resources,
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
							Name: "config",
							VolumeSource: corev1.VolumeSource{
								ConfigMap: &corev1.ConfigMapVolumeSource{
									LocalObjectReference: corev1.LocalObjectReference{
										Name: name,
									},
									Items: []corev1.KeyToPath{
										{
											Key:  KeyMtkConfig,
											Path: KeyMtkConfig,
										},
									},
								},
							},
						},
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
	}, nil
}
