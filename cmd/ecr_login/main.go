package main

import (
	"encoding/base64"
	"flag"
	"fmt"
	"strings"

	"github.com/mitchellh/cli"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ecr"

	"github.com/johnt337/awscli/logger"
)

// ECRLogin -
type ECRLogin struct {
	UI cli.Ui
}

// Help -
func (c *ECRLogin) Help() string {
	helpText := `
Usage: awscli ecr [options] name
  
  ECR.....

Options:
  
  -verbose=true  Display additional information from 
                 behind the scenes.
`
	return strings.TrimSpace(helpText)
}

func main(args []string) {
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

	ecr := args[0]

	args = cmdFlags.Args()
	if len(args) < 1 {
		c.UI.Error("arguments must be specified.")
		c.UI.Error("")
		c.UI.Error(c.Help())
		return 1
	}

	log := logger.NewCLILogger(level, logfile, "ecr", format, c.UI)

	awscli.ECRInfo(account)
	token, _ := awscli.ECRLogin(account)
	decoded, err := base64.StdEncoding.DecodeString(token)
	if err != nil {
		fmt.Println("decode error:", err)
		return 255
	}
	fmt.Println(string(decoded))
	log.Flush()

	return 0

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

// login - login to aws ecr registry
func login(registryID string) (token string, err error) {
	svc := ecr.New(session.New())

	params := &ecr.GetAuthorizationTokenInput{
		RegistryIds: []*string{
			aws.String(registryID), // Required
			// More values...
		},
	}
	resp, err := svc.GetAuthorizationToken(params)

	if err != nil {
		// Print the error, cast err to awserr.Error to get the Code and
		// Message from an error.
		fmt.Println(err.Error())
		return token, err
	}

	// Pretty-print the response data.
	fmt.Println(resp)
	token = fmt.Sprintf("%s", resp)
	return token, nil
}
