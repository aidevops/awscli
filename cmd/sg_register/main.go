// Package main - sg_register application
package main

// import - import our dependencies
import (
	"flag"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
)

// Unit - this application's name
const Unit = "sg_register"

// verbose - control debug output
var verbose bool

// main - log us in...
func main() {
	var (
		ip         string
		port       string
		id         string
		name       string
		region     string
		register   bool
		deregister bool
		version    bool
	)

	flag.StringVar(&ip, "ip", "", "ip address to register")
	flag.StringVar(&port, "port", "", "port to register access to...")
	flag.StringVar(&id, "id", "", "security group id to work against (mutually exclusive to name)")
	flag.StringVar(&name, "name", "", "security group name to work against (mutually exclusive to id")
	flag.StringVar(&region, "region", "us-east-1", "region sg lives in...")
	flag.BoolVar(&verbose, "verbose", false, "be more verbose.....")
	flag.BoolVar(&version, "version", false, "print version and exit")
	flag.Parse()

	if version == true {
		fmt.Println(versionInfo())
		os.Exit(0)
	}

	if len(id) <= 0 && len(name) <= 0 {
		fmt.Println("sg_register: you need to specify either -id or -name")
		os.Exit(1)
	}

	if !register && !deregister {
		fmt.Println("sg_register: you need to specify either -register or -deregister")
		os.Exit(1)
	}

	if register && deregister {
		fmt.Println("sg_register: register and deregister are mutually exclusive")
		os.Exit(1)
	}

	debugf("[DEBUG]: using ip address(s): %s\n", ip)
	if ip == "" || len(ip) < 7 {
		fmt.Printf("sg_register: missing or invalid ip: -ip='1.1.1.1..', received: '%s'\n", ip)
		os.Exit(253)
	}

	debugf("[DEBUG]: using region: %s\n", region)

	var ok bool
	var err error
	if register {
		ok, err = Register(region, verbose, ip, port, id, name)
	}

	if deregister {
		ok, err = Deregister(region, verbose, ip, port, id, name)
	}

	if !ok {
		fmt.Printf("[ERROR]: failed while processing request: %s", err)
		os.Exit(253)
	}

	// success!!!
	os.Exit(0)

}

// Register - register instance ip with security group
func Register(region string, verbose bool, ip, port, id, name string) (ok bool, err error) {
	debugf("[DEBUG]: creating new session...\n")
	svc := ec2.New(session.New(&aws.Config{Region: aws.String(region)}))

	params := &ec2.AuthorizeSecurityGroupIngressInput{
		CidrIp:    aws.String("String"),
		DryRun:    aws.Bool(true),
		FromPort:  aws.Int64(1),
		GroupId:   aws.String("String"),
		GroupName: aws.String("String"),
		IpPermissions: []*ec2.IpPermission{
			{ // Required
				FromPort:   aws.Int64(1),
				IpProtocol: aws.String("String"),
				IpRanges: []*ec2.IpRange{
					{ // Required
						CidrIp: aws.String("String"),
					},
					// More values...
				},
				PrefixListIds: []*ec2.PrefixListId{
					{ // Required
						PrefixListId: aws.String("String"),
					},
					// More values...
				},
				ToPort: aws.Int64(1),
				UserIdGroupPairs: []*ec2.UserIdGroupPair{
					{ // Required
						GroupId:       aws.String("String"),
						GroupName:     aws.String("String"),
						PeeringStatus: aws.String("String"),
						UserId:        aws.String("String"),
						VpcId:         aws.String("String"),
						VpcPeeringConnectionId: aws.String("String"),
					},
					// More values...
				},
			},
			// More values...
		},
		IpProtocol:                 aws.String("String"),
		SourceSecurityGroupName:    aws.String("String"),
		SourceSecurityGroupOwnerId: aws.String("String"),
		ToPort: aws.Int64(1),
	}
	debugf("[DEBUG]: registering...\n")
	resp, err := svc.AuthorizeSecurityGroupIngress(params)
	debugf("[DEBUG]: response: %v\n", resp)

	if err != nil {
		// Print the error, cast err to awserr.Error to get the Code and
		// Message from an error.
		fmt.Println(err.Error())
		return false, err
	}

	// Pretty-print the response data.
	fmt.Println(resp)
	return true, nil
}

// Deregister - deregister instance ip from security group
func Deregister(region string, verbose bool, ip, port, id, name string) (ok bool, err error) {
	debugf("[DEBUG]: creating new session and s3manager object...\n")
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
