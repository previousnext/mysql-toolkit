package acquia

import (
	"os"

	"github.com/fatih/color"
	"github.com/pkg/errors"
	"gopkg.in/alecthomas/kingpin.v2"

	cmdenv "github.com/previousnext/mysql-toolkit/cmd/env"
	"github.com/previousnext/mysql-toolkit/internal/acquia/backups"
)

type cmdDump struct {
	Username    string
	Password    string
	Site        string
	Environment string
	Database    string
	File        string
}

func (cmd *cmdDump) run(c *kingpin.ParseContext) error {
	color.Red("WARNING: Backups must be santized on the Acquia platform")

	tmp, err := os.Create(cmd.File)
	if err != nil {
		return errors.Wrap(err, "failed to create local tmp file")
	}
	defer tmp.Close()

	return backups.New(cmd.Username, cmd.Password, cmd.Site, cmd.Environment, cmd.Database).Download(tmp)
}

// Dump declares the "dump" subcommand.
func Dump(app *kingpin.CmdClause) {
	c := new(cmdDump)

	cmd := app.Command("dump", "Dump a database using Backups on the Acquia platform").Action(c.run)
	cmd.Flag("username", "Acquia uesrname used for authentication").Required().Envar(cmdenv.AcquiaUsername).StringVar(&c.Username)
	cmd.Flag("password", "Acquia password used for authentication").Required().Envar(cmdenv.AcquiaPassword).StringVar(&c.Password)
	cmd.Flag("site", "Acquia site to dump the database from eg. prod:example").Required().Envar(cmdenv.AcquiaSite).StringVar(&c.Site)
	cmd.Flag("environment", "Acquia environment to dump the database from eg. dev/test/prod").Required().Envar(cmdenv.AcquiaEnvironment).StringVar(&c.Environment)
	cmd.Flag("database", "Name of the database from the Acquia environment").Required().Envar(cmdenv.AcquiaDatabase).StringVar(&c.Database)
	cmd.Flag("file", "Location to save the dumped database").Required().Envar(cmdenv.MySQLFile).StringVar(&c.File)
}
