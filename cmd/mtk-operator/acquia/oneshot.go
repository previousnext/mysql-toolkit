package acquia

import (
	"context"

	"github.com/operator-framework/operator-sdk/pkg/sdk"
	"github.com/sirupsen/logrus"
	"gopkg.in/alecthomas/kingpin.v2"
	corev1 "k8s.io/api/core/v1"

	"github.com/previousnext/mysql-toolkit/internal/envar"
	acquiaoneshot "github.com/previousnext/mysql-toolkit/internal/operator/handler/acquia/oneshot"
)

type cmdOneShot struct {
	Namespace string
	Resync    int
}

func (cmd *cmdOneShot) run(c *kingpin.ParseContext) error {
	var (
		resource = "mtk.skpr.io/v1alpha1"
		kind     = "AcquiaOneShot"
	)

	logrus.Infof("Watching %s, %s, %s, %d", resource, kind, cmd.Namespace, cmd.Resync)

	sdk.Watch(resource, kind, cmd.Namespace, cmd.Resync)
	sdk.Handle(acquiaoneshot.NewHandler())
	sdk.Run(context.TODO())

	return nil
}

// OneShot declares the "operator-oneshot" subcommand.
func OneShot(app *kingpin.CmdClause) {
	c := new(cmdOneShot)

	cmd := app.Command("oneshot", "Operator for running a single 'OneShot' snapshot").Action(c.run)
	cmd.Flag("namespace", "Namespace which to watch for new AcquiaOneShot objects").Default(corev1.NamespaceAll).Envar(envar.OperatorNamespace).StringVar(&c.Namespace)
	cmd.Flag("resync", "How often to resync all the AcquiaOneShot objects").Default("30").Envar(envar.OperatorResync).IntVar(&c.Resync)
}
