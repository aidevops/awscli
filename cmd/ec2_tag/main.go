// Package main - ecr_login application
package main

// import - import our dependencies
import (
	// "bufio"
	// "encoding/base64"
	"flag"
	"fmt"
	// "io"
	"os"
	// "os/exec"
	"regexp"
	"strconv"
	"strings"
	// "time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
)

// Unit - this application's name
const Unit = "ec2_tag"

// verbose - control debug output
var verbose bool

// for handling tags
var tag map[string]string

// main - log us in...
func main() {
	var (
		account   string
		region    string
		resources string
		tags      string
		version   bool
	)

	var empty string
	flag.StringVar(&account, "account", "", "AWS account #. E.g. -account='1234556790123'")
	flag.StringVar(&region, "region", "us-east-1", "AWS region. E.g. -region=us-east-1")
	flag.BoolVar(&verbose, "verbose", false, "be more verbose.....")
	flag.BoolVar(&version, "version", false, "print version and exit")
	flag.StringVar(&resources, "resources", empty, "-resources 'one two three four five'")
	flag.StringVar(&tags, "tags", empty, "-tags 'foo=bar,bar=foo,hello=world'")
	flag.Parse()

	if version == true {
		fmt.Println(versionInfo())
		os.Exit(0)
	}

	debugf("[DEBUG]: using account: %s\n", account)
	if account == "" || len(account) < 12 {
		fmt.Printf("ecr_login: missing or invalid account length: -account='1234556790123', received: '%s'\n", account)
		os.Exit(255)
	}

	debugf("[DEBUG]: using resource(s): %s\n", resources)
	if resources == "" || len(resources) < 10 {
		fmt.Printf("ecr_login: missing or invalid resource(s): -resources='i-86424106 i-864241.. i-864242..', received: '%s'\n", resources)
		os.Exit(255)
	}

	debugf("[DEBUG]: using region: %s\n", region)

	t := ToMap(tags)
	debugf("[DEBUG]: raw input: %s\n", tags)
	for k, v := range t {
		debugf("[DEBUG]: mapped: Key=%s,Value=%s\n", k, v)
	}
	ok, err := Tag(account, region, verbose, ToSlice(resources), t)
	if !ok {
		fmt.Printf("[ERROR]: failed to tag: %s", err)
		os.Exit(254)
	}
	// success!!!
	os.Exit(0)

}

// Login - login to aws ecr registry
func Tag(account, region string, verbose bool, resources []string, tags map[string]string) (ok bool, err error) {

	debugf("[DEBUG]: creating new session...\n")
	svc := ec2.New(session.New(), &aws.Config{Region: aws.String(region)})

	debugf("[DEBUG]: creating tag(s) input...\n")
	debugf("[DEBUG]: total tag pair(s): %d\n", len(tags))
	counter := 0
	ec2Tags := make([]*ec2.Tag, len(tags))
	for key, value := range tags {
		debugf("[DEBUG]: processing tag #%d: '%s'='%s'\n", counter, key, value)
		ec2Tags[counter] = &ec2.Tag{
			Key:   aws.String(key),
			Value: aws.String(value),
		}
		counter++
	}

	ec2Resources := make([]*string, len(resources))
	for pos, resource := range resources {
		debugf("[DEBUG]: tagging resource: '%s'\n", resource)
		ec2Resources[pos] = &resource
	}

	resp, err := svc.CreateTags(&ec2.CreateTagsInput{
		Resources: ec2Resources,
		Tags:      ec2Tags,
	})

	debugf("[DEBUG]: response: %v\n", resp)

	if err != nil {
		return false, fmt.Errorf("Could not create tags for instance(s): '%s': %s\n", strings.Join(resources, " "), err)
	}

	debugf("Successfully tagged instance(s) '%s'\n", strings.Join(resources, " "))
	return true, nil
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

// ToMap - Convert options into a go map
func ToMap(data string) map[string]string {
	opts := make(map[string]string)
	if data == "" {
		return opts
	}

	// var sanitized string
	// var err error
	sanitized, err := strconv.Unquote(data)
	if err != nil {
		// fmt.Printf("failed to strip quotes: '%s'\n", err)
		sanitized = data
	}

	re1, err := regexp.Compile(",")
	if err != nil {
		return opts
	}
	pairs := re1.Split(sanitized, -1)
	for _, field := range pairs {
		re2, err := regexp.Compile("=")
		if err != nil {
			return opts
		}
		pair := re2.Split(field, 2)
		key := pair[0]
		var val string
		if len(pair) == 2 {
			cleaned, err := strconv.Unquote(pair[1])
			if err != nil {
				val = pair[1]
			} else {
				val = cleaned
			}
		} else {
			val = ""
		}

		opts[key] = val
	}
	return opts
}

// ToSlice - return a string of space delimited arguments as a []string slice
func ToSlice(data string) (slice []string) {
	if data == "" {
		return slice
	}

	list := strings.Fields(data)
	slice = make([]string, len(list))
	for pos, field := range list {
		slice[pos] = field
	}
	return slice
}
