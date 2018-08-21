package acquia

import (
	"context"

	"github.com/operator-framework/operator-sdk/pkg/sdk"
	"gopkg.in/alecthomas/kingpin.v2"
	corev1 "k8s.io/api/core/v1"

	"github.com/previousnext/mysql-toolkit/internal/envar"
	acquiasnapshot "github.com/previousnext/mysql-toolkit/internal/operator/handler/acquia/snapshot"
)

type cmdSnapshot struct {
	Watch     string
	Namespace string
	Secret    string
	Resync    int
	JobImage  string
	JobCPU    string
	JobMemory string
}

func (cmd *cmdSnapshot) run(c *kingpin.ParseContext) error {
	sdk.Watch("mtk.skpr.io/v1alpha1", "AcquiaSnapshot", cmd.Watch, cmd.Resync)
	sdk.Handle(acquiasnapshot.NewHandler(cmd.Namespace, cmd.Secret, cmd.JobImage, cmd.JobCPU, cmd.JobMemory))
	sdk.Run(context.TODO())
	return nil
}

// Snapshot declares the "operator-snapshot" subcommand.
func Snapshot(app *kingpin.CmdClause) {
	c := new(cmdSnapshot)

	cmd := app.Command("snapshot", "Operator for running a single 'Snapshot' snapshot").Action(c.run)
	cmd.Flag("namespace-watch", "Namespace to watch for new AcquiaSnapshot objects").Default(corev1.NamespaceAll).Envar(envar.OperatorNamespaceWatch).StringVar(&c.Watch)
	cmd.Flag("namespace", "Namespace to execute AcquiaSnapshot Jobs").Required().Envar(envar.OperatorNamespace).StringVar(&c.Namespace)
	cmd.Flag("secret", "Secret for loading Operator configuration").Required().Envar(envar.OperatorSecret).StringVar(&c.Secret)
	cmd.Flag("resync", "How often to resync all the AcquiaSnapshot objects").Default("30").Envar(envar.OperatorResync).IntVar(&c.Resync)
	cmd.Flag("job-image", "Image to use when executing a Job").Default("previousnext/mtk:latest").Envar(envar.OperatorJobImage).StringVar(&c.JobImage)
	cmd.Flag("job-cpu", "CPU allocated when executing of a Job").Default("250m").Envar(envar.OperatorJobCPU).StringVar(&c.JobCPU)
	cmd.Flag("job-memory", "Memory allocated when executing of a Job").Default("512Mi").Envar(envar.OperatorJobMemory).StringVar(&c.JobMemory)
}
