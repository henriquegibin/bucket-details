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
	}

	tests := []struct {
		name string
		args args
	}{
		{"Check if error is not nil and log the message", args{errors.New("failed"), "functionName"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			CheckError(tt.args.e, tt.args.functionName)
		})
	}
}

func TestAWSErrors(t *testing.T) {
	err := aws.ErrMissingRegion
	type args struct {
		err error
	}

	tests := []struct {
		name string
		args args
	}{
		{"Print errors from aws SDK", args{err}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			AWSErrors(tt.args.err)
		})
	}
}
