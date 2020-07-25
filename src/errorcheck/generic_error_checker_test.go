package errorchecker

import (
	"errors"
	"testing"
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
