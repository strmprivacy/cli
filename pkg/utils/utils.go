package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	flag "github.com/spf13/pflag"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"streammachine.io/strm/pkg/common"
)

var ConfigPath string

func Print(m proto.Message) {
	// As protojson.Marshal adds random spaces, we use json.Compact to omit the random spaces in the output.
	// Linked issue in google/protobuf: https://github.com/golang/protobuf/issues/1082
	marshal, _ := protojson.Marshal(m)
	buffer := bytes.Buffer{}
	err := json.Compact(&buffer, marshal)

	if err != nil {
		common.CliExit(err)
	}

	fmt.Println(string(buffer.Bytes()))
}

func MapInt32ToString(vs []int32, f func(a ...interface{}) string) []string {
	vsm := make([]string, len(vs))
	for i, v := range vs {
		vsm[i] = f(v)
	}
	return vsm
}

func atoi(s string) int {
	v, err := strconv.Atoi(s)
	common.CliExit(err)
	return v
}
func atoi32(s string) int32 {
	return int32(atoi(s))
}

func MapStringsToInt(vs []string, f func(string) int) []int {
	vsm := make([]int, len(vs))
	for i, v := range vs {
		vsm[i] = f(v)
	}
	return vsm
}
func MapStringsToInt32(vs []string, f func(string) int32) []int32 {
	vsm := make([]int32, len(vs))
	for i, v := range vs {
		vsm[i] = f(v)
	}
	return vsm
}

func StringsArrayToInt(vs []string) []int {
	return MapStringsToInt(vs, atoi)
}
func StringsArrayToInt32(vs []string) []int32 {
	return MapStringsToInt32(vs, atoi32)
}

func GetStringAndErr(f *flag.FlagSet, k string) string {
	v, err := f.GetString(k)
	common.CliExit(err)
	return v
}
func GetBoolAndErr(f *flag.FlagSet, k string) bool {
	v, err := f.GetBool(k)
	common.CliExit(err)
	return v
}
func GetInt32AndErr(f *flag.FlagSet, k string) int32 {
	v, err := f.GetInt32(k)
	common.CliExit(err)
	return v
}
func GetInt64AndErr(f *flag.FlagSet, k string) int64 {
	v, err := f.GetInt64(k)
	common.CliExit(err)
	return v
}
func GetIntAndErr(f *flag.FlagSet, k string) int {
	v, err := f.GetInt(k)
	common.CliExit(err)
	return v
}

func TryLoad(m proto.Message, name *string) error {
	filename := getSaveFilename(m, name)
	bytes, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}
	err = protojson.Unmarshal(bytes, m)
	return err
}

func Load(m proto.Message, name *string) {
	err := TryLoad(m, name)
	common.CliExit(err)
}

func Save(m proto.Message, name *string) {
	filename := getSaveFilename(m, name)
	err := os.MkdirAll(filepath.Dir(filename), 0700)
	common.CliExit(err)
	jsonData, err := protojson.Marshal(m)
	common.CliExit(err)
	err = ioutil.WriteFile(filename, jsonData, 0644)
	common.CliExit(err)
}

func DeleteSaved(m proto.Message, name *string) {
	filename := getSaveFilename(m, name)
	_ = os.Remove(filename)
}

func getSaveFilename(m proto.Message, name *string) string {
	cat := fmt.Sprint(m.ProtoReflect().Descriptor().Name())
	return path.Join(ConfigPath, cat, *name+".json")
}
