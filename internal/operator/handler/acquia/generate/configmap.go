package generate

import (
	"github.com/pkg/errors"
	yaml "gopkg.in/yaml.v2"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/previousnext/mysql-toolkit/internal/dumper"
)

// ConfigMap associated with a Job.
func ConfigMap(namespace, name string, config dumper.Config) (*corev1.ConfigMap, error) {
	raw, err := yaml.Marshal(&config)
	if err != nil {
		return nil, errors.Wrap(err, "failed to generate mtk.yml")
	}

	return &corev1.ConfigMap{
		TypeMeta: metav1.TypeMeta{
			Kind:       "ConfigMap",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Namespace: namespace,
			Name:      name,
		},
		Data: map[string]string{
			KeyMtkConfig: string(raw),
		},
	}, nil
}
