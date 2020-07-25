package aws

import (
	"fmt"
	"os"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/costexplorer"
	"github.com/aws/aws-sdk-go/service/costexplorer/costexploreriface"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3iface"
)

// CreateNewAwsSession Create a new aws session using access_key
// and secret_access_key from environment variables.
//
//Return one aws.session
func CreateNewAwsSession() *session.Session {
	awsConfig := &aws.Config{
		Credentials: credentials.NewStaticCredentials(os.Getenv("AWS_ACCESS_KEY_ID"), os.Getenv("AWS_SECRET_ACCESS_KEY"), ""),
		Region:      aws.String("us-east-1"),
	}

	return session.Must(session.NewSession(awsConfig))
}

// ListBuckets Receive s3 instance then list all bucker inside AWS S3.
//
// Return array with all buckets
func ListBuckets(client s3iface.S3API) []*s3.Bucket {
	listBucketsOutput, err := client.ListBuckets(&s3.ListBucketsInput{})
	if err != nil {
		fmt.Printf("Something goes wrong: %v", err)
	}

	return listBucketsOutput.Buckets
}

// ListObjects Receive a bucket name than iterate from all file pages counting
//how many files have inside and put the object size in an array.
//
// Return array with all object sizes and an integer with the files count
func ListObjects(bucketName *string, client s3iface.S3API, size *[]int64, filesCount *int, lastModified *time.Time) ([]int64, int, time.Time) {
	listObjectsV2Input := s3.ListObjectsV2Input{
		Bucket: bucketName,
	}

	client.ListObjectsV2Pages(&listObjectsV2Input,
		func(page *s3.ListObjectsV2Output, lastPage bool) bool {
			for _, item := range page.Contents {
				*size = append(*size, *item.Size)
				*filesCount++
				if item.LastModified.After(*lastModified) { // Otimizar no futuro. o(n)
					lastModified = item.LastModified
				}
			}
			return *page.IsTruncated
		})

	return *size, *filesCount, *lastModified
}

// CheckPrice Receive coast explorer instance and tagvalue(bucket name) than query in cost explorer
// to check how much this bucket spend this month.
//
// Return one string with the amount in dolars
func CheckPrice(client costexploreriface.CostExplorerAPI, tagValue string) string { // Melhorar essa função. Esta horrivel
	service := "Amazon Simple Storage Service"
	metricsValue := "BlendedCost"

	obj1 := &costexplorer.DimensionValues{
		Key:    aws.String("SERVICE"),
		Values: []*string{&service},
	}

	obj2 := &costexplorer.TagValues{
		Key:    aws.String("Name"),
		Values: []*string{&tagValue},
	}

	input := costexplorer.GetCostAndUsageInput{
		TimePeriod: &costexplorer.DateInterval{
			Start: aws.String("2020-07-01"),
			End:   aws.String("2020-07-31"),
		},
		Granularity: aws.String("MONTHLY"),
		Metrics:     []*string{&metricsValue},
		Filter: &costexplorer.Expression{
			And: []*costexplorer.Expression{
				{
					Dimensions: obj1,
				},
				{
					Tags: obj2,
				},
			},
		},
	}

	output, err := client.GetCostAndUsage(&input)
	if err != nil {
		fmt.Println(err)
	}

	return *output.ResultsByTime[0].Total["BlendedCost"].Amount
}
