// Package main - s3_util application
package main

// import - import our dependencies
import (
	"flag"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

// Unit - this application's name
const Unit = "s3_util"

// verbose - control debug output
var verbose bool

// main - log us in...
func main() {
	var (
		account string
		region  string
		bucket  string
		retry   int64
		src     string
		dst     string
		put     bool
		get     bool
		version bool
	)

	var empty string
	flag.StringVar(&account, "account", "", "AWS account #. E.g. -account='1234556790123'")
	flag.StringVar(&region, "region", "us-east-1", "AWS region. E.g. -region=us-east-1")
	flag.StringVar(&bucket, "bucket", "", "mybucket-name...")
	flag.Int64Var(&retry, "retry", 3, "number of times to attempt the operation - not implemented")
	flag.StringVar(&src, "src", "", "/path/to/my/object")
	flag.StringVar(&dst, "dst", "", "/path/to/my/object")
	flag.BoolVar(&put, "put", false, "put object")
	flag.BoolVar(&get, "get", false, "get object")
	flag.BoolVar(&verbose, "verbose", false, "be more verbose.....")
	flag.BoolVar(&version, "version", false, "print version and exit")
	flag.Parse()

	if version == true {
		fmt.Println(versionInfo())
		os.Exit(0)
	}

	if !put && !get {
		fmt.Println("s3_util: you need to specify either -put or -get")
		os.Exit(1)
	}

	if put && get {
		fmt.Println("s3_util: put and get are mutually exclusive")
		os.Exit(1)
	}

	debugf("[DEBUG]: using retry count: %d\n", retry)
	if retry < 0 || retry > 10 {
		fmt.Printf("s3_util: invalid count valid values 1 - 10\n")
		os.Exit(255)
	}

	debugf("[DEBUG]: using account: %s\n", account)
	if account == "" || len(account) < 12 {
		fmt.Printf("s3_util: missing or invalid account length: -account='1234556790123', received: '%s'\n", account)
		os.Exit(254)
	}

	debugf("[DEBUG]: using bucket name(s): %s\n", bucket)
	if bucket == "" || len(bucket) < 3 {
		fmt.Printf("s3_util: missing or invalid bucket: -bucket='some-fancy-bucket..', received: '%s'\n", bucket)
		os.Exit(253)
	}

	debugf("[DEBUG]: using region: %s\n", region)

	var ok bool
	var err error
	if put {
		ok, err = Put(account, region, verbose, bucket, retry, src, dst)
	}

	if get {
		ok, err = Get(account, region, verbose, bucket, retry, src, dst)
	}

	if !ok {
		fmt.Printf("[ERROR]: failed while processing request: %s", err)
		os.Exit(253)
	}

	// success!!!
	os.Exit(0)

}

// Put - place a file in aws s3
func Put(account, region string, verbose bool, bucket string, retry int64, src, dst string) (ok bool, err error) {
	file, err := os.Open(src)
	if err != nil {
		return false, fmt.Errorf("Failed to open source file '%s': %s\n", src, err)
	}

	// Not required, but you could zip the file before uploading it
	// using io.Pipe read/writer to stream gzip'd file contents.
	// reader, writer := io.Pipe()
	// go func() {
	// 	gw := gzip.NewWriter(writer)
	// 	io.Copy(gw, file)

	// 	file.Close()
	// 	gw.Close()
	// 	writer.Close()
	// }()

	debugf("[DEBUG]: creating new session and s3manager object...\n")
	svc := s3manager.NewUploader(session.New(&aws.Config{Region: aws.String(region)}))
	debugf("[DEBUG]: uploading...\n")
	resp, err := svc.Upload(&s3manager.UploadInput{
		Body:   file,
		Bucket: aws.String(bucket),
		Key:    aws.String(src),
	})
	debugf("[DEBUG]: response: %v\n", resp)

	if err != nil {
		return false, fmt.Errorf("Could not put file '%s' into bucket '%s:%s': %s\n", src, bucket, dst, err)
	}

	debugf("[DEBUG]: Successfully placed file(s) '%s' into '%s'\n", src, resp.Location)
	return true, nil
}

// Get - Get file from aws s3
func Get(account, region string, verbose bool, bucket string, retry int64, src, dst string) (ok bool, err error) {
	debugf("[DEBUG]: creating new session and s3manager object...\n")
	svc := s3manager.NewDownloader(session.New(&aws.Config{Region: aws.String(region)}))

	debugf("[DEBUG]: downloading...\n")
	resp, err := svc.Download(file, &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(src),
	})
	debugf("[DEBUG]: response: %v\n", resp)

	if err != nil {
		return false, fmt.Errorf("Could not get file '%s' from bucket '%s:%s': %s\n", dst, bucket, src, err)
	}

	debugf("[DEBUG]: Successfully retrieved file(s) '%s' from '%s%s' - %d\n", dst, bucket, src, resp)
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
