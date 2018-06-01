package builder

import (
	"fmt"
	"io"
	"os"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/pkg/errors"
	"golang.org/x/net/context"

	"github.com/previousnext/mysql-toolkit/internal/docker/auth"
	"github.com/previousnext/mysql-toolkit/internal/docker/stream"
	"github.com/previousnext/mysql-toolkit/internal/packager"
)

// BuildParams to pass to the Build function.
type BuildParams struct {
	Dockerfile string
	Image      string
	Username   string
	Password   string
	Database   string
}

// Build a MySQL image using Docker.
func Build(w io.Writer, params BuildParams) error {
	cli, err := client.NewEnvClient()
	if err != nil {
		return errors.Wrap(err, "failed to setup client")
	}

	fmt.Println("Building the image...")

	path, err := packager.Tar(packager.TarParams{
		Files: []packager.File{
			{
				Local:  params.Dockerfile,
				Remote: packager.Dockerfile,
			},
			{
				Local:  params.Database,
				Remote: packager.Database,
			},
		},
	})
	if err != nil {
		errors.Wrap(err, "failed to build context")
	}

	ctx, err := os.Open(path)
	if err != nil {
		errors.Wrap(err, "failed to open build context")
	}
	defer os.Remove(path)

	build, err := cli.ImageBuild(context.Background(), ctx, types.ImageBuildOptions{
		Tags: []string{
			params.Image,
		},
	})
	if err != nil {
		return errors.Wrap(err, "failed to build image")
	}

	err = stream.Build(os.Stdout, build.Body)
	if err != nil {
		errors.Wrap(err, "failed to stream build output")
	}

	fmt.Println("Pushing the image to the registry...")

	authBase64, err := auth.Base64(params.Username, params.Password)
	if err != nil {
		return errors.Wrap(err, "failed to setup credentials")
	}

	push, err := cli.ImagePush(context.Background(), params.Image, types.ImagePushOptions{
		RegistryAuth: authBase64,
	})
	if err != nil {
		errors.Wrap(err, "failed to stream push output")
	}

	return stream.Push(os.Stdout, push)
}
