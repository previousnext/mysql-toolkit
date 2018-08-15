package tests

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/nickschuch/k8sutils/configmap"
)

func TestMarshal(t *testing.T) {
	backend := Backend{
		Host:   "1.1.1.1",
		Port:   443,
		Secure: true,
	}

	data, err := configmap.Marshal(backend)
	assert.Nil(t, err)

	assert.Equal(t, "1.1.1.1", data["backend.host"])
	assert.Equal(t, "443", data["backend.port"])
	assert.Equal(t, "true", data["backend.secure"])
}
