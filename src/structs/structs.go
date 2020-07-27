package structs

import (
	"time"
)

// BucketInfo struct represent all bucket info retrieved from bucket
type BucketInfo struct {
	Name                       string
	CreationDate               time.Time
	FilesCount                 int
	Size                       int64
	LastModifiedFromNewestFile time.Time
	Cost                       string
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
}
