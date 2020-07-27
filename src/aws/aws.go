package aws

import (
	errorchecker "bucket-details/src/errorcheck"
	genfun "bucket-details/src/genericfunctions"
	"bucket-details/src/structs"
	"errors"
	"os"
	"regexp"
	"strings"
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
// Return one aws.session.
func CreateNewAwsSession() *session.Session {
	awsConfig := &aws.Config{
		Credentials: credentials.NewStaticCredentials(os.Getenv("AWS_ACCESS_KEY_ID"), os.Getenv("AWS_SECRET_ACCESS_KEY"), ""),
		Region:      aws.String("us-east-1"),
	}

	return session.Must(session.NewSession(awsConfig))
}

// ListBuckets Receive s3 instance then list all bucker inside AWS S3.
//
// Return array with all buckets.
func ListBuckets(client s3iface.S3API, flags structs.Flags) []*s3.Bucket {
	listBucketsOutput, err := client.ListBuckets(&s3.ListBucketsInput{})
	errorchecker.CheckFatal(err, "ListBuckets")

	if flags.FilterType != "" {
		buckets := filterBuckets(listBucketsOutput.Buckets, flags)
		return buckets
	}

	return listBucketsOutput.Buckets
}

// ListObjects Receive a bucket name than iterate from all file pages counting
//how many files have inside and put the object size in an array.
//
// Return array with all object sizes and an integer with the files count.
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

// CheckS3BucketCost Receive coast explorer instance and tagvalue(bucket name) than query in cost explorer
// to check how much this bucket spend this month.
//
// Return one string with the amount in dolars.
func CheckS3BucketCost(client costexploreriface.CostExplorerAPI, tagValue string) string { // Melhorar essa função. Esta horrivel
	service := "Amazon Simple Storage Service"
	metricsValue := "BlendedCost"

	var dimensionValues costexplorer.DimensionValues
	var tagValues costexplorer.TagValues
	var dateInterval costexplorer.DateInterval
	var filterObject costexplorer.Expression

	dimensionValues.SetKey("SERVICE")
	dimensionValues.SetValues([]*string{&service})

	tagValues.SetKey("Name")
	tagValues.SetValues([]*string{&tagValue})

	dateInterval.SetStart(genfun.GetFirstLastDayOfMonth("first").Format("2006-01-02"))
	dateInterval.SetEnd(genfun.GetFirstLastDayOfMonth("last").Format("2006-01-02"))

	filterObject.SetAnd(createFilterExpressionSlice(dimensionValues, tagValues))

	input := costexplorer.GetCostAndUsageInput{}
	input.SetTimePeriod(&dateInterval)
	input.SetGranularity("MONTHLY")
	input.SetMetrics([]*string{&metricsValue})
	input.SetFilter(&filterObject)

	output, err := client.GetCostAndUsage(&input)
	errorchecker.CheckError(err, "CheckS3BucketCost")

	return *output.ResultsByTime[0].Total["BlendedCost"].Amount
}

// GetBucketLifeCycle Receive s3 client and bucket name.
//
// Return a lifecycle output.
func GetBucketLifeCycle(client s3iface.S3API, bucketName *string) ([]*s3.LifecycleRule, error) {
	lifeCycleOutput, err := client.GetBucketLifecycleConfiguration(&s3.GetBucketLifecycleConfigurationInput{Bucket: bucketName})
	if err != nil {
		return nil, err
	}

	return lifeCycleOutput.Rules, err
}

// GetBucketACL Receive s3 client and bucket name.
//
// Return a list of users and level permissions.
func GetBucketACL(client s3iface.S3API, bucketName *string) ([]*s3.Grant, error) {
	aclOutput, err := client.GetBucketAcl(&s3.GetBucketAclInput{Bucket: bucketName})
	if err != nil {
		return nil, err
	}

	return aclOutput.Grants, nil
}

// GetBucketEncryption Receive s3 client and bucket name.
//
// Return a slice with all server side encriptyon rules.
func GetBucketEncryption(client s3iface.S3API, bucketName *string) ([]*s3.ServerSideEncryptionRule, error) {
	encryptionOutput, err := client.GetBucketEncryption(&s3.GetBucketEncryptionInput{Bucket: bucketName})
	if err != nil {
		return nil, err
	}

	return encryptionOutput.ServerSideEncryptionConfiguration.Rules, nil
}

// GetBucketLocation Receive s3 client and bucket name.
//
// Return a string representation of bucket location.
func GetBucketLocation(client s3iface.S3API, bucketName *string) (string, error) {
	locationOutput, err := client.GetBucketLocation(&s3.GetBucketLocationInput{Bucket: bucketName})
	if err != nil {
		return "", err
	}

	return *locationOutput.LocationConstraint, nil
}

// GetBucketTagging Receive s3 client and bucket name.
//
// Return a string representation of a bucket tagSet.
func GetBucketTagging(client s3iface.S3API, bucketName *string) ([]*s3.Tag, error) {
	taggingOutput, err := client.GetBucketTagging(&s3.GetBucketTaggingInput{Bucket: bucketName})
	if err != nil {
		return nil, err
	}

	return taggingOutput.TagSet, nil
}

// createFilterExpressionSlice Receive one DimensionValue object and one TagValue object then create
// a filterExpressionSlice to use during the GetCostAndUsage.
//
// Return a slice with Expression pointers.
func createFilterExpressionSlice(dimensionValues costexplorer.DimensionValues, tagValue costexplorer.TagValues) []*costexplorer.Expression {
	var expressionSlice []*costexplorer.Expression
	var dimensionValuesExpression costexplorer.Expression
	var tagValueExpression costexplorer.Expression

	dimensionValuesExpression.SetDimensions(&dimensionValues)
	tagValueExpression.SetTags(&tagValue)

	expressionSlice = append(expressionSlice, &dimensionValuesExpression, &tagValueExpression)
	return expressionSlice
}

func filterBuckets(buckets []*s3.Bucket, flags structs.Flags) []*s3.Bucket {
	var filteredBuckets []*s3.Bucket

	switch flags.FilterType {
	case "prefix":
		for _, bucket := range buckets {
			if strings.HasPrefix(*bucket.Name, flags.FilterValue) {
				filteredBuckets = append(filteredBuckets, bucket)
			}
		}
	case "suffix":
		for _, bucket := range buckets {
			if strings.HasSuffix(*bucket.Name, flags.FilterValue) {
				filteredBuckets = append(filteredBuckets, bucket)
			}
		}
	case "regex":
		for _, bucket := range buckets {
			regex := regexp.MustCompile(flags.FilterValue)
			if regex.MatchString(*bucket.Name) {
				filteredBuckets = append(filteredBuckets, bucket)
			}
		}
	default:
		e := errors.New("Invalid filterType value. Please use prefix or suffix or regex")
		errorchecker.CheckFatal(e, "filterBuckets")
	}

	return filteredBuckets
}
