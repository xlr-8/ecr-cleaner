package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ecr"
)

type Result struct {
	Data interface{} `json:"data"`
}

func main() {
	var (
		discoveryType = flag.String("type", "", "type of discovery. EC2, ELB, RDS, CloudFront, Lambda or ECSClusters")
		awsRegion     = flag.String("aws.region", "eu-central-1", "AWS region")
		list          interface{}
		err           error
	)

	flag.Parse()

	ecrCli := session.New(session.New(), aws.NewConfig().WithRegion(*awsRegion))
    
    repos, error := getAllRepos(ecrCli)

	err = json.NewEncoder(os.Stdout).Encode(Result{Data: list})

	if err != nil {
		log.Fatal(err)
	}
}

func getAllRepos(ecrCli ecr.ECR) ([]string, error) {

	resp, err := ecrCli.DescribeRepositories(&ecr.DescribeRepositoriesInput{})

	if err != nil {
		return nil, fmt.Errorf("getting EC2 instances: %v", err)
	}

	ec2Identifiers := make([]string, 0, len(resp.))

	for _, reservation := range resp.Reservations {


	}
	return ec2Identifiers, nil
}