package errorchecker

import "log"

// CheckFatal Receive one error object and check if something go wrong
// and stop the execution
func CheckFatal(e error, functionName string) {
	if e != nil {
		log.Fatalf("Function where error appear: %s\nError: %v", functionName, e)
	}
}

// CheckError Receive one error object and check if something
// go wrong, log the error, but don`t stop the execution
func CheckError(e error, functionName string) {
	if e != nil {
		log.Printf("Function where error appear: %s\nError: %v", functionName, e)
	}
}
