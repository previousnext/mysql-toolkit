package k8s

import (
	"fmt"

	"github.com/pkg/errors"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
)

// Helper function to extract a mysql connection/image from a ConfigMap key/value set.
func getBuildInfo(configmap corev1.ConfigMap, keys Keys) (string, string, string, string, string, error) {
	if _, ok := configmap.Data[keys.Hostname]; !ok {
		return "", "", "", "", "", fmt.Errorf("not found: %s", keys.Hostname)
	}

	if _, ok := configmap.Data[keys.Username]; !ok {
		return "", "", "", "", "", fmt.Errorf("not found: %s", keys.Username)
	}

	if _, ok := configmap.Data[keys.Password]; !ok {
		return "", "", "", "", "", fmt.Errorf("not found: %s", keys.Password)
	}

	if _, ok := configmap.Data[keys.Database]; !ok {
		return "", "", "", "", "", fmt.Errorf("not found: %s", keys.Database)
	}

	if _, ok := configmap.Data[keys.Image]; !ok {
		return "", "", "", "", "", fmt.Errorf("not found: %s", keys.Image)
	}

	return configmap.Data[keys.Hostname], configmap.Data[keys.Username], configmap.Data[keys.Password], configmap.Data[keys.Database], configmap.Data[keys.Image], nil
}

// Helper function generate resource requirements.
func resourceRequirements(r Resources) (corev1.ResourceRequirements, error) {
	cpu, err := resource.ParseQuantity(r.CPU)
	if err != nil {
		return corev1.ResourceRequirements{}, errors.Wrap(err, "failed to parse resource: cpu")
	}

	memory, err := resource.ParseQuantity(r.Memory)
	if err != nil {
		return corev1.ResourceRequirements{}, errors.Wrap(err, "failed to parse resource: memory")
	}

	return corev1.ResourceRequirements{
		Requests: corev1.ResourceList{
			corev1.ResourceCPU:    cpu,
			corev1.ResourceMemory: memory,
		},
		Limits: corev1.ResourceList{
			corev1.ResourceCPU:    cpu,
			corev1.ResourceMemory: memory,
		},
	}, nil
}
