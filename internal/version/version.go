package version

import (
	"fmt"
	"os"
	"runtime"

	"gopkg.in/alecthomas/kingpin.v2"

	"github.com/previousnext/gopher/pkg/version"
)

var (
	// GitVersion overridden at build time by:
	//   -ldflags='-X github.com/previousnext/gopher/cmd.GitVersion=$(git describe --tags --always)'
	GitVersion string
	// GitCommit overridden at build time by:
	//   -ldflags='-X github.com/previousnext/gopher/cmd.GitCommit=$(git rev-list -1 HEAD)'
	GitCommit string
)

type cmdVersion struct {
	BuildDate    string
	BuildVersion string
	GOARCH       string
	GOOS         string
}

func (cmd *cmdVersion) run(c *kingpin.ParseContext) error {
	return version.Print(os.Stdin, version.PrintParams{
		Version: GitVersion,
		Commit:  GitCommit,
		OS:      runtime.GOOS,
		Arch:    runtime.GOARCH,
	})
}

// Command declares the "version" sub command.
func Command(app *kingpin.Application) {
	cmd := new(cmdVersion)
	app.Command("version", fmt.Sprintf("Prints %s version", app.Name)).Action(cmd.run)
}
