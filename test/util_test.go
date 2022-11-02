package test

import (
	"reflect"
	"strmprivacy/strm/pkg/util"
	"testing"
)

func TestInt32map(t *testing.T) {
	in := []string{"1", "2", "3"}
	out := util.StringsArrayToInt32(in)
	if !reflect.DeepEqual([]int32{1, 2, 3}, out) {
		t.Fail()
	}
}

func TestStringMap(t *testing.T) {
	in := []string{"1", "2", "3"}
	out := util.MapStrings(in, func(i string) string {
		return i + "a"
	})
	if !reflect.DeepEqual([]string{"1a", "2a", "3a"}, out) {
		t.Fail()
	}
}
