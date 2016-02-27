package command

import (
	"flag"
	"fmt"
	"strings"

	"github.com/mitchellh/cli"
)

// EC2Command -
type EC2Command struct {
	UI cli.Ui
}

// Help -
func (c *EC2Command) Help() string {
	helpText := `
Usage: awscli ec2 [options] name
  
  EC2.....

Options:
  
  -verbose=true  Display additional information from 
                 behind the scenes.
`
	return strings.TrimSpace(helpText)
}

// Run -
func (c *EC2Command) Run(args []string) int {
	var verbose bool

	cmdFlags := flag.NewFlagSet("ec2", flag.ContinueOnError)
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

	ec2 := args[0]

	c.UI.Output(fmt.Sprintf("Setting ec2 to '%s'! Verbosity enabled: %#v",
		ec2, verbose))

	return 0
}

// Synopsis -
func (c *EC2Command) Synopsis() string {
	return "ec2....."
}
