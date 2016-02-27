// Package awscli -
package awscli

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ecr"
)

// ECRInfo -
func ECRInfo(registryID string) {
	svc := ecr.New(session.New())

	params := &ecr.DescribeRepositoriesInput{
		MaxResults: aws.Int64(100),
		// NextToken:  aws.String("NextToken"),
		// RegistryId: aws.String("RegistryId"),
		RegistryId: aws.String(registryID),
		// RepositoryNames: []*string{
		// 	aws.String("awscli"), // Required
		// 	// More values...
		// },
	}
	resp, err := svc.DescribeRepositories(params)

	if err != nil {
		// Print the error, cast err to awserr.Error to get the Code and
		// Message from an error.
		fmt.Println(err.Error())
		return
	}

	// Pretty-print the response data.
	fmt.Println(resp)
}
