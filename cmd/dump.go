package cmd

import (
	"os"
	"log"

	"gopkg.in/alecthomas/kingpin.v2"

	cmdenv "github.com/previousnext/mysql-toolkit/cmd/env"
	"github.com/previousnext/mysql-toolkit/internal/dumper"
	"github.com/pkg/errors"
)

type cmdDump struct {
	params dumper.DumpParams
}

func (cmd *cmdDump) run(c *kingpin.ParseContext) error {
	cfg, err := dumper.Load(cmd.params.Config)
	if err != nil {
		return errors.Wrap(err, "failed to load config")
	}

	var logger = log.New(os.Stdout, "", 0)

	logger.Println("Opening file for writing:", cmd.params.File)
	f, err := os.Create(cmd.params.File)
	if err != nil {
		return err
	}
	defer f.Close()

	return dumper.Dump(dumper.DumpArgs{
		Logger:     logger,
		SQLWriter:  f,
		Config:     cfg,
		Connection: cmd.params.Connection,
	})
}

// Dump declares the "dump" subcommand.
func Dump(app *kingpin.Application) {
	c := new(cmdDump)

	cmd := app.Command("dump", "Dump the database").Action(c.run)

	cmd.Flag("hostname", "Hostname for connecting to Mysql").Required().Envar(cmdenv.MySQLHostname).StringVar(&c.params.Connection.Hostname)
	cmd.Flag("username", "Username for connecting to Mysql").Required().Envar(cmdenv.MySQLUsername).StringVar(&c.params.Connection.Username)
	cmd.Flag("password", "Password for connecting to Mysql").Required().Envar(cmdenv.MySQLPassword).StringVar(&c.params.Connection.Password)
	cmd.Flag("database", "Database for connecting to Mysql").Required().Envar(cmdenv.MySQLDatabase).StringVar(&c.params.Connection.Database)
	cmd.Flag("protocol", "Protocol for connecting to Mysql").Default("tcp").Envar(cmdenv.MySQLProtocol).StringVar(&c.params.Connection.Protocol)
	cmd.Flag("port", "Port for connecting to Mysql").Default("3306").Envar(cmdenv.MySQLPort).StringVar(&c.params.Connection.Port)
	cmd.Flag("max-conn", "Maximum amount of open connections").Default("50").Envar(cmdenv.MySQLMaxConn).IntVar(&c.params.Connection.MaxConn)
	cmd.Flag("config", "Policy for dumping the database").Default("config.yml").Envar(cmdenv.MySQLConfig).StringVar(&c.params.Config)
	cmd.Flag("file", "Location to save the dumped database").Required().Envar(cmdenv.MySQLFile).StringVar(&c.params.File)
}
