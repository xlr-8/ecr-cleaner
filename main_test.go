package main

import (
	"testing"

	"github.com/aws/aws-sdk-go/aws"

	"github.com/aws/aws-sdk-go/service/ecr"
)

func TestSeparateImages(t *testing.T) {
	testcases := map[string]struct {
		images       []*ecr.ImageDetail
		expTagsLen   int
		expNoTagsLen int
	}{
		"empty": {images: []*ecr.ImageDetail{},
			expTagsLen:   0,
			expNoTagsLen: 0,
		},
		"simple": {images: []*ecr.ImageDetail{
			{ImageDigest: aws.String("foo"), ImageTags: []*string{aws.String("foo")}},
			{ImageDigest: aws.String("bar")},
		},
			expTagsLen:   1,
			expNoTagsLen: 1,
		},
	}

	for _, testcase := range testcases {

		noTag, withTag := separateHavingTag(testcase.images)
		if len(noTag) != testcase.expNoTagsLen {
			t.Errorf("want %d with no tag got %d", testcase.expNoTagsLen, len(noTag))
		}
		if len(withTag) != testcase.expTagsLen {
			t.Errorf("want %d with tag got %d", testcase.expTagsLen, len(withTag))
		}
	}
}
