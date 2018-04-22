package cmd

import (
	"os"

	"gopkg.in/alecthomas/kingpin.v2"

	cmdenv "github.com/previousnext/mysql-toolkit/cmd/env"
	"github.com/previousnext/mysql-toolkit/codebuilder"
)

type cmdCodeBuild struct {
	Region string
	params codebuilder.BuildParams
}

func (cmd *cmdCodeBuild) run(c *kingpin.ParseContext) error {
	return codebuilder.Build(os.Stdout, &cmd.params)
}

// CodeBuild declares the "codebuild" subcommand.
func CodeBuild(app *kingpin.Application) {
	c := new(cmdCodeBuild)

	cmd := app.Command("codebuild", "Build a container using AWS CodeBuild").Action(c.run)

	cmd.Flag("aws-region", "Region to run ").Default("ap-southeast-2").Envar(cmdenv.AWSRegion).StringVar(&c.params.Region)

	cmd.Flag("project", "Name for the CodeBuild project").Required().Envar(cmdenv.AWSCodeBuildProject).StringVar(&c.params.Project)
	cmd.Flag("compute", "Size of the compute for the build").Default(codebuilder.DefaultComputeSize).Envar(cmdenv.AWSCodeBuildCompute).StringVar(&c.params.Compute)
	cmd.Flag("image", "CodeBuild image to use for executing the build").Default(codebuilder.DefaultBuildImage).Envar(cmdenv.AWSCodeBuildImage).StringVar(&c.params.Image)
	cmd.Flag("dockerfile", "Path to the Dockerfile use to build the image").Required().Envar(cmdenv.AWSCodeBuildDockerfile).StringVar(&c.params.Dockerfile)
	cmd.Flag("spec", "Path to the BuildSpec use to build the image").Required().Envar(cmdenv.AWSCodeBuildSpec).StringVar(&c.params.BuildSpec)
	cmd.Flag("bucket", "Bucket to upload the file temporarily before CodeBuild runs").Required().Envar(cmdenv.AWSS3Bucket).StringVar(&c.params.Bucket)
	cmd.Flag("role", "ServiceRole or IAM resource which grants access to the S3 bucket").Required().Envar(cmdenv.AWSIAMRole).StringVar(&c.params.Role)
	cmd.Flag("docker-username", "Username for the Docker Registry").Required().Envar(cmdenv.DockerUsername).StringVar(&c.params.Docker.Username)
	cmd.Flag("docker-password", "Password for the Docker Registry").Required().Envar(cmdenv.DockerPassword).StringVar(&c.params.Docker.Password)
	cmd.Flag("docker-image", "Name to push to the registry").Required().Envar(cmdenv.DockerImage).StringVar(&c.params.Docker.Image)
	cmd.Flag("file", "Path to the Mysql database use to build the image").Required().Envar(cmdenv.MySQLFile).StringVar(&c.params.Database)
}
