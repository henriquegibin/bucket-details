package aws

import (
	"reflect"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3iface"
)

type mockS3Client struct {
	s3iface.S3API
}

var creationDate = time.Date(2020, time.July, 23, 22, 41, 0, 0, time.UTC)
var buckets = []*s3.Bucket{
	{
		CreationDate: &creationDate,
		Name:         aws.String("bucket1"),
	}, {
		CreationDate: &creationDate,
		Name:         aws.String("bucket2"),
	},
}

func (m *mockS3Client) ListBuckets(input *s3.ListBucketsInput) (*s3.ListBucketsOutput, error) {
	s3Output := s3.ListBucketsOutput{
		Buckets: buckets,
	}

	return &s3Output, nil
}

func TestListBuckets(t *testing.T) {
	mockS3Client := &mockS3Client{}

	type args struct {
		client s3iface.S3API
	}

	tests := []struct {
		name string
		args args
		want []*s3.Bucket
	}{
		{"List all buckets in S3", args{mockS3Client}, buckets},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ListBuckets(tt.args.client); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ListBuckets() = %v, want %v", got, tt.want)
			}
		})
	}
}
