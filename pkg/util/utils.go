package util

import (
	"fmt"
	"github.com/bykof/gostradamus"
	"github.com/lithammer/dedent"
	"github.com/samber/lo"
	"github.com/spf13/pflag"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"strmprivacy/strm/pkg/common"
	"time"
)

func atoi(s string) int {
	v, err := strconv.Atoi(s)
	common.CliExit(err)
	return v
}
func atoi32(s string) int32 {
	return int32(atoi(s))
}

func MapStrings[T any](vs []string, f func(string) T) []T {
	return lo.Map[string, T](vs, func(s string, _ int) T {
		return f(s)
	})
}
func StringsArrayToInt32(vs []string) []int32 {
	return MapStrings[int32](vs, atoi32)
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

func GetInt32AndErr(f *pflag.FlagSet, k string) int32 {
	v, err := f.GetInt32(k)
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
	bytes, err := os.ReadFile(filename)
	if err != nil {
		return err
	}
	err = protojson.UnmarshalOptions{DiscardUnknown: true}.Unmarshal(bytes, m)
	return err
}

func Save(m proto.Message, name *string) {
	filename := getSaveFilename(m, name)
	err := os.MkdirAll(filepath.Dir(filename), 0700)
	common.CliExit(err)
	jsonData, err := protojson.Marshal(m)
	common.CliExit(err)
	err = os.WriteFile(filename, jsonData, 0644)
	common.CliExit(err)
}

func DeleteSaved(m proto.Message, name *string) {
	filename := getSaveFilename(m, name)
	_ = os.Remove(filename)
}

func getSaveFilename(m proto.Message, name *string) string {
	cat := fmt.Sprint(m.ProtoReflect().Descriptor().Name())
	return path.Join(common.ConfigPath, common.SavedEntitiesDirectory, cat, *name+".json")
}

func CreateConfigDirAndFileIfNotExists() {
	err := os.MkdirAll(filepath.Dir(common.ConfigPath), 0700)
	common.CliExit(err)

	configFilepath := path.Join(common.ConfigPath, common.DefaultConfigFilename+common.DefaultConfigFileSuffix)

	if _, err := os.Stat(configFilepath); os.IsNotExist(err) {
		writeFileError := os.WriteFile(
			configFilepath,
			common.DefaultConfigFileContents,
			0644,
		)

		common.CliExit(writeFileError)
	}
}

// LongDocs dedents, trims surrounding whitespace, changes !strm for the command Name and changes ° for `
func LongDocs(s string) string {
	s2 := DedentTrim(strings.Replace(
		strings.Replace(s, "!strm", common.RootCommandName, -1), "°", "`", -1))
	return s2
}

func LongDocsUsage(s string) string {
	return LongDocs(s) + "\n\n### Usage"
}

func DedentTrim(s string) string {
	return strings.TrimSpace(dedent.Dedent(s))

}

func IsoFormat(tz gostradamus.Timezone, t *timestamppb.Timestamp) string {
	tt := time.Unix(t.Seconds, int64(t.Nanos))
	n := gostradamus.DateTimeFromTime(tt)
	return n.InTimezone(tz).IsoFormatTZ()
}
