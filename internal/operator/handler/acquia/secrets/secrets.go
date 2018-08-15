package secrets

import (
	"github.com/nickschuch/k8sutils/secret"
	"github.com/pkg/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

// Load Acquia config from multiple ConfigMaps.
func Load(opNamespace, opSecret, crNamespace, crSecret string) (Secrets, error) {
	var compiled Secrets

	k8sconfig, err := rest.InClusterConfig()
	if err != nil {
		return compiled, errors.Wrap(err, "failed to get K8s config")
	}

	clientset, err := kubernetes.NewForConfig(k8sconfig)
	if err != nil {
		return compiled, errors.Wrap(err, "failed to get K8s client")
	}

	err = load(clientset, opNamespace, opSecret, &compiled)
	if err != nil {
		return compiled, errors.Wrap(err, "failed to load operator configuration")
	}

	err = load(clientset, crNamespace, crSecret, &compiled)
	if err != nil {
		return compiled, errors.Wrap(err, "failed to load project configuration")
	}

	err = validate(compiled)
	if err != nil {
		return compiled, errors.Wrap(err, "valildation failed")
	}

	return compiled, nil
}

// Helper function to load a config.
func load(clientset *kubernetes.Clientset, namespace, secretName string, cfg *Secrets) error {
	if namespace == "" {
		return errors.New("not found: namespace")
	}

	if secretName == "" {
		return errors.New("not found: secret")
	}

	s, err := clientset.CoreV1().Secrets(namespace).Get(secretName, metav1.GetOptions{})
	if err != nil {
		return errors.Wrap(err, "failed to load Secret")
	}

	err = secret.Unmarshal(s, cfg)
	if err != nil {
		return errors.Wrap(err, "failed to unmarshal Secret")
	}

	return err
}

// Helper function to validate config.
func validate(s Secrets) error {
	if s.AcquiaUsername == "" {
		return errors.New("not found: acquia: username")
	}

	if s.AcquiaPassword == "" {
		return errors.New("not found: acquia: password")
	}

	if s.DockerUsername == "" {
		return errors.New("not found: docker: username")
	}

	if s.DockerPassword == "" {
		return errors.New("not found: docker: password")
	}

	if s.AWSRole == "" {
		return errors.New("not found: aws: role")
	}

	if s.AWSBucket == "" {
		return errors.New("not found: aws: bucket")
	}

	if s.AWSKey == "" {
		return errors.New("not found: aws: key")
	}

	if s.AWSAccess == "" {
		return errors.New("not found: aws: access")
	}

	return nil
}
