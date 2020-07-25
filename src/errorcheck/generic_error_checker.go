package errorchecker

import "log"

// CheckFatal Receive one error object and check if something go wrong
// and stop the execution
func CheckFatal(e error) {
	if e != nil {
		log.Fatal(e)
	}
}

// CheckError Receive one error object and check if something
// go wrong, log the error, but don`t stop the execution
func CheckError(e error) {
	if e != nil {
		log.Println(e)
	}
}
