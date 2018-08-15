package dumper

import (
	"database/sql"
	"io"
	"log"

	"github.com/hgfischer/mysqlsuperdump/dumper"
	"github.com/pkg/errors"

	"github.com/previousnext/mysql-toolkit/internal/dumper/tableglob"
)

// DumpParams store arguments provided by CLI.
type DumpParams struct {
	Connection Connection
	Config     Config
	Logger     *log.Logger
	SQLWriter  io.Writer
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

// Dump the MySQL database.
func Dump(args DumpParams) error {
	logger := args.Logger

	logger.Println("Connecting to Mysql database:", args.Connection.Database)

	conn := args.Connection.String()

	db, err := sql.Open("mysql", conn)
	if err != nil {
		return errors.Wrap(err, "cannot connect to database")
	}
	defer db.Close()

	logger.Println("Setting maximum connection to:", args.Connection.MaxConn)

	db.SetMaxOpenConns(args.Connection.MaxConn)

	d := dumper.NewMySQLDumper(db, logger)

	// Get a list of tables to nodata, passed through a globber.
	nodata, err := tableglob.Show(conn, args.Config.NoData)
	if err != nil {
		return err
	}

	// Get a list of tables to ignore, passed through a globber.
	ignore, err := tableglob.Show(conn, args.Config.Ignore)
	if err != nil {
		return err
	}

	// Assign nodata tables.
	d.FilterMap = make(map[string]string)
	for _, table := range nodata {
		d.FilterMap[table] = "nodata"
	}

	// Assign ignore tables.
	for _, table := range ignore {
		d.FilterMap[table] = "ignore"
	}

	// Assign our sanitization rules to the dumper.
	d.SelectMap = args.Config.Sanitize.Map()

	logger.Println("Starting to dump database:", args.Connection.Database)

	return d.Dump(args.SQLWriter)
}
