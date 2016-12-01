// Package main - sqs_util application
package main

// import - import our dependencies
import (
	"encoding/json"
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

// VisibilityTimeout - http://docs.aws.amazon.com/AWSSimpleQueueService/latest/SQSDeveloperGuide/AboutVT.html
const VisibilityTimeout = 1

// WaitTimeSeconds - The duration (in seconds) for which the call will wait for a message to arrive in the queue before returning.
const WaitTimeSeconds = 10

// verbose - control debug output
var verbose bool

// main - log us in...
func main() {
	var (
		account    string
		attributes string
		build      bool
		count      int64
		region     string
		queue      string
		message    string
		send       bool
		recv       bool
		url        bool
		version    bool
	)

	var empty string
	flag.StringVar(&account, "account", "", "AWS account #. E.g. -account='1234556790123'")
	flag.StringVar(&attributes, "attributes", empty, "-attributes 'foo=bar,bar=foo,hello=world'")
	flag.BoolVar(&build, "build", false, "build the url instead of looking it up against aws (less permission required)")
	flag.Int64Var(&count, "count", 1, "number of messages to retrieve from queue")
	flag.StringVar(&message, "message", "", "-message 'hello world'")
	flag.StringVar(&queue, "queue", "", "vault-registration, consul-registration, serviceN-registration...")
	flag.StringVar(&region, "region", "us-east-1", "AWS region. E.g. -region=us-east-1")
	flag.BoolVar(&send, "send", false, "send message")
	flag.BoolVar(&recv, "recv", false, "receive messages")
	flag.BoolVar(&verbose, "verbose", false, "be more verbose.....")
	flag.BoolVar(&version, "version", false, "print version and exit")
	flag.BoolVar(&url, "url", false, "lookup the url for -queue='...' and exit")
	flag.Parse()

	if version == true {
		fmt.Println(versionInfo())
		os.Exit(0)
	}

	if !send && !recv {
		fmt.Println("sqs_util: you need to specify either -send or -recv")
		os.Exit(1)
	}

	if send && recv {
		fmt.Println("sqs_util: send and recv are mutually exclusive")
		os.Exit(1)
	}

	debugf("[DEBUG]: using count: %d\n", count)
	if count < 0 || count > 10 {
		fmt.Printf("sqs_util: invalid count valid values 1 - 10\n", account)
		os.Exit(255)
	}

	debugf("[DEBUG]: using account: %s\n", account)
	if account == "" || len(account) < 12 {
		fmt.Printf("sqs_util: missing or invalid account length: -account='1234556790123', received: '%s'\n", account)
		os.Exit(254)
	}

	debugf("[DEBUG]: using queue name(s): %s\n", queue)
	if queue == "" || len(queue) < 3 {
		fmt.Printf("sqs_util: missing or invalid queue(s): -queue='some-fancy-queue..', received: '%s'\n", queue)
		os.Exit(253)
	}

	debugf("[DEBUG]: using region: %s\n", region)

	var ok bool
	var err error
	if send {
		ok, err = Send(account, region, verbose, queue, message, ToMap(attributes), url, build)
	}

	if recv {
		ok, err = Receive(account, region, verbose, queue, message, url, build, count)
	}

	if !ok {
		fmt.Printf("[ERROR]: failed while processing request: %s", err)
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
		return "", fmt.Errorf("failed to lookup queue by name '%s' %s", queue, err.Error())
	}
	return fmt.Sprintf("%s", *resp.QueueUrl), nil
}

// BuildQueueURL - Builds the url based on provided input instead of querying AWS sqs.
func BuildQueueURL(account, region, queue string) string {
	return fmt.Sprintf("https://sqs.%s.amazonaws.com/%s/%s", region, account, queue)
}

// Send - send a messsage to aws sqs destination
func Send(account, region string, verbose bool, queue string, message string, attributes map[string]string, url, build bool) (ok bool, err error) {

	var queueURL string
	ses := session.New()

	if build {
		queueURL = BuildQueueURL(account, region, queue)
	} else {
		queueURL, err = GetQueueURL(ses, account, region, queue)
		if err != nil {
			return false, fmt.Errorf("[ERROR] lookup queue url for queue '%s': %s", queue, err.Error())
		}

		debugf("[DEBUG]: found url: '%s' for queue '%s'\n", queueURL, queue)
	}

	if url {
		fmt.Println(queueURL)
		os.Exit(0)
	}

	debugf("[DEBUG]: creating new session...\n")
	svc := sqs.New(ses, &aws.Config{Region: aws.String(region)})
	debugf("[DEBUG]: creating send message(s) input...\n")

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
		return false, fmt.Errorf("Could not send message '%s' to queue '%s'@'%s': %s", message, queue, queueURL, err)
	}

	debugf("[DEBUG]: Successfully sent message(s) '%s'\n", message)
	return true, nil
}

// Receive - receive messsages from aws sqs destination
func Receive(account, region string, verbose bool, queue string, message string, url, build bool, count int64) (ok bool, err error) {
	var queueURL string
	ses := session.New()

	if build {
		queueURL = BuildQueueURL(account, region, queue)
	} else {
		queueURL, err = GetQueueURL(ses, account, region, queue)
		if err != nil {
			return false, fmt.Errorf("[ERROR] lookup queue url for queue '%s': %s", queue, err.Error())
		}

		debugf("[DEBUG]: found url: '%s' for queue '%s'\n", queueURL, queue)
	}

	if url {
		fmt.Println(queueURL)
		os.Exit(0)
	}

	debugf("[DEBUG]: creating new session...\n")
	svc := sqs.New(ses, &aws.Config{Region: aws.String(region)})
	debugf("[DEBUG]: creating receive message(s) input...\n")

	params := &sqs.ReceiveMessageInput{
		QueueUrl:            aws.String(queueURL),
		MaxNumberOfMessages: aws.Int64(count),
		MessageAttributeNames: []*string{
			aws.String("id"),         // Required
			aws.String("node"),       // Required
			aws.String("role"),       // Required
			aws.String("instance"),   // Required
			aws.String("registered"), // Required
			// More values...
		},
		VisibilityTimeout: aws.Int64(VisibilityTimeout),
		WaitTimeSeconds:   aws.Int64(WaitTimeSeconds),
	}
	resp, err := svc.ReceiveMessage(params)

	total := len(resp.Messages)
	for pos, msg := range resp.Messages {
		debugf("[DEBUG]: [%d of %d] body: %s\n", pos+1, total, *msg.Body)
		attributes := msg.MessageAttributes

		if verbose {

		}

		if len(attributes) > 0 {
			for k, v := range attributes {
				debugf("[DEBUG]: %s=%s\n", k, *v.StringValue)
			}
		}

		b, err := json.MarshalIndent(msg, "", " ")
		if err != nil {
			fmt.Printf("Error: %s", err)
			return false, err
		}
		fmt.Println(string(b))

	}

	if err != nil {
		return false, fmt.Errorf("Could not receive message(s) from queue '%s'@'%s': %s", message, queue, queueURL, err)
	}

	debugf("[DEBUG]: Successfully received %d message(s)\n", total)
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
