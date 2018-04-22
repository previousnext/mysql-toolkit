package main

import (
	"os"

	"github.com/previousnext/mysql-toolkit/cmd"
	"gopkg.in/alecthomas/kingpin.v2"
)

func main() {
	app := kingpin.New("mysql-toolkit", "Toolkit for working with MySQL databases")

	cmd.Dump(app)
	cmd.CodeBuild(app)
	cmd.CodeBuildK8s(app)
	cmd.Version(app)

	kingpin.MustParse(app.Parse(os.Args[1:]))
}
