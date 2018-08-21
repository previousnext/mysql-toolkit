package job

import (
	"fmt"

	"github.com/pkg/errors"
	yaml "gopkg.in/yaml.v2"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/previousnext/mysql-toolkit/internal/dumper"
	"github.com/previousnext/mysql-toolkit/internal/envar"
	"github.com/previousnext/mysql-toolkit/internal/operator/apis/mtk/v1alpha1"
	"github.com/previousnext/mysql-toolkit/internal/operator/handler/acquia/secrets"
)

const (
	// KeyAcquiaUsername for authentication.
	KeyAcquiaUsername = "acquia.username"
	// KeyAcquiaPassword for authentication.
	KeyAcquiaPassword = "acquia.password"
	// KeyDockerUsername for authentication.
	KeyDockerUsername = "docker.username"
	// KeyDockerPassword for authentication.
	KeyDockerPassword = "docker.password"
	// KeyAWSRole for AWS CodeBuild role assume.
	KeyAWSRole = "aws.role"
	// KeyAWSKey for authentication.
	KeyAWSKey = "aws.key.id"
	// KeyAWSAccess for authentication.
	KeyAWSAccess = "aws.key.access"
	// KeyAWSBucket for AWS CodeBuild to consume during a built.
	KeyAWSBucket = "aws.bucket"
	// KeyMtkConfig file for mtk ruleset.
	KeyMtkConfig = "mtk.yml"
)

// Params passed to GenerateSpec.
type Params struct {
	Namespace string
	Name      string
	Database  v1alpha1.AcquiaDatabase
	Docker    v1alpha1.Docker
	Secrets   secrets.Secrets
	Config    dumper.Config
	Image     string
	CPU       string
	Memory    string
}

// Validate the params.
func (p Params) Validate() error {
	if p.Namespace == "" {
		return errors.New("not found: namespace")
	}

	if p.Name == "" {
		return errors.New("not found: project")
	}

	if p.Database.Site == "" {
		return errors.New("not found: database: site")
	}

	if p.Database.Environment == "" {
		return errors.New("not found: database: environment")
	}

	if p.Database.Name == "" {
		return errors.New("not found: database: database")
	}

	if p.Docker.Image == "" {
		return errors.New("not found: docker: image")
	}

	if p.Image == "" {
		return errors.New("not found: image")
	}

	if p.CPU == "" {
		return errors.New("not found: cpu")
	}

	if p.Memory == "" {
		return errors.New("not found: memory")
	}

	return nil
}

// Generate a ConfigMap, Secret and Job object for executing a snapshot.
func Generate(params Params) (*batchv1.Job, *corev1.ConfigMap, *corev1.Secret, error) {
	err := params.Validate()
	if err != nil {
		return nil, nil, nil, errors.Wrap(err, "validation failed")
	}

	configmap, err := generateConfigMap(params)
	if err != nil {
		return nil, nil, nil, errors.Wrap(err, "failed to generate ConfigMap")
	}

	secret, err := generateSecret(params)
	if err != nil {
		return nil, nil, nil, errors.Wrap(err, "failed to generate Secret")
	}

	job, err := generateJob(params)
	if err != nil {
		return nil, nil, nil, errors.Wrap(err, "failed to generate Job")
	}

	return job, configmap, secret, nil
}

// Helper function to generate a ConfigMap.
func generateConfigMap(params Params) (*corev1.ConfigMap, error) {
	config, err := yaml.Marshal(&params.Config)
	if err != nil {
		return nil, errors.Wrap(err, "failed to generate mtk.yml")
	}

	return &corev1.ConfigMap{
		TypeMeta: metav1.TypeMeta{
			Kind:       "ConfigMap",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Namespace: params.Namespace,
			Name:      params.Name,
		},
		Data: map[string]string{
			KeyMtkConfig: string(config),
		},
	}, nil
}

// Helper function to generate a Secret.
func generateSecret(params Params) (*corev1.Secret, error) {
	return &corev1.Secret{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Secret",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Namespace: params.Namespace,
			Name:      params.Name,
		},
		StringData: map[string]string{
			KeyAcquiaUsername: params.Secrets.AcquiaUsername,
			KeyAcquiaPassword: params.Secrets.AcquiaPassword,
			KeyDockerUsername: params.Secrets.DockerUsername,
			KeyDockerPassword: params.Secrets.DockerPassword,
			KeyAWSRole:        params.Secrets.AWSRole,
			KeyAWSKey:         params.Secrets.AWSKey,
			KeyAWSAccess:      params.Secrets.AWSAccess,
			KeyAWSBucket:      params.Secrets.AWSBucket,
		},
	}, nil
}

// Helper function to generate a Job.
func generateJob(params Params) (*batchv1.Job, error) {
	// Backoff determines how many times the build fails before it does not get recreated.
	// This amount allows for any "transient" issues that could be fixed with a rerun.
	backoff := int32(2)

	cpu, err := resource.ParseQuantity(params.CPU)
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse: job: cpu")
	}

	memory, err := resource.ParseQuantity(params.Memory)
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse: job: memory")
	}

	resources := corev1.ResourceRequirements{
		Requests: corev1.ResourceList{
			corev1.ResourceCPU:    cpu,
			corev1.ResourceMemory: memory,
		},
		Limits: corev1.ResourceList{
			corev1.ResourceCPU:    cpu,
			corev1.ResourceMemory: memory,
		},
	}

	return &batchv1.Job{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Job",
			APIVersion: "batch/v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Namespace: params.Namespace,
			Name:      params.Name,
		},
		Spec: batchv1.JobSpec{
			BackoffLimit: &backoff,
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Namespace: params.Namespace,
				},
				Spec: corev1.PodSpec{
					RestartPolicy: corev1.RestartPolicyNever,
					InitContainers: []corev1.Container{
						{
							Name:  "version",
							Image: params.Image,
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
							Image: params.Image,
							Env: []corev1.EnvVar{
								{
									Name: envar.AcquiaUsername,
									ValueFrom: &corev1.EnvVarSource{
										SecretKeyRef: &corev1.SecretKeySelector{
											Key: KeyAcquiaUsername,
											LocalObjectReference: corev1.LocalObjectReference{
												Name: params.Name,
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
												Name: params.Name,
											},
										},
									},
								},
								{
									Name:  envar.AcquiaSite,
									Value: params.Database.Site,
								},
								{
									Name:  envar.AcquiaEnvironment,
									Value: params.Database.Environment,
								},
								{
									Name:  envar.AcquiaDatabase,
									Value: params.Database.Name,
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
							Image: params.Image,
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
							Image:           params.Image,
							ImagePullPolicy: "Always",
							Env: []corev1.EnvVar{
								{
									Name:  envar.AWSCodeBuildProject,
									Value: params.Name,
								},
								{
									Name: envar.AWSIAMRole,
									ValueFrom: &corev1.EnvVarSource{
										SecretKeyRef: &corev1.SecretKeySelector{
											Key: KeyAWSRole,
											LocalObjectReference: corev1.LocalObjectReference{
												Name: params.Name,
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
												Name: params.Name,
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
												Name: params.Name,
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
												Name: params.Name,
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
												Name: params.Name,
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
												Name: params.Name,
											},
										},
									},
								},
								{
									Name:  envar.DockerImage,
									Value: params.Docker.Image,
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
										Name: params.Name,
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
