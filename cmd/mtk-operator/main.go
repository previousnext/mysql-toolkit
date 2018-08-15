package main

import (
	"os"

	_ "github.com/go-sql-driver/mysql"
	"gopkg.in/alecthomas/kingpin.v2"

	cmdacquia "github.com/previousnext/mysql-toolkit/cmd/mtk-operator/acquia"
	"github.com/previousnext/mysql-toolkit/internal/version"
)

func main() {
	app := kingpin.New("mtk-operator", "MySQL Toolkit Operator: Kubernetes operators for working with MySQL databases")

	version.Command(app)

	acquia := app.Command("acquia", "Acquia Platform tools")
	cmdacquia.Snapshot(acquia)
	cmdacquia.SnapshotScheduled(acquia)

	kingpin.MustParse(app.Parse(os.Args[1:]))
}
