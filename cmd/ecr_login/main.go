// Package main - ecr_login application
package main

// import - import our dependencies
import (
	"encoding/base64"
	"flag"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ecr"
)

// Unit - this application's name
const Unit = "ecr_login"

// verbose - control debug output
var verbose bool

// main - log us in...
func main() {
	var (
		account string
		region  string
		version bool
	)

	flag.StringVar(&account, "account", "", "AWS account #. E.g. -account='1234556790123'")
	flag.StringVar(&region, "region", "us-east-1", "AWS region. E.g. -region=us-east-1")
	flag.BoolVar(&verbose, "verbose", false, "be more verbose.....")
	flag.BoolVar(&version, "version", false, "print version and exit")
	flag.Parse()

	if version == true {
		fmt.Println(versionInfo())
		os.Exit(0)
	}

	debugf("[DEBUG]: using account: %s\n", account)
	debugf("[DEBUG]: checking length: %d\n", len(account))
	if account == "" || len(account) < 12 {
		fmt.Printf("ecr_login: missing or invalid account length: -account='1234556790123', received: '%s'\n", account)
		os.Exit(255)
	}

	debugf("[DEBUG]: using region: %s", region)
	debugf("[DEBUG]: logging in...\n")
	token, err := Login(account, region, verbose)
	if err != nil {
		fmt.Printf("ecr_login: login error: %s\n", err)
		os.Exit(254)
	}

	debugf("[DEBUG]: decoding creds...\n")
	decoded, err := base64.StdEncoding.DecodeString(token)
	if err != nil {
		fmt.Printf("ecr_login: decode error: %s\n", err)
		os.Exit(253)
	}

	fmt.Println(string(decoded))
	os.Exit(0)

}

// Login - login to aws ecr registry
func Login(registryID, region string, verbose bool) (token string, err error) {

	debugf("[DEBUG]: creating new session...\n")
	svc := ecr.New(session.New(), &aws.Config{Region: aws.String(region)})

	debugf("[DEBUG]: creating auth token input...\n")
	params := &ecr.GetAuthorizationTokenInput{
		RegistryIds: []*string{
			aws.String(registryID), // Required
			// More values...
		},
	}

	debugf("[DEBUG]: fetching auth token...\n")
	resp, err := svc.GetAuthorizationToken(params)

	if err != nil {
		// Print the error, cast err to awserr.Error to get the Code and
		// Message from an error.
		fmt.Printf("ecr_login: %s\n", err.Error())
		return token, err
	}

	debugf("[DEBUG]: formatting and returning login token...\n")

	// Pretty-print the response data.
	fmt.Println(resp)
	token = fmt.Sprintf("%s", resp)
	return token, nil
}

// helper functions....

// debugf - print to stdout if verbose is enabled....
func debugf(format string, args ...interface{}) {
	if verbose == true {
		fmt.Printf(format, args...)
	}
}

// versionInfo - vendoring version info
func versionInfo() string {
	return fmt.Sprintf("%s v%s.%s (%s)", Unit, Version, VersionPrerelease, GitCommit)
}
