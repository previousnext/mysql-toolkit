package auth

import (
	"encoding/base64"
	"encoding/json"

	"github.com/docker/docker/api/types"
	"github.com/pkg/errors"
)

// Base64 encode the Docker registry authentication credentials.
func Base64(username, password string) (string, error) {
	auth := types.AuthConfig{
		Username: username,
		Password: password,
	}

	authBytes, err := json.Marshal(auth)
	if err != nil {
		return "", errors.Wrap(err, "failed to marshal auth")
	}

	return base64.URLEncoding.EncodeToString(authBytes), nil
}
