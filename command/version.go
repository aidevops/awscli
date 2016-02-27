// Package command -
package command

// Imports -
import (
	"bytes"
	"fmt"
	"github.com/mitchellh/cli"
)

// VersionCommand - is a Command implementation prints the version.
type VersionCommand struct {
	Revision          string
	Version           string
	VersionPrerelease string
	UI                cli.Ui
}

// Help -
func (c *VersionCommand) Help() string {
	return ""
}

// Run -
func (c *VersionCommand) Run(_ []string) int {
	var versionString bytes.Buffer
	fmt.Fprintf(&versionString, "awscli v%s", c.Version)
	if c.VersionPrerelease != "" {
		fmt.Fprintf(&versionString, ".%s", c.VersionPrerelease)

		if c.Revision != "" {
			fmt.Fprintf(&versionString, " (%s)", c.Revision)
		}
	}

	c.UI.Output(versionString.String())
	return 0
}

// Synopsis -
func (c *VersionCommand) Synopsis() string {
	return "Prints the awscli version"
}
