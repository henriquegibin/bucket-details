package aws

import (
	"reflect"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/costexplorer"
	"github.com/aws/aws-sdk-go/service/costexplorer/costexploreriface"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3iface"
)

type mockS3Client struct {
	s3iface.S3API
}

type mockCostExplorerClient struct {
	costexploreriface.CostExplorerAPI
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

var size []int64
var filesCount = 0
var lastModified time.Time

func (m *mockS3Client) ListObjectsV2Pages(input *s3.ListObjectsV2Input, fn func(page *s3.ListObjectsV2Output, lastPage bool) bool) error {
	size = []int64{int64(100)}
	filesCount = 1
	lastModified = time.Date(2020, time.July, 23, 22, 41, 0, 0, time.UTC)
	return nil
} // melhorar esse mock. Tentar iterar entre dois buckets no minimo. Hoje ele só valida se ta retornando certo

func TestListObjects(t *testing.T) {
	mockS3Client := &mockS3Client{}

	type args struct {
		bucketName   *string
		client       s3iface.S3API
		size         *[]int64
		filesCount   *int
		lastModified *time.Time
	}

	tests := []struct {
		name  string
		args  args
		want  []int64
		want1 int
		want2 time.Time
	}{
		{
			"Retrieve metadata from each bucket object",
			args{aws.String("bucket"), mockS3Client, &size, &filesCount, &lastModified},
			[]int64{int64(100)},
			1,
			time.Date(2020, time.July, 23, 22, 41, 0, 0, time.UTC),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, got2 := ListObjects(tt.args.bucketName, tt.args.client, tt.args.size, tt.args.filesCount, tt.args.lastModified)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ListObjects() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("ListObjects() got1 = %v, want %v", got1, tt.want1)
			}
			if got2 != tt.want2 {
				t.Errorf("ListObjects() got2 = %v, want %v", got2, tt.want2)
			}
		})
	}
}

func (m *mockCostExplorerClient) GetCostAndUsage(input *costexplorer.GetCostAndUsageInput) (*costexplorer.GetCostAndUsageOutput, error) {
	metricValue := costexplorer.MetricValue{}
	metricValue.SetAmount("10.00")

	total := map[string]*costexplorer.MetricValue{
		"BlendedCost": &metricValue,
	}

	resultByTime := *&costexplorer.ResultByTime{
		Total: total,
	}
	resultByTimeColletion := make([]*costexplorer.ResultByTime, 1)
	resultByTimeColletion[0] = &resultByTime

	output := &costexplorer.GetCostAndUsageOutput{}
	output.SetResultsByTime(resultByTimeColletion)

	return output, nil
}

func TestCheckPrice(t *testing.T) {
	mockCostExplorerClient := &mockCostExplorerClient{}

	type args struct {
		client   costexploreriface.CostExplorerAPI
		tagValue string
	}

	tests := []struct {
		name string
		args args
		want string
	}{
		{"Get Cost from a bucket", args{mockCostExplorerClient, "bucketName"}, "10.00"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CheckPrice(tt.args.client, tt.args.tagValue); got != tt.want {
				t.Errorf("CheckPrice() = %v, want %v", got, tt.want)
			}
		})
	}
}
