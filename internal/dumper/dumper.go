package dumper

import (
	"database/sql"
	"io"
	"log"

	"github.com/hgfischer/mysqlsuperdump/dumper"
	"github.com/pkg/errors"
)

// DumpParams store arguments provided by CLI.
type DumpParams struct {
	Connection Connection
	Config     string
	File       string
}

// Connection details for MySQL.
type Connection struct {
	Hostname string
	Username string
	Password string
	Database string
	Protocol string
	Port     string
	MaxConn  int
}

// DumpArgs are passed to Dump method.
type DumpArgs struct {
	Logger       *log.Logger
	SQLWriter    io.Writer
	Config       Config
	Connection   Connection
}

// Dump the MySQL database.
func Dump(args DumpArgs) error {
	logger := args.Logger

	logger.Println("Connecting to Mysql database:", args.Connection.Database)

	db, err := sql.Open("mysql", args.Connection.String())
	if err != nil {
		return errors.Wrap(err, "cannot connect to database")
	}
	defer db.Close()



	logger.Println("Setting maximum connection to:", args.Connection.MaxConn)

	db.SetMaxOpenConns(args.Connection.MaxConn)

	d := dumper.NewMySQLDumper(db, logger)

	// Assign nodata tables.
	d.FilterMap = make(map[string]string)
	for _, table := range args.Config.NoData {
		d.FilterMap[table] = "nodata"
	}

	// Assign ignore tables.
	for _, table := range args.Config.Ignore {
		d.FilterMap[table] = "ignore"
	}

	// Assign our sanitization rules to the dumper.
	d.SelectMap = args.Config.Sanitize.Map()

	logger.Println("Starting to dump database:", args.Connection.Database)

	return d.Dump(args.SQLWriter)
}
