package dumper

import (
	"database/sql"
	"io"
	"log"
	"os"

	"github.com/hgfischer/mysqlsuperdump/dumper"
	"github.com/pkg/errors"
)

// DumpParams are passed to the Dump method.
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

// Dump the MySQL database.
func Dump(w io.Writer, params DumpParams) error {
	var logger = log.New(w, "", 0)

	logger.Println("Connecting to Mysql database:", params.Connection.Database)

	db, err := sql.Open("mysql", params.Connection.String())
	if err != nil {
		return errors.Wrap(err, "cannot connect to database")
	}
	defer db.Close()

	logger.Println("Opening file for writing:", params.File)

	f, err := os.Create(params.File)
	if err != nil {
		return err
	}
	defer f.Close()

	logger.Println("Setting maximum connection to:", params.Connection.MaxConn)

	db.SetMaxOpenConns(params.Connection.MaxConn)

	d := dumper.NewMySQLDumper(db, logger)

	cfg, err := Load(params.Config)
	if err != nil {
		return errors.Wrap(err, "failed to load config")
	}

	// Assign nodata tables.
	d.FilterMap =  make(map[string]string)
	for _, table := range cfg.NoData {
		d.FilterMap[table] = "nodata"
	}
	// Assign our sanitization rules to the dumper.
	d.SelectMap = cfg.Sanitize.Map()

	logger.Println("Starting to dump database:", params.Connection.Database)

	return d.Dump(f)
}
