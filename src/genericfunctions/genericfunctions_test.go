package genericfunctions

import (
	errorchecker "bucket-details/src/errorcheck"
	"bucket-details/src/structs"
	"io/ioutil"
	"os"
	"reflect"
	"testing"
	"time"
)

func TestBucketSize(t *testing.T) {
	size1 := int64(200)
	size2 := int64(300)
	size3 := int64(500)

	type args struct {
		array []int64
	}
	tests := []struct {
		name string
		args args
		want int64
	}{
		{"Sum all itens inside array", args{[]int64{size1, size2, size3}}, int64(1)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := BucketSize(tt.args.array); got != tt.want {
				t.Errorf("BucketSize() = %v, want %v", got, tt.want)
			}
		})
	}
}

// If you need to change something in the fixture file, pay attention.
// Depending of the configurations in your editor, you can accidentally create
// diffs, For exemplo, if your editor remove empty spaces at the end of a line.
func TestPrint(t *testing.T) {
	var extras structs.Extras
	rescueStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	date := time.Date(2020, time.July, 23, 22, 41, 0, 0, time.UTC)
	infos := structs.BucketInfo{
		Name:                       "bucket-name",
		CreationDate:               date,
		FilesCount:                 130,
		Size:                       1024,
		LastModifiedFromNewestFile: date,
		Cost:                       "10.00",
		Extras:                     extras,
	}
	Print(infos)

	w.Close()
	out, _ := ioutil.ReadAll(r)
	os.Stdout = rescueStdout

	content, err := ioutil.ReadFile("../../test/fixtures/printOutput.txt")
	errorchecker.CheckError(err, "TestPrint")

	if string(out) != string(content) {
		t.Errorf("Expected %s, got %s", string(content), out)
	}
}

func TestFlagsStructCreator(t *testing.T) {
	var flags = structs.Flags{
		FilterType:       "prefix",
		FilterValue:      "harry",
		LifeCycle:        true,
		BucketACL:        true,
		BucketEncryption: true,
		BucketLocation:   true,
		BucketTagging:    true,
	}
	var arrayInput = []string{
		"prefix",
		"harry",
		"true",
		"true",
		"true",
		"true",
		"true",
	}

	type args struct {
		flags []string
	}

	tests := []struct {
		name string
		args args
		want structs.Flags
	}{
		{"Return a struct flags", args{arrayInput}, flags},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := FlagsStructCreator(tt.args.flags...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FlagsStructCreator() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_parseBool(t *testing.T) {
	type args struct {
		flag string
	}

	tests := []struct {
		name string
		args args
		want bool
	}{
		{"Parse a string into bool type", args{"true"}, true},
		{"Parse a string into bool type", args{"false"}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := parseBool(tt.args.flag); got != tt.want {
				t.Errorf("parseBool() = %v, want %v", got, tt.want)
			}
		})
	}
}
