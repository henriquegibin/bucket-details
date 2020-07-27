package genericfunctions

import (
	"bucket-details/src/structs"
	"fmt"
	"time"
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

// GetFirstLastDayOfMonth Receive one string containing first or last.
//
// Return a time object with first or last day of the current month
func GetFirstLastDayOfMonth(day string) time.Time {
	if day == "first" {
		y, m, _ := time.Now().Date()
		d := time.Date(y, m, 1, 0, 0, 0, 0, time.UTC)
		return d
	}

	firstDay := GetFirstLastDayOfMonth("first")
	lastDay := firstDay.AddDate(0, 1, 0).Add(-time.Millisecond)
	return lastDay
}

// FlagsStructCreator receive all flags values and return one flag struct
func FlagsStructCreator(flags ...string) structs.Flags { // Melhorar depois
	var flagsStruct structs.Flags
	flagsStruct.FilterType = flags[0]
	flagsStruct.FilterValue = flags[1]
	flagsStruct.LifeCycle = flags[2]
	flagsStruct.BucketACL = flags[3]
	flagsStruct.BucketEncryption = flags[4]
	flagsStruct.BucketLocation = flags[5]
	flagsStruct.BucketTagging = flags[6]
	return flagsStruct
}
