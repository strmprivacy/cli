package util

import (
	"fmt"
	"github.com/spf13/pflag"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"streammachine.io/strm/pkg/common"
	"streammachine.io/strm/pkg/constants"
	"streammachine.io/strm/pkg/demoschema"
)

func atoi(s string) int {
	v, err := strconv.Atoi(s)
	common.CliExit(err)
	return v
}
func atoi32(s string) int32 {
	return int32(atoi(s))
}

func MapStringsToInt32(vs []string, f func(string) int32) []int32 {
	if len(vs) == 0 {
		return []int32{}

	}
	vsm := make([]int32, len(vs))
	for i, v := range vs {
		vsm[i] = f(v)
	}
	return vsm
}

func StringsArrayToInt32(vs []string) []int32 {
	return MapStringsToInt32(vs, atoi32)
}

func GetStringAndErr(f *pflag.FlagSet, k string) string {
	v, err := f.GetString(k)
	common.CliExit(err)
	return v
}
func GetBoolAndErr(f *pflag.FlagSet, k string) bool {
	v, err := f.GetBool(k)
	common.CliExit(err)
	return v
}
func GetInt64AndErr(f *pflag.FlagSet, k string) int64 {
	v, err := f.GetInt64(k)
	common.CliExit(err)
	return v
}
func GetIntAndErr(f *pflag.FlagSet, k string) int {
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
	return path.Join(constants.ConfigPath, constants.SavedEntitiesDirectory, cat, *name+".json")
}

func CreateUnionString(s string) *demoschema.UnionNullString {
	v := demoschema.NewUnionNullString()
	v.UnionType = demoschema.UnionNullStringTypeEnumString
	v.String = s
	return v
}

func CreateConfigDirAndFileIfNotExists() {
	err := os.MkdirAll(filepath.Dir(constants.ConfigPath), 0700)
	common.CliExit(err)

	configFilepath := path.Join(constants.ConfigPath, constants.DefaultConfigFilename+constants.DefaultConfigFileSuffix)

	if _, err := os.Stat(configFilepath); os.IsNotExist(err) {
		writeFileError := ioutil.WriteFile(
			configFilepath,
			constants.DefaultConfigFileContents,
			0644,
		)

		common.CliExit(writeFileError)
	}
}
