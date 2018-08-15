package tests

import (
	"testing"

	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"

	"github.com/nickschuch/k8sutils/secret"
)

func TestUnmarshal(t *testing.T) {
	auth := Authentication{
		Username: "example",
		Password: "password123",
	}

	data, err := secret.Marshal(auth)
	assert.Nil(t, err)

	var newAuth Authentication

	err = secret.Unmarshal(corev1.Secret{Data: data}, &newAuth)
	assert.Nil(t, err)

	assert.Equal(t, auth.Username, newAuth.Username)
	assert.Equal(t, auth.Password, newAuth.Password)
}
