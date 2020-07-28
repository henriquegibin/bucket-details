package genericfunctions

import (
	errorchecker "bucket-details/src/errorcheck"
	"bucket-details/src/structs"
	"encoding/json"
	"fmt"
	"strconv"
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

// Print Receive one BucketInfo object then print as json
func Print(bucketDetails structs.BucketInfo) {
	json, err := json.Marshal(bucketDetails)
	if err != nil {
		errorchecker.CheckError(err, "Print")
	}

	fmt.Println(string(json))
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
	flagsStruct.LifeCycle = parseBool(flags[2])
	flagsStruct.BucketACL = parseBool(flags[3])
	flagsStruct.BucketEncryption = parseBool(flags[4])
	flagsStruct.BucketLocation = parseBool(flags[5])
	flagsStruct.BucketTagging = parseBool(flags[6])
	return flagsStruct
}

// parseBool convert string into bool
func parseBool(flag string) bool {
	res, err := strconv.ParseBool(flag)
	errorchecker.CheckFatal(err, "parseBool")
	return res
}
