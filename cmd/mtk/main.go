package main

import (
	"os"

	_ "github.com/go-sql-driver/mysql"
	"gopkg.in/alecthomas/kingpin.v2"

	cmdacquia "github.com/previousnext/mysql-toolkit/cmd/mtk/acquia"
	cmdbuild "github.com/previousnext/mysql-toolkit/cmd/mtk/build"
	cmddb "github.com/previousnext/mysql-toolkit/cmd/mtk/db"
	"github.com/previousnext/mysql-toolkit/internal/version"
)

func main() {
	app := kingpin.New("mtk", "MySQL Toolkit: utility for working with MySQL databases")

	version.Command(app)

	acquia := app.Command("acquia", "Acquia Platform tools")
	cmdacquia.Dump(acquia)

	db := app.Command("db", "MySQL database tools")
	cmddb.Dump(db)

	build := app.Command("build", "Builders for creating MySQL images")
	cmdbuild.AWS(build)
	cmdbuild.Docker(build)

	kingpin.MustParse(app.Parse(os.Args[1:]))
}
