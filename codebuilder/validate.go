package codebuilder

import "github.com/pkg/errors"

func (b *BuildParams) Validate() error {
	if b.Project == "" {
		return errors.New("not found: project")
	}

	if b.Region == "" {
		return errors.New("not found: region")
	}

	if b.Dockerfile == "" {
		return errors.New("not found: source: Dockerfile")
	}

	if b.BuildSpec == "" {
		return errors.New("not found: source: buildspec")
	}

	if b.Database == "" {
		return errors.New("not found: source: database")
	}

	if b.Docker.Username == "" {
		return errors.New("not found: Docker username")
	}

	if b.Docker.Password == "" {
		return errors.New("not found: Docker password")
	}

	if b.Docker.Image == "" {
		return errors.New("not found: Docker image")
	}

	return nil
}
