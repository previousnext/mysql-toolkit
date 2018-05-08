package cmd

import (
	"os"

	"gopkg.in/alecthomas/kingpin.v2"
	corev1 "k8s.io/api/core/v1"

	cmdenv "github.com/previousnext/mysql-toolkit/cmd/env"
	codebuilderk8s "github.com/previousnext/mysql-toolkit/internal/codebuilder/k8s"
)

type cmdCodeBuildK8s struct {
	params codebuilderk8s.SyncParams
}

func (cmd *cmdCodeBuildK8s) run(c *kingpin.ParseContext) error {
	return codebuilderk8s.Sync(os.Stdout, cmd.params)
}

// CodeBuildK8s declares the "codebuild-k8s" subcommand.
func CodeBuildK8s(app *kingpin.Application) {
	c := new(cmdCodeBuildK8s)

	cmd := app.Command("codebuild-k8s", "Build a container using AWS CodeBuild").Action(c.run)

	cmd.Flag("namespace", "Namespace to lookup ConfigMaps").Default(corev1.NamespaceAll).Envar(cmdenv.K8sNamespace).StringVar(&c.params.Namespace)
	cmd.Flag("frequency", "How ofter CronJobs should create new CodeBuild project builds").Default("@daily").Envar(cmdenv.K8sCronJobFrequency).StringVar(&c.params.Frequency)
	cmd.Flag("image", "Image to use for running the CronJob").Default("previousnext/mysql-toolkit:latest").Envar(cmdenv.K8sCronJobImage).StringVar(&c.params.Image)
	cmd.Flag("cpu", "How much CPU resource should be assigned to the CronJob").Default("250m").Envar(cmdenv.K8sCronJobCPU).StringVar(&c.params.Resources.CPU)
	cmd.Flag("memory", "How much memory resource should be assigned to the CronJob").Default("256Mi").Envar(cmdenv.K8sCronJobMemory).StringVar(&c.params.Resources.Memory)
	cmd.Flag("key-hostname", "ConfigMap key which containers the MySQL hostname").Default("mysql.hostname").Envar(cmdenv.K8sConfigMapKeyHostname).StringVar(&c.params.Keys.Hostname)
	cmd.Flag("key-username", "ConfigMap key which containers the MySQL username").Default("mysql.username").Envar(cmdenv.K8sConfigMapKeyUsername).StringVar(&c.params.Keys.Username)
	cmd.Flag("key-password", "ConfigMap key which containers the MySQL password").Default("mysql.password").Envar(cmdenv.K8sConfigMapKeyPassword).StringVar(&c.params.Keys.Password)
	cmd.Flag("key-database", "ConfigMap key which containers the MySQL database").Default("mysql.database").Envar(cmdenv.K8sConfigMapKeyDatabase).StringVar(&c.params.Keys.Database)
	cmd.Flag("key-image", "ConfigMap key which containers the MySQL image").Default("mysql.docker.image").Envar(cmdenv.K8sConfigMapKeyImage).StringVar(&c.params.Keys.Image)
	cmd.Flag("role", "ServiceRole or IAM resource which grants access to the S3 bucket").Required().Envar(cmdenv.AWSIAMRole).StringVar(&c.params.Role)
	cmd.Flag("bucket", "Bucket to upload the file temporarily before CodeBuild runs").Required().Envar(cmdenv.AWSS3Bucket).StringVar(&c.params.Bucket)
	cmd.Flag("key-id", "AWS Credentials").Required().Envar(cmdenv.AWSAccessKeyID).StringVar(&c.params.AWS.KeyID)
	cmd.Flag("access-key", "AWS Credentials").Required().Envar(cmdenv.AWSSecretAccessKey).StringVar(&c.params.AWS.AccessKey)
	cmd.Flag("docker-username", "Username for the Docker Registry").Required().Envar(cmdenv.DockerUsername).StringVar(&c.params.Docker.Username)
	cmd.Flag("docker-password", "Password for the Docker Registry").Required().Envar(cmdenv.DockerPassword).StringVar(&c.params.Docker.Password)
}
