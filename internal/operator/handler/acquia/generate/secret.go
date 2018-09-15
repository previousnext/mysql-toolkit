package generate

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/previousnext/mysql-toolkit/internal/operator/handler/acquia/secrets"
)

// Secret associated with a Job.
func Secret(namespace, name string, values secrets.Secrets) (*corev1.Secret, error) {
	return &corev1.Secret{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Secret",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Namespace: namespace,
			Name:      name,
		},
		StringData: map[string]string{
			KeyAcquiaUsername: values.AcquiaUsername,
			KeyAcquiaPassword: values.AcquiaPassword,
			KeyDockerUsername: values.DockerUsername,
			KeyDockerPassword: values.DockerPassword,
			KeyAWSRole:        values.AWSRole,
			KeyAWSKey:         values.AWSKey,
			KeyAWSAccess:      values.AWSAccess,
			KeyAWSBucket:      values.AWSBucket,
		},
	}, nil
}
