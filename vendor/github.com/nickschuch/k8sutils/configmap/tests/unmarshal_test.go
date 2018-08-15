package tests

import (
	"testing"

	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"

	"github.com/nickschuch/k8sutils/configmap"
)

func TestUnmarshal(t *testing.T) {
	cfg := corev1.ConfigMap{
		Data: map[string]string{
			"backend.host":   "1.1.1.1",
			"backend.port":   "443",
			"backend.secure": "true",
		},
	}

	var backend Backend

	err := configmap.Unmarshal(cfg, &backend)
	assert.Nil(t, err)

	assert.Equal(t, "1.1.1.1", backend.Host)
	assert.Equal(t, int64(443), backend.Port)
	assert.True(t, backend.Secure)
}
