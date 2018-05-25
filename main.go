package main

import (
	"os"

	_ "github.com/go-sql-driver/mysql"
	"gopkg.in/alecthomas/kingpin.v2"

	"github.com/previousnext/mysql-toolkit/cmd"
	cmdbuild "github.com/previousnext/mysql-toolkit/cmd/build"
	cmddb "github.com/previousnext/mysql-toolkit/cmd/db"
)

func main() {
	app := kingpin.New("mysql-toolkit", "Toolkit for working with MySQL databases")

	cmd.Version(app)

	db := app.Command("db", "Dump the database")
	cmddb.Dump(db)
	cmddb.Operator(db)

	build := app.Command("build", "Build the database image")
	cmdbuild.CodeBuild(build)

	kingpin.MustParse(app.Parse(os.Args[1:]))
}
