package genericfunctions

import (
	"bucket-details/src/structs"
	"fmt"
)

// BucketSize Receive int64 array then sum all itens inside.
//
// Return int64 value.
func BucketSize(array []int64) int64 {
	result := int64(0)
	for _, v := range array {
		result += v
	}

	result = result / 1000
	return result
}

// BeautyPrint Receive one BucketInfo object then print every property
func BeautyPrint(bucketDetails structs.BucketInfo) {
	fmt.Println("----------------")
	fmt.Printf("Bucket Name:                           %s\n", bucketDetails.Name)
	fmt.Printf("Bucket Creation:                       %s\n", bucketDetails.CreationDate)
	fmt.Printf("Bucket Files:                          %d\n", bucketDetails.FilesCount)
	fmt.Printf("Bucket Size:                           %d\n", bucketDetails.Size)
	fmt.Printf("Bucket Last Modified From Newest File: %s\n", bucketDetails.LastModifiedFromNewestFile)
	fmt.Printf("Bucket Cost:                           $%s\n", bucketDetails.Cost)
	fmt.Println("----------------")
}