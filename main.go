package main

import (
	"bucket-details/src/aws"
	"bucket-details/src/genericfunctions"
	"bucket-details/src/structs"
	"log"
	"os"

	"github.com/aws/aws-sdk-go/service/costexplorer"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "Bucket Details"
	app.Usage = "Use this CLI to retreive metadata from all AWS Buckets. Use flags to specifie filters and more options"

	app.Flags = []cli.Flag{
		&cli.StringFlag{
			Name:  "lang",
			Value: "english",
			Usage: "language for the greeting",
		},
	}

	app.Action = func(c *cli.Context) error {
		getMetadata()
		// var output string
		// if c.String("lang") == "spanish" {
		// 	output = "Hola"
		// } else {
		// 	output = "Hello"
		// }
		// fmt.Println(output)
		return nil
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func getMetadata() {
	awsSession := aws.CreateNewAwsSession()
	s3Instance := s3.New(awsSession)
	costexplorerInstance := costexplorer.New(awsSession)
	var size []int64
	var filesCount = 0

	buckets := aws.ListBuckets(s3Instance)
	for _, bucket := range buckets {
		var bucketDetails structs.BucketInfo

		sizeEachObject, filesCount := aws.ListObjects(bucket.Name, s3Instance, &size, &filesCount)
		bucketPrice := aws.CheckPrice(costexplorerInstance, *bucket.Name)

		bucketDetails.Name = *bucket.Name
		bucketDetails.CreationDate = *bucket.CreationDate
		bucketDetails.FilesCount = filesCount
		bucketDetails.Size = genericfunctions.BucketSize(sizeEachObject)
		bucketDetails.Cost = bucketPrice
		genericfunctions.BeautyPrint(bucketDetails)
	}
}
