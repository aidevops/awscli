// Package main - sg_register application
package main

// import - import our dependencies
import (
	"flag"
	"fmt"
	"net"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
)

// Unit - this application's name
const Unit = "sg_register"

// verbose - control debug output
var verbose bool
var dryrun bool

// main - log us in...
func main() {
	var (
		ip         string
		protocol   string
		fromPort   int64
		toPort     int64
		sid        string
		name       string
		region     string
		register   bool
		deregister bool
		version    bool
	)

	flag.StringVar(&ip, "ip", "0.0.0.0/0", "ip address to register")
	flag.StringVar(&protocol, "protocol", "tcp", "protocol to register 'tcp','udp','icmp','all'")
	flag.Int64Var(&fromPort, "from-port", 443, "start port range to register access to...")
	flag.Int64Var(&toPort, "to-port", -1, "end port range to register access to...")
	flag.StringVar(&sid, "sg-id", "", "security group id to work against (mutually exclusive to name - not implemented)")
	flag.StringVar(&name, "sg-name", "", "security group name to work against (mutually exclusive to sg-id)")
	flag.StringVar(&region, "region", "us-east-1", "region sg lives in...")
	flag.BoolVar(&register, "register", false, "register with security group ingress.....")
	flag.BoolVar(&deregister, "deregister", false, "deregister with security group ingress.....")
	flag.BoolVar(&verbose, "verbose", false, "be more verbose.....")
	flag.BoolVar(&version, "version", false, "print version and exit")
	flag.BoolVar(&dryrun, "dryrun", false, "perform dryrun and exit")
	flag.Parse()

	if version == true {
		fmt.Println(versionInfo())
		os.Exit(0)
	}

	if len(sid) <= 0 && len(name) <= 0 {
		fmt.Println("sg_register: you need to specify either -sg-id or -name")
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
	if _, _, err := net.ParseCIDR(ip); err != nil {
		fmt.Printf("sg_register: %s\n", err)
		os.Exit(1)
	}

	debugf("[DEBUG]: using region: %s\n", region)

	var ok bool
	var err error
	if register {
		ok, err = Register(region, ip, protocol, fromPort, toPort, sid, name)
	}

	if deregister {
		ok, err = Deregister(region, ip, protocol, fromPort, toPort, sid, name)
	}

	if !ok {
		fmt.Printf("[ERROR]: failed while processing request: %s", err)
		os.Exit(253)
	}

	fmt.Printf("success: '%s'", ip)
	os.Exit(0)

}

// Register - register instance ip with security group
func Register(region, ip, protocol string, fromPort, toPort int64, sid, name string) (ok bool, err error) {
	debugf("[DEBUG]: creating new session...\n")
	svc := ec2.New(session.New(&aws.Config{Region: aws.String(region)}))

	if sid == "" {
		sid, err = LookupSGID(name, svc)
		if err != nil {
			fmt.Printf("[ERROR]: failed to lookup sg '%s' by name: %s", name, err)
		}
	}

	params := &ec2.AuthorizeSecurityGroupIngressInput{
		CidrIp:     aws.String(ip),
		DryRun:     aws.Bool(dryrun),
		GroupId:    aws.String(sid),
		IpProtocol: aws.String(protocol),
		ToPort:     aws.Int64(toPort),
		FromPort:   aws.Int64(fromPort),
	}
	debugf("[DEBUG]: registering...\n")
	resp, err := svc.AuthorizeSecurityGroupIngress(params)
	debugf("[DEBUG]: response: %v\n", resp)

	if err != nil {
		return false, err
	}

	return true, nil
}

// LookupSGID - Lookup security group by name, return its id
func LookupSGID(name string, svc *ec2.EC2) (sid string, err error) {
	params := &ec2.DescribeSecurityGroupsInput{
		DryRun: aws.Bool(dryrun),
		Filters: []*ec2.Filter{
			{
				Name: aws.String("group-name"),
				Values: []*string{
					aws.String(name),
				},
			},
		},
	}
	resp, err := svc.DescribeSecurityGroups(params)

	if err != nil {
		return sid, err
	}

	// read the first one and exit
	for _, res := range resp.SecurityGroups {
		sid = *res.GroupId
		break
	}
	return sid, err
}

// Deregister - deregister instance ip from security group
func Deregister(region, ip, protocol string, fromPort, toPort int64, sid, name string) (ok bool, err error) {
	debugf("[DEBUG]: creating new session...\n")
	svc := ec2.New(session.New(&aws.Config{Region: aws.String(region)}))

	if sid == "" {
		sid, err = LookupSGID(name, svc)
		if err != nil {
			fmt.Printf("[ERROR]: failed to lookup sg '%s' by name: %s", name, err)
		}
	}

	params := &ec2.RevokeSecurityGroupIngressInput{
		CidrIp:     aws.String(ip),
		DryRun:     aws.Bool(dryrun),
		GroupId:    aws.String(sid),
		IpProtocol: aws.String(protocol),
		ToPort:     aws.Int64(toPort),
		FromPort:   aws.Int64(fromPort),
	}
	debugf("[DEBUG]: deregistering...\n")
	resp, err := svc.RevokeSecurityGroupIngress(params)
	debugf("[DEBUG]: response: %v\n", resp)

	if err != nil {
		return false, err
	}

	return true, nil
}

// debugf - print to stdout if verbose is enabled....
func debugf(format string, args ...interface{}) {
	if verbose == true {
		fmt.Printf(format, args...)
	}
}

// versionInfo - return version info
func versionInfo() string {
	return fmt.Sprintf("%s v%s.%s (%s)", Unit, Version, VersionPrerelease, GitCommit)
}
