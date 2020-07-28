package errorchecker

import (
	"errors"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
)

func TestCheckError(t *testing.T) {
	type args struct {
		e            error
		functionName string
		debug        bool
	}

	tests := []struct {
		name string
		args args
	}{
		{"Check if error is not nil and log the message if debug is true", args{errors.New("failed"), "functionName", true}},
		{"Do nothing if debug is false", args{errors.New("failed"), "functionName", false}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			CheckError(tt.args.e, tt.args.functionName, tt.args.debug)
		})
	}
}

func TestAWSErrors(t *testing.T) {
	err := aws.ErrMissingRegion
	type args struct {
		err   error
		debug bool
	}

	tests := []struct {
		name string
		args args
	}{
		{"Print errors from aws SDK", args{err, true}},
		{"Do nothing if debug is false", args{err, false}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			AWSErrors(tt.args.err, tt.args.debug)
		})
	}
}
