package main

import (
	"flag"
	"fmt"
	"log"
	"sort"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ecr"
	"github.com/xlab/handysort"
)

type ByAlphabet []*ecr.ImageIdentifier

func main() {
	var (
		amountToKeep  = flag.Int("amount-to-keep", 100, "amount of images / repo you want to keep")
		awsRegion     = flag.String("aws.region", "eu-central-1", "AWS region")
		repoToProcess = flag.String("repository", "", "AWS region")
		dryRun        = flag.Bool("dry-run", false, "dry run")
		error         error
	)

	flag.Parse()

	ecrCli := ecr.New(session.New(), aws.NewConfig().WithRegion(*awsRegion))

	var repos []string

	if *repoToProcess != "" {
		repos = []string{*repoToProcess}
	} else {
		repos, error = getAllRepoNames(ecrCli)
	}

	if error != nil {
		log.Fatal(error)
	}

	log.Printf("Repositories to process: %v", repos)

	for _, repo := range repos {
		images, error := getRepoImages(ecrCli, repo, "")
		if error != nil {
			log.Fatal(error)
		}
		log.Printf("Number of images in %v: %v", repo, len(images))

		_, error = processRepo(ecrCli, repo, images, *dryRun, *amountToKeep)

		if error != nil {
			log.Fatal(error)
		}
	}
}

func (s ByAlphabet) Len() int {
	return len(s)
}
func (s ByAlphabet) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s ByAlphabet) Less(i, j int) bool {
	return handysort.StringLess(*s[i].ImageTag, *s[j].ImageTag)
}

func processRepo(ecrCli *ecr.ECR, repoName string, images []*ecr.ImageIdentifier, dryRun bool, amountToKeep int) (bool, error) {
	var imageIdents []*ecr.ImageIdentifier

	tagFiltered, remainingImages := filterByTag(images)

	for _, image := range tagFiltered {
		imageIdents = append(imageIdents, image)
	}

	amountFiltered := filterByAmount(remainingImages, amountToKeep)

	for _, image := range amountFiltered {
		imageIdents = append(imageIdents, image)
	}

	log.Printf("number of images to delete: %v", len(imageIdents))

	if !dryRun {
		_, err := ecrCli.BatchDeleteImage(&ecr.BatchDeleteImageInput{
			RepositoryName: aws.String(repoName),
			ImageIds:       imageIdents,
		})

		if err != nil {
			return false, fmt.Errorf("deleting %v images: %v", repoName, err)
		}

		log.Printf("deleted %v images in %v", len(imageIdents), repoName)
	} else {
		log.Print("dry run ...")
		log.Print("images to delete: $v", imageIdents)
	}

	return true, nil
}

func filterByAmount(images []*ecr.ImageIdentifier, amountToKeep int) []*ecr.ImageIdentifier {
	sort.Sort(ByAlphabet(images))

	if len(images) < amountToKeep {
		var a []*ecr.ImageIdentifier
		return a
	}

	return images[0 : len(images)-amountToKeep]
}

func filterByTag(images []*ecr.ImageIdentifier) (imagesWithoutTag []*ecr.ImageIdentifier, imagesWithTag []*ecr.ImageIdentifier) {
	for _, image := range images {
		if image.ImageTag == nil {
			imagesWithoutTag = append(imagesWithoutTag, image)
		} else {
			imagesWithTag = append(imagesWithTag, image)
		}
	}

	return imagesWithoutTag, imagesWithTag
}

func getRepoImages(ecrCli *ecr.ECR, repoName string, token string) ([]*ecr.ImageIdentifier, error) {
	var resp *ecr.ListImagesOutput
	var err error

	if token != "" {
		resp, err = ecrCli.ListImages(&ecr.ListImagesInput{
			RepositoryName: aws.String(repoName),
			NextToken:      aws.String(token),
		})
	} else {
		resp, err = ecrCli.ListImages(&ecr.ListImagesInput{
			RepositoryName: aws.String(repoName),
		})
	}

	if err != nil {
		return nil, fmt.Errorf("getting %v images: %v", repoName, err)
	}

	imageIdents := resp.ImageIds

	if resp.NextToken != nil {
		resp2, err2 := getRepoImages(ecrCli, repoName, *resp.NextToken)

		if err2 != nil {
			return nil, fmt.Errorf("getting %v images: %v", repoName, err2)
		}

		imageIdents = append(imageIdents, resp2...)
	}

	return imageIdents, nil
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
