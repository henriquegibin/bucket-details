package structs

import (
	"time"

	"github.com/aws/aws-sdk-go/service/s3"
)

// BucketInfo struct represent all bucket info retrieved from bucket
type BucketInfo struct {
	Name                       string
	CreationDate               time.Time
	FilesCount                 int
	Size                       int64
	LastModifiedFromNewestFile time.Time
	Cost                       string
	Extras                     Extras
}

// Flags struct represent all flags received from the user
type Flags struct {
	FilterType       string
	FilterValue      string
	LifeCycle        bool
	BucketACL        bool
	BucketEncryption bool
	BucketLocation   bool
	BucketTagging    bool
	Debug            bool
}

// Extras struct represent all aditional data asked via flags
type Extras struct {
	LifeCycle        []*s3.LifecycleRule
	BucketACL        []*s3.Grant
	BucketEncryption []*s3.ServerSideEncryptionRule
	BucketLocation   string
	BucketTagging    []*s3.Tag
}
