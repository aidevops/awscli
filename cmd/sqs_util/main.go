// Package main - sqs_util application
package main

// import - import our dependencies
import (
	"flag"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
)

// Unit - this application's name
const Unit = "sqs_util"

// verbose - control debug output
var verbose bool

// main - log us in...
func main() {
	var (
		account    string
		attributes string
		region     string
		queue      string
		message    string
		url        bool
		version    bool
	)

	var empty string
	flag.StringVar(&account, "account", "", "AWS account #. E.g. -account='1234556790123'")
	flag.StringVar(&region, "region", "us-east-1", "AWS region. E.g. -region=us-east-1")
	flag.StringVar(&queue, "queue", "", "vault-registration, consul-registration, serviceN-registration...")
	flag.StringVar(&message, "message", "", "-message 'hello world'")
	flag.StringVar(&attributes, "attributes", empty, "-attributes 'foo=bar,bar=foo,hello=world'")
	flag.BoolVar(&verbose, "verbose", false, "be more verbose.....")
	flag.BoolVar(&version, "version", false, "print version and exit")
	flag.BoolVar(&url, "url", false, "url the url of the queue and exit")
	flag.Parse()

	if version == true {
		fmt.Println(versionInfo())
		os.Exit(0)
	}

	debugf("[DEBUG]: using account: %s\n", account)
	if account == "" || len(account) < 12 {
		fmt.Printf("sqs_util: missing or invalid account length: -account='1234556790123', received: '%s'\n", account)
		os.Exit(255)
	}

	debugf("[DEBUG]: using queue name(s): %s\n", queue)
	if queue == "" || len(queue) < 3 {
		fmt.Printf("sqs_util: missing or invalid queue(s): -queue='some-fancy-queue..', received: '%s'\n", queue)
		os.Exit(254)
	}

	debugf("[DEBUG]: using region: %s\n", region)
	ok, err := Send(account, region, verbose, queue, message, ToMap(attributes), url)
	if !ok {
		fmt.Printf("[ERROR]: failed to send: %s", err)
		os.Exit(253)
	}
	// success!!!
	os.Exit(0)

}

// GetQueueURL - return the url of a queue by its name
func GetQueueURL(ses *session.Session, account, region, queue string) (string, error) {
	svc := sqs.New(ses, &aws.Config{Region: aws.String(region)})
	params := &sqs.GetQueueUrlInput{
		QueueName:              aws.String(queue), // Required
		QueueOwnerAWSAccountId: aws.String(account),
	}
	resp, err := svc.GetQueueUrl(params)
	if err != nil {
		// Print the error, cast err to awserr.Error to get the Code and
		// Message from an error.
		return "", fmt.Errorf("failed to get lookup queue by name '%s' %s\n", queue, err.Error())
	}
	return fmt.Sprintf("%s", *resp.QueueUrl), nil
}

// Send - send a messsage to aws sqs destination
func Send(account, region string, verbose bool, queue string, message string, attributes map[string]string, url bool) (ok bool, err error) {

	ses := session.New()
	queueURL, err := GetQueueURL(ses, account, region, queue)
	if err != nil {
		return false, fmt.Errorf("[ERROR] lookup queue url for queue '%s': %s\n", queue, err.Error())
	}

	debugf("[DEBUG]: found url: '%s' for queue '%s'\n", queueURL, queue)

	if url {
		fmt.Println(queueURL)
		os.Exit(0)
	}

	debugf("[DEBUG]: creating new session...\n")
	svc := sqs.New(ses, &aws.Config{Region: aws.String(region)})
	debugf("[DEBUG]: creating message(s) input...\n")

	params := &sqs.SendMessageInput{
		MessageBody:  aws.String(message),
		QueueUrl:     aws.String(queueURL),
		DelaySeconds: aws.Int64(1),
	}
	if len(attributes) > 0 {
		attrs := make(map[string]*sqs.MessageAttributeValue, len(attributes))
		for k, v := range attributes {
			attrs[k] = &sqs.MessageAttributeValue{DataType: aws.String("String"), StringValue: aws.String(v)}
		}
		params.MessageAttributes = attrs
	}
	resp, err := svc.SendMessage(params)
	debugf("[DEBUG]: response: %v\n", resp)

	if err != nil {
		return false, fmt.Errorf("Could not send message '%s' to queue '%s'@'%s': %s\n", message, queue, queueURL, err)
	}

	debugf("[DEBUG]: Successfully sent message(s) '%s'\n", message)
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
