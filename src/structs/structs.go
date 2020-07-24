package structs

import "time"

// BucketInfo struct represent all bucket info retrieved from bucket
type BucketInfo struct {
	Name                       string
	CreationDate               time.Time
	FilesCount                 int
	Size                       int64
	LastModifiedFromNewestFile string
	Cost                       string
}
