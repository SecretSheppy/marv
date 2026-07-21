package pathutil

import (
	"reflect"
	"testing"
)

func TestSplit(t *testing.T) {
	split := Split("/user/test/dir/here")
	if !reflect.DeepEqual(split, []string{"user", "test", "dir", "here"}) {
		t.Errorf("Split did not produce expected output")
	}
}
