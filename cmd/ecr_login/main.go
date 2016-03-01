// Package main - ecr_login application
package main

// import - import our dependencies
import (
	"bufio"
	"encoding/base64"
	"flag"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
	"time"

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
		login   bool
	)

	flag.StringVar(&account, "account", "", "AWS account #. E.g. -account='1234556790123'")
	flag.StringVar(&region, "region", "us-east-1", "AWS region. E.g. -region=us-east-1")
	flag.BoolVar(&verbose, "verbose", false, "be more verbose.....")
	flag.BoolVar(&version, "version", false, "print version and exit")
	flag.BoolVar(&login, "login", false, "docker login on your behalf, otherwise return login string")
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
	debugf("[DEBUG]: generating login credentials...\n")
	token, endpoint, expires, err := Login(account, region, verbose, login)

	if err != nil {
		fmt.Printf("[ERROR]: generating login credentials: %s\n", err)
		os.Exit(254)
	}

	debugf("[DEBUG]: credentials valid until: %s...\n", expires.String())
	debugf("[DEBUG]: decoding creds...\n")
	decoded, err := base64.StdEncoding.DecodeString(*token)
	if err != nil {
		fmt.Printf("[ERROR]: decode error: %s\n", err)
		os.Exit(253)
	}

	creds := strings.Split(string(decoded), ":")

	debugf("[DEBUG]: creds length: %d\n", len(creds))
	debugf("[DEBUG]: generating login command\n")
	args := []string{"login", "-u", creds[0], "-p", creds[1], "-e", "none", *endpoint}

	if login == true {
		debugf("[DEBUG]: executing command: 'docker %s'\n", strings.Join(args, " "))
		cmd := exec.Command("docker", args...)
		stdout, err := cmd.StdoutPipe()
		if err != nil {
			fmt.Printf("[ERROR]: failed to open stdout: %s\n", err)
			os.Exit(252)
		}
		stderr, err := cmd.StderrPipe()
		if err != nil {
			fmt.Printf("[ERROR]: failed to open stderr: %s\n", err)
			os.Exit(251)
		}

		// start the command after having set up the pipes
		if err = cmd.Start(); err != nil {
			fmt.Printf("[ERROR]: failed to start command: %s\n", err)
			os.Exit(250)
		}

		// collect both pipes together
		multi := io.MultiReader(stdout, stderr)
		// read command's stdout & stderr line by line
		in := bufio.NewScanner(multi)

		for in.Scan() {
			line := in.Text()
			fmt.Println(line)
		}

		err = cmd.Wait()
		if err != nil {
			fmt.Printf("[ERROR]: failed while waiting for command to complete: %s\n", err)
			os.Exit(249)
		}
	} else {
		fmt.Print("docker ")
		fmt.Println(strings.Join(args, " "))
	}

	// success!!!
	os.Exit(0)

}

// Login - login to aws ecr registry
func Login(registryID, region string, verbose, login bool) (token, endpoint *string, expires *time.Time, err error) {

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
		fmt.Printf("[ERROR]: %s\n", err.Error())
		return token, endpoint, expires, err
	}

	debugf("[DEBUG]: formatting and returning login token...\n")

	// Pretty-print the response data.
	debugf("[DEBUG]: raw aws response: %s\n", resp)

	token = resp.AuthorizationData[0].AuthorizationToken
	endpoint = resp.AuthorizationData[0].ProxyEndpoint
	expires = resp.AuthorizationData[0].ExpiresAt
	return token, endpoint, expires, nil
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
