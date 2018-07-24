package cmd

import (
	"log"
	"os"

	"gopkg.in/alecthomas/kingpin.v2"

	"github.com/pkg/errors"
	"github.com/previousnext/mysql-toolkit/internal/dumper"
	"github.com/previousnext/mysql-toolkit/internal/envar"
)

type cmdDump struct {
	params dumper.DumpParams
	file   string
	config string
}

func (cmd *cmdDump) run(c *kingpin.ParseContext) error {
	cfg, err := dumper.Load(cmd.config)
	if err != nil {
		return errors.Wrap(err, "failed to load config")
	}

	cmd.params.Config = cfg
	cmd.params.Logger = log.New(os.Stderr, "", 0)
	cmd.params.SQLWriter = os.Stdout

	// Use stdout if no output file specified.
	if cmd.file != "" {
		cmd.params.Logger.Println("Opening file for writing:", cmd.file)

		writer, err := os.Create(cmd.file)
		if err != nil {
			return err
		}
		defer writer.Close()

		cmd.params.SQLWriter = writer
	}

	return dumper.Dump(cmd.params)
}

// Dump declares the "dump" subcommand.
func Dump(app *kingpin.CmdClause) {
	c := new(cmdDump)

	cmd := app.Command("dump", "Dump the database using a MySQL connection").Action(c.run)

	cmd.Flag("hostname", "Hostname for connecting to Mysql").Required().Envar(envar.MySQLHostname).StringVar(&c.params.Connection.Hostname)
	cmd.Flag("username", "Username for connecting to Mysql").Required().Envar(envar.MySQLUsername).StringVar(&c.params.Connection.Username)
	cmd.Flag("password", "Password for connecting to Mysql").Required().Envar(envar.MySQLPassword).StringVar(&c.params.Connection.Password)
	cmd.Flag("database", "Database for connecting to Mysql").Required().Envar(envar.MySQLDatabase).StringVar(&c.params.Connection.Database)
	cmd.Flag("protocol", "Protocol for connecting to Mysql").Default("tcp").Envar(envar.MySQLProtocol).StringVar(&c.params.Connection.Protocol)
	cmd.Flag("port", "Port for connecting to Mysql").Default("3306").Envar(envar.MySQLPort).StringVar(&c.params.Connection.Port)
	cmd.Flag("max-conn", "Maximum amount of open connections").Default("50").Envar(envar.MySQLMaxConn).IntVar(&c.params.Connection.MaxConn)
	cmd.Flag("config", "Policy for dumping the database").Default("config.yml").Envar(envar.MySQLConfig).StringVar(&c.config)
	cmd.Flag("file", "Location to save the dumped database. Omit to send to stdout.").Envar(envar.MySQLFile).StringVar(&c.file)
}
