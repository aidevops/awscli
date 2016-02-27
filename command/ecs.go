package command

import (
	"flag"
	"fmt"
	"strings"

	"github.com/mitchellh/cli"
)

// ECSCommand -
type ECSCommand struct {
	UI cli.Ui
}

// Help -
func (c *ECSCommand) Help() string {
	helpText := `
Usage: awscli ecs [options] name
  
  ECS.....

Options:
  
  -verbose=true  Display additional information from 
                 behind the scenes.
`
	return strings.TrimSpace(helpText)
}

// Run -
func (c *ECSCommand) Run(args []string) int {
	var verbose bool

	cmdFlags := flag.NewFlagSet("ecs", flag.ContinueOnError)
	cmdFlags.Usage = func() { c.UI.Output(c.Help()) }
	cmdFlags.BoolVar(&verbose, "verbose", false, "verbose")

	if err := cmdFlags.Parse(args); err != nil {
		return 1
	}

	args = cmdFlags.Args()
	if len(args) < 1 {
		c.UI.Error("arguments must be specified.")
		c.UI.Error("")
		c.UI.Error(c.Help())
		return 1
	} else if len(args) > 1 {
		c.UI.Error("Too many command line arguments.")
		c.UI.Error("")
		c.UI.Error(c.Help())
		return 1
	}

	ecs := args[0]

	c.UI.Output(fmt.Sprintf("Setting ecs to '%s'! Verbosity enabled: %#v",
		ecs, verbose))

	return 0
}

// Synopsis -
func (c *ECSCommand) Synopsis() string {
	return "ecs....."
}
