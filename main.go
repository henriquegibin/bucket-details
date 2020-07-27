package main

import (
	"bucket-details/src/aws"
	errorchecker "bucket-details/src/errorcheck"
	"bucket-details/src/genericfunctions"
	"bucket-details/src/structs"
	"os"
	"time"

	"github.com/aws/aws-sdk-go/service/costexplorer"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "Bucket Details"
	app.Usage = "Use this CLI to retrieve metadata from all AWS Buckets. Use flags to specifie filters and more options"
	app.UsageText = "bucket-details [global options]"

	app.Flags = []cli.Flag{
		&cli.StringFlag{
			Name:  "filter-type",
			Usage: "Choose the method used to filter buckets. Default is without filter. Possible values are: prefix, regexp, suffix. This flag needs to be used with filterValue",
		},
		&cli.StringFlag{
			Name:  "filter-value",
			Usage: "Pass your string to use as filter. This flag needs to be used with filterType",
		},
		&cli.StringFlag{
			Name:  "life-cycle",
			Usage: "Pass this flag to retrieve the bucket life cycle",
			Value: "false",
		},
		&cli.StringFlag{
			Name:  "bucket-acl",
			Usage: "Pass this flag to retrieve the bucket bucket acl",
			Value: "false",
		},
		&cli.StringFlag{
			Name:  "bucket-encryption",
			Usage: "Pass this flag to retrieve the bucket encryption",
			Value: "false",
		},
		&cli.StringFlag{
			Name:  "bucket-location",
			Usage: "Pass this flag to retrieve the bucket location",
			Value: "false",
		},
		&cli.StringFlag{
			Name:  "bucket-tagging",
			Usage: "Pass this flag to retrieve the bucket tagging",
			Value: "false",
		},
	}

	app.Action = func(c *cli.Context) error {
		flags := genericfunctions.FlagsStructCreator(
			c.String("filter-type"),
			c.String("filter-value"),
			c.String("life-cycle"),
			c.String("bucket-acl"),
			c.String("bucket-encryption"),
			c.String("bucket-location"),
			c.String("bucket-tagging"),
		)
		getMetadata(flags)
		return nil
	}

	err := app.Run(os.Args)
	errorchecker.CheckFatal(err, "main")
}

func getMetadata(flags structs.Flags) {
	awsSession := aws.CreateNewAwsSession()
	s3Instance := s3.New(awsSession)
	costexplorerInstance := costexplorer.New(awsSession)
	var size []int64
	var filesCount = 0
	var lastModified time.Time

	buckets := aws.ListBuckets(s3Instance, flags)
	for _, bucket := range buckets {
		var bucketDetails structs.BucketInfo

		sizeEachObject, filesCount, lastModified := aws.ListObjects(bucket.Name, s3Instance, &size, &filesCount, &lastModified)
		bucketPrice := aws.CheckS3BucketCost(costexplorerInstance, *bucket.Name)

		bucketDetails.Name = *bucket.Name
		bucketDetails.CreationDate = *bucket.CreationDate
		bucketDetails.FilesCount = filesCount
		bucketDetails.Size = genericfunctions.BucketSize(sizeEachObject)
		bucketDetails.LastModifiedFromNewestFile = lastModified
		bucketDetails.Cost = bucketPrice
		genericfunctions.BeautyPrint(bucketDetails)
	}
}
