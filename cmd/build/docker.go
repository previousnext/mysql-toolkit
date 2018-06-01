package build

import (
	"os"

	"gopkg.in/alecthomas/kingpin.v2"

	cmdenv "github.com/previousnext/mysql-toolkit/cmd/env"
	"github.com/previousnext/mysql-toolkit/internal/docker/builder"
)

type cmdDocker struct {
	params builder.BuildParams
}

func (cmd *cmdDocker) run(c *kingpin.ParseContext) error {
	return builder.Build(os.Stdout, cmd.params)
}

// Docker declares the "docker" subcommand.
func Docker(app *kingpin.CmdClause) {
	c := new(cmdDocker)

	cmd := app.Command("docker", "Build an image using Docker").Action(c.run)

	cmd.Flag("dockerfile", "Path to the Dockerfile use to build the image").Required().Envar(cmdenv.Dockerfile).StringVar(&c.params.Dockerfile)
	cmd.Flag("username", "Username for the Docker Registry").Required().Envar(cmdenv.DockerUsername).StringVar(&c.params.Username)
	cmd.Flag("password", "Password for the Docker Registry").Required().Envar(cmdenv.DockerPassword).StringVar(&c.params.Password)
	cmd.Flag("image", "Name to push to the registry").Required().Envar(cmdenv.DockerImage).StringVar(&c.params.Image)
	cmd.Flag("file", "Path to the Mysql database use to build the image").Required().Envar(cmdenv.MySQLFile).StringVar(&c.params.Database)
}
