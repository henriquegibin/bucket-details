package errorchecker

import (
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go/aws/awserr"
)

// CheckFatal Receive one error object and check if something go wrong
// and stop the execution
func CheckFatal(e error, functionName string) {
	if e != nil {
		log.Fatalf("Function where error appear: %s\nError: %v", functionName, e)
	}
}

// CheckError Receive one error object and check if something
// go wrong, log the error, but don`t stop the execution
func CheckError(e error, functionName string, debug bool) {
	if !debug {
		return
	}

	if e != nil {
		log.Printf("Function where error appear: %s\nError: %v", functionName, e)
	}
}

// AWSErrors Receive errors from aws sdk
func AWSErrors(err error, debug bool) {
	if !debug {
		return
	}

	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			fmt.Println("error", aerr.Code())
		}
	}
}
