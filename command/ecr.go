package command

import (
	"flag"
	"fmt"
	"strings"

	"github.com/mitchellh/cli"

	"github.com/johnt337/awscli"
	"github.com/johnt337/awscli/logger"
)

// ECRCommand -
type ECRCommand struct {
	UI cli.Ui
}

// Help -
func (c *ECRCommand) Help() string {
	helpText := `
Usage: awscli ecr [options] name
  
  ECR.....

Options:
  
  -verbose=true  Display additional information from 
                 behind the scenes.
`
	return strings.TrimSpace(helpText)
}

// Run -
func (c *ECRCommand) Run(args []string) int {
	var (
		account string
		region  string
		format  string
		level   string
		logfile string
		verbose bool
	)

	cmdFlags := flag.NewFlagSet("ecr", flag.ContinueOnError)
	cmdFlags.Usage = func() { c.UI.Output(c.Help()) }

	cmdFlags.StringVar(&account, "account", "", "AWS account #.")
	cmdFlags.StringVar(&region, "region", "", "AWS region.")
	cmdFlags.StringVar(&format, "format", "text", "Format response as either json or regular text.")
	cmdFlags.StringVar(&level, "level", "info", "logging level: error, warn, info, or debug")
	cmdFlags.StringVar(&logfile, "log", "/tmp/cloudconfig.log", "logfile path")
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

	ecr := args[0]

	c.UI.Output(fmt.Sprintf("Setting ecr to '%s'! Verbosity enabled: %#v",
		ecr, verbose))

	log := logger.NewCLILogger(level, logfile, "ecr", format, c.UI)

	awscli.ECRInfo(account)

	log.Flush()

	return 0
}

// Synopsis -
func (c *ECRCommand) Synopsis() string {
	return "ecr....."
}
