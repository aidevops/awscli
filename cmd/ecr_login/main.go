package main

import (
	"encoding/base64"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ecr"
)

// Help -
func Help() string {
	helpText := `
Usage: awscli ecr [options] name
  
  ECR.....

Options:
  
  -verbose=true  Display additional information from 
                 behind the scenes.
`
	return strings.TrimSpace(helpText)
}

func main() {
	var (
		account string
		region  string
		verbose bool
	)

	if verbose == true {
		fmt.Println("[DEBUG]: starting up...")
	}

	cmdFlags := flag.NewFlagSet("ecr_login", flag.ContinueOnError)
	cmdFlags.Usage = func() { Help() }
	cmdFlags.StringVar(&account, "account", "", "AWS account #.")
	cmdFlags.StringVar(&region, "region", "us-east-1", "AWS region.")
	cmdFlags.BoolVar(&verbose, "verbose", false, "verbose")

	if verbose == true {
		fmt.Println("[DEBUG]: logging in...")
	}
	token, err := Login(account, region, verbose)
	if err != nil {
		fmt.Errorf("ecr_login: login error: %s\n", err)
	}

	if verbose == true {
		fmt.Println("[DEBUG]: decoding creds...")
	}
	decoded, err := base64.StdEncoding.DecodeString(token)
	if err != nil {
		fmt.Errorf("ecr_login: decode error: %s\n", err)
		os.Exit(255)
	}

	fmt.Println(string(decoded))
	os.Exit(0)

}

// login - login to aws ecr registry
func Login(registryID, region string, verbose bool) (token string, err error) {

	if verbose {
		fmt.Println("[DEBUG]: creating new session...")
	}
	svc := ecr.New(session.New(), &aws.Config{Region: aws.String(region)})

	if verbose {
		fmt.Println("[DEBUG]: creating auth token input...")
	}
	params := &ecr.GetAuthorizationTokenInput{
		RegistryIds: []*string{
			aws.String(registryID), // Required
			// More values...
		},
	}

	if verbose {
		fmt.Println("[DEBUG]: fetching auth token...")
	}
	resp, err := svc.GetAuthorizationToken(params)

	if err != nil {
		// Print the error, cast err to awserr.Error to get the Code and
		// Message from an error.
		fmt.Errorf("ecr_login: %s\n", err.Error())
		return token, err
	}

	if verbose {
		fmt.Println("[DEBUG]: formatting and returning login token...")
	}

	// Pretty-print the response data.
	fmt.Println(resp)
	token = fmt.Sprintf("%s", resp)
	return token, nil
}
