package main

import (
	"os"

	_ "github.com/go-sql-driver/mysql"
	"gopkg.in/alecthomas/kingpin.v2"

	"github.com/previousnext/mysql-toolkit/cmd"
	cmdacquia "github.com/previousnext/mysql-toolkit/cmd/acquia"
	cmdbuild "github.com/previousnext/mysql-toolkit/cmd/build"
	cmddb "github.com/previousnext/mysql-toolkit/cmd/db"
)

func main() {
	app := kingpin.New("mtk", "MySQL Toolkit: utility for working with MySQL databases")

	cmd.Version(app)

	acquia := app.Command("acquia", "Acquia Platform tools")
	cmdacquia.Dump(acquia)

	db := app.Command("db", "MySQL database tools")
	cmddb.Dump(db)
	cmddb.Operator(db)

	build := app.Command("build", "Builders for creating MySQL images")
	cmdbuild.AWS(build)

	kingpin.MustParse(app.Parse(os.Args[1:]))
}
