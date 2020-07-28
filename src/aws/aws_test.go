package aws

import (
	"bucket-details/src/structs"
	"errors"
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
var allBuckets = []*s3.Bucket{
	{
		CreationDate: &creationDate,
		Name:         aws.String("1bucket"),
	}, {
		CreationDate: &creationDate,
		Name:         aws.String("bucket2"),
	}, {
		CreationDate: &creationDate,
		Name:         aws.String("4-bucket.example.net"),
	},
}
var prefixBuckets = []*s3.Bucket{allBuckets[0]}
var suffixBuckets = []*s3.Bucket{allBuckets[1]}
var regexBuckets = []*s3.Bucket{allBuckets[2]}

func (m *mockS3Client) ListBuckets(input *s3.ListBucketsInput) (*s3.ListBucketsOutput, error) {
	s3Output := s3.ListBucketsOutput{
		Buckets: allBuckets,
	}

	return &s3Output, nil
}

func TestListBuckets(t *testing.T) {
	mockS3Client := &mockS3Client{}

	type args struct {
		client s3iface.S3API
		flags  structs.Flags
	}

	tests := []struct {
		name string
		args args
		want []*s3.Bucket
	}{
		{"List all buckets in S3 without filter", args{mockS3Client, structs.Flags{FilterType: "", FilterValue: ""}}, allBuckets},
		{"List all buckets in S3 with prefix filter", args{mockS3Client, structs.Flags{FilterType: "prefix", FilterValue: "1"}}, prefixBuckets},
		{"List all buckets in S3 with suffix filter", args{mockS3Client, structs.Flags{FilterType: "suffix", FilterValue: "2"}}, suffixBuckets},
		{"List all buckets in S3 with regex filter", args{mockS3Client, structs.Flags{FilterType: "regex", FilterValue: `^[0-9]-[A-z]+\.example\.net$`}}, regexBuckets},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ListBuckets(tt.args.client, tt.args.flags); !reflect.DeepEqual(got, tt.want) {
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
		"UnblendedCost": &metricValue,
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

func TestCheckS3BucketCost(t *testing.T) {
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
			if got := CheckS3BucketCost(tt.args.client, tt.args.tagValue, false); got != tt.want {
				t.Errorf("CheckS3BucketCost() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_createFilterExpressionSlice(t *testing.T) {
	var dimension costexplorer.DimensionValues
	var tags costexplorer.TagValues
	var values = aws.StringSlice([]string{"test"})

	dimension.SetKey("Key")
	dimension.SetValues(values)

	tags.SetKey("Key")
	tags.SetValues(values)

	var expect = []*costexplorer.Expression{
		{Dimensions: &dimension},
		{Tags: &tags},
	}

	type args struct {
		dimensionValues costexplorer.DimensionValues
		tagValue        costexplorer.TagValues
	}

	tests := []struct {
		name string
		args args
		want []*costexplorer.Expression
	}{
		{"Generate a slice of Expressions with dimensionValues and tagValues", args{dimension, tags}, expect},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := createFilterExpressionSlice(tt.args.dimensionValues, tt.args.tagValue); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("createFilterExpressionSlice() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_filterBuckets(t *testing.T) {
	type args struct {
		buckets []*s3.Bucket
		flags   structs.Flags
	}

	tests := []struct {
		name    string
		args    args
		want    []*s3.Bucket
		wantErr bool
	}{
		{"Return buckets after prefix filter", args{allBuckets, structs.Flags{FilterType: "prefix", FilterValue: "1"}}, prefixBuckets, false},
		{"Return buckets after suffix filter", args{allBuckets, structs.Flags{FilterType: "suffix", FilterValue: "2"}}, suffixBuckets, false},
		{"Return buckets after regex filter", args{allBuckets, structs.Flags{FilterType: "regex", FilterValue: `^[0-9]-[A-z]+\.example\.net$`}}, regexBuckets, false},
		{"Return error after invalid filter type", args{allBuckets, structs.Flags{FilterType: "harry"}}, nil, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := filterBuckets(tt.args.buckets, tt.args.flags)
			if (err != nil) != tt.wantErr {
				t.Errorf("filterBuckets() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("filterBuckets() = %v, want %v", got, tt.want)
			}
		})
	}
}

func (m *mockS3Client) GetBucketLifecycleConfiguration(input *s3.GetBucketLifecycleConfigurationInput) (*s3.GetBucketLifecycleConfigurationOutput, error) {
	var lifecycleExpiration = s3.LifecycleExpiration{}
	lifecycleExpiration.SetDays(2)

	var rule = s3.LifecycleRule{}
	rule.SetExpiration(&lifecycleExpiration)

	var lifeCycleOutput = s3.GetBucketLifecycleConfigurationOutput{}
	if *input.Bucket == "hasLifeCycle" {
		lifeCycleOutput.SetRules([]*s3.LifecycleRule{&rule})
		return &lifeCycleOutput, nil
	}

	return &lifeCycleOutput, errors.New("Fail")
}

func TestGetBucketLifeCycle(t *testing.T) {
	mockS3Client := &mockS3Client{}
	var lifeCycleRules = []*s3.LifecycleRule{
		{
			Expiration: &s3.LifecycleExpiration{
				Days: aws.Int64(2),
			},
		},
	}

	type args struct {
		client     s3iface.S3API
		bucketName *string
	}

	tests := []struct {
		name    string
		args    args
		want    []*s3.LifecycleRule
		wantErr bool
	}{
		{"Return lifeCycle rules", args{mockS3Client, aws.String("hasLifeCycle")}, lifeCycleRules, false},
		{"Return error", args{mockS3Client, aws.String("noLifeCycle")}, nil, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetBucketLifeCycle(tt.args.client, tt.args.bucketName)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetBucketLifeCycle() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetBucketLifeCycle() = %v, want %v", got, tt.want)
			}
		})
	}
}

func (m *mockS3Client) GetBucketAcl(input *s3.GetBucketAclInput) (*s3.GetBucketAclOutput, error) {
	var displayName = "harry"
	var permission = "FULL_CONTROL"

	var grantee = s3.Grantee{}
	grantee.SetDisplayName(displayName)

	var grants = s3.Grant{}
	grants.SetPermission(permission)
	grants.SetGrantee(&grantee)

	var ACLOutput = s3.GetBucketAclOutput{}
	if *input.Bucket == "hasGrants" {
		ACLOutput.SetGrants([]*s3.Grant{&grants})
		return &ACLOutput, nil
	}

	return &ACLOutput, errors.New("Fail")
}

func TestGetBucketACL(t *testing.T) {
	mockS3Client := &mockS3Client{}
	var grants = []*s3.Grant{
		{
			Grantee: &s3.Grantee{
				DisplayName: aws.String("harry"),
			},
			Permission: aws.String("FULL_CONTROL"),
		},
	}

	type args struct {
		client     s3iface.S3API
		bucketName *string
	}

	tests := []struct {
		name    string
		args    args
		want    []*s3.Grant
		wantErr bool
	}{
		{"Return all grants", args{mockS3Client, aws.String("hasGrants")}, grants, false},
		{"Return error", args{mockS3Client, aws.String("noGrants")}, nil, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetBucketACL(tt.args.client, tt.args.bucketName)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetBucketACL() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetBucketACL() = %v, want %v", got, tt.want)
			}
		})
	}
}

func (m *mockS3Client) GetBucketEncryption(input *s3.GetBucketEncryptionInput) (*s3.GetBucketEncryptionOutput, error) {
	var encryptionByDefault = s3.ServerSideEncryptionByDefault{}
	encryptionByDefault.SetSSEAlgorithm("harry")

	var encryptionRule = s3.ServerSideEncryptionRule{}
	encryptionRule.SetApplyServerSideEncryptionByDefault(&encryptionByDefault)

	var encryptionConfiguration = s3.ServerSideEncryptionConfiguration{}
	encryptionConfiguration.SetRules([]*s3.ServerSideEncryptionRule{&encryptionRule})

	var encryptionOutput = s3.GetBucketEncryptionOutput{}
	if *input.Bucket == "hasEncryption" {
		encryptionOutput.SetServerSideEncryptionConfiguration(&encryptionConfiguration)
		return &encryptionOutput, nil
	}

	return &encryptionOutput, errors.New("Fail")
}

func TestGetBucketEncryption(t *testing.T) {
	mockS3Client := &mockS3Client{}
	var serverSideEncryptionRule = []*s3.ServerSideEncryptionRule{
		{
			ApplyServerSideEncryptionByDefault: &s3.ServerSideEncryptionByDefault{
				SSEAlgorithm: aws.String("harry"),
			},
		},
	}

	type args struct {
		client     s3iface.S3API
		bucketName *string
	}

	tests := []struct {
		name    string
		args    args
		want    []*s3.ServerSideEncryptionRule
		wantErr bool
	}{
		{"Return all bucket Encryption rules", args{mockS3Client, aws.String("hasEncryption")}, serverSideEncryptionRule, false},
		{"Return error", args{mockS3Client, aws.String("noEncryption")}, nil, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetBucketEncryption(tt.args.client, tt.args.bucketName)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetBucketEncryption() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetBucketEncryption() = %v, want %v", got, tt.want)
			}
		})
	}
}

func (m *mockS3Client) GetBucketLocation(input *s3.GetBucketLocationInput) (*s3.GetBucketLocationOutput, error) {
	var locationOutput = s3.GetBucketLocationOutput{}

	if *input.Bucket == "withoutErr" {
		locationOutput.SetLocationConstraint(s3.BucketLocationConstraintSaEast1)
		return &locationOutput, nil
	}

	return &locationOutput, errors.New("Fail")
}

func TestGetBucketLocation(t *testing.T) {
	mockS3Client := &mockS3Client{}

	type args struct {
		client     s3iface.S3API
		bucketName *string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{"Return where bucket is stored", args{mockS3Client, aws.String("withoutErr")}, "sa-east-1", false},
		{"Return error", args{mockS3Client, aws.String("locationSP")}, "", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetBucketLocation(tt.args.client, tt.args.bucketName)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetBucketLocation() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GetBucketLocation() = %v, want %v", got, tt.want)
			}
		})
	}
}

func (m *mockS3Client) GetBucketTagging(input *s3.GetBucketTaggingInput) (*s3.GetBucketTaggingOutput, error) {
	var tagSet = s3.Tag{}
	tagSet.SetKey("name")
	tagSet.SetValue("harry")

	var taggingOutput = s3.GetBucketTaggingOutput{}

	if *input.Bucket == "withoutErr" {
		taggingOutput.SetTagSet([]*s3.Tag{&tagSet})
		return &taggingOutput, nil
	}

	return &taggingOutput, errors.New("Fail")
}

func TestGetBucketTagging(t *testing.T) {
	mockS3Client := &mockS3Client{}
	var tags = []*s3.Tag{
		{
			Key:   aws.String("name"),
			Value: aws.String("harry"),
		},
	}

	type args struct {
		client     s3iface.S3API
		bucketName *string
	}

	tests := []struct {
		name    string
		args    args
		want    []*s3.Tag
		wantErr bool
	}{
		{"Return where bucket is stored", args{mockS3Client, aws.String("withoutErr")}, tags, false},
		{"Return error", args{mockS3Client, aws.String("locationSP")}, nil, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetBucketTagging(tt.args.client, tt.args.bucketName)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetBucketTagging() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetBucketTagging() = %v, want %v", got, tt.want)
			}
		})
	}
}
