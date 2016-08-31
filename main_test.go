package main

import (
	"testing"

	"github.com/aws/aws-sdk-go/aws"

	"github.com/aws/aws-sdk-go/service/ecr"
)

func TestSeparateImages(t *testing.T) {
	images := []*ecr.ImageIdentifier{
		{ImageDigest: aws.String("foo"), ImageTag: aws.String("foo")},
		{ImageDigest: aws.String("bar")},
	}
	noTag, withTag := separateHavingTag(images)
	if len(noTag) != 1 {
		t.Errorf("want 1 with no tag got %s", len(noTag))
	}
	if len(withTag) != 1 {
		t.Errorf("want 1 with tag got %s", len(withTag))
	}
}
