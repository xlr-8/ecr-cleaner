package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ecr"
)

type Result struct {
	Data interface{} `json:"data"`
}

func main() {
	var (
		amountToKeep = flag.int("amount-to-keep", 100, "amount of images / repo you want to keep")
		awsRegion    = flag.String("aws.region", "eu-central-1", "AWS region")
		simulate     = flag.Bool("simulate", true, "simulate or not")
		list         interface{}
		err          error
	)

	flag.Parse()

	ecrCli := ecr.New(session.New(), aws.NewConfig().WithRegion(*awsRegion))

	repos, error := getAllRepoNames(ecrCli)

	if error != nil {
		log.Fatal(err)
	}

	log.Printf("repos %v", repos)

	if !*simulate {
		log.Print("jop")
	}

	for _, repo := range repos {
		images, error := getRepoImages(ecrCli, repo)

		if error != nil {
			log.Fatal(err)
		}

		resp, error := processRepo(ecrCli, repo, images, *simulate, *amountToKeep)

		log.Printf("response %v", resp)
	}

	err = json.NewEncoder(os.Stdout).Encode(Result{Data: list})

	if err != nil {
		log.Fatal(err)
	}
}

func processRepo(ecrCli *ecr.ECR, repoName string, images []*ecr.ImageIdentifier, simulate bool, amountToKeep int) (bool, error) {
	var imageIds []*ecr.ImageIdentifier

	tagFiltered := filterByTag(images)
	log.Printf("images without tag %v", len(tagFiltered))

	for _, image := range tagFiltered {
		imageIds = append(imageIds, image)
	}

	if !simulate {
		resp, err := ecrCli.BatchDeleteImage(&ecr.BatchDeleteImageInput{
			RepositoryName: aws.String(repoName),
			ImageIds:       images,
		})

		if err != nil {
			return false, fmt.Errorf("deleting %v images: %v", repoName, err)
		}

		log.Printf("jop: %v", resp)
	} else {
		log.Print("simulation...")
		log.Printf("repo %v contains %v images and %v will be deleted", repoName, len(images), len(imageIds))
		log.Print("images to delete: $v", imageIds)
	}

	return true, nil
}

func filterByTag(images []*ecr.ImageIdentifier) (imgs []*ecr.ImageIdentifier) {
	for _, image := range images {
		if image.ImageTag == nil {
			imgs = append(imgs, image)
		}
	}

	return imgs
}

func getRepoImages(ecrCli *ecr.ECR, repoName string) ([]*ecr.ImageIdentifier, error) {
	resp, err := ecrCli.ListImages(&ecr.ListImagesInput{
		RepositoryName: aws.String(repoName),
	})

	if err != nil {
		return nil, fmt.Errorf("getting %v images: %v", repoName, err)
	}

	return resp.ImageIds, nil
}

func getAllRepoNames(ecrCli *ecr.ECR) ([]string, error) {

	resp, err := ecrCli.DescribeRepositories(&ecr.DescribeRepositoriesInput{})

	if err != nil {
		return nil, fmt.Errorf("getting ecr repos: %v", err)
	}

	repos := make([]string, 0, len(resp.Repositories))

	for _, repo := range resp.Repositories {
		repos = append(repos, *repo.RepositoryName)
	}
	return repos, nil
}
