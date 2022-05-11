package util

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
	"os"
	"strings"
	"strmprivacy/strm/pkg/common"
)

type Printer interface {
	Print(data interface{})
}

var DefaultPrinters = map[string]Printer{
	common.OutputFormatJson + common.ListCommandName:      ProtoMessageJsonPrettyPrinter{},
	common.OutputFormatJson + common.GetCommandName:       ProtoMessageJsonPrettyPrinter{},
	common.OutputFormatJson + common.DeleteCommandName:    ProtoMessageJsonPrettyPrinter{},
	common.OutputFormatJson + common.CreateCommandName:    ProtoMessageJsonPrettyPrinter{},
	common.OutputFormatJsonRaw + common.ListCommandName:   ProtoMessageJsonRawPrinter{},
	common.OutputFormatJsonRaw + common.GetCommandName:    ProtoMessageJsonRawPrinter{},
	common.OutputFormatJsonRaw + common.DeleteCommandName: ProtoMessageJsonRawPrinter{},
	common.OutputFormatJsonRaw + common.CreateCommandName: ProtoMessageJsonRawPrinter{},
}

type ProtoMessageJsonRawPrinter struct{}
type ProtoMessageJsonPrettyPrinter struct{}

func (p ProtoMessageJsonRawPrinter) Print(content interface{}) {
	protoContent, _ := (content).(proto.Message)
	rawJson := protoMessageToRawJson(protoContent)
	fmt.Println(string(rawJson.Bytes()))
}

func (p ProtoMessageJsonPrettyPrinter) Print(content interface{}) {
	protoContent, _ := (content).(proto.Message)
	prettyJson := protoMessageToPrettyJson(protoContent)
	fmt.Println(string(prettyJson.Bytes()))
}

func protoMessageToPrettyJson(proto proto.Message) bytes.Buffer {
	return PrettifyJson(protoMessageToRawJson(proto))
}

func protoMessageToRawJson(proto proto.Message) bytes.Buffer {
	// As protojson.Marshal adds random spaces, we use json.Compact to omit the random spaces in the output.
	// Linked issue in google/protobuf: https://github.com/golang/protobuf/issues/1082
	marshal, _ := protojson.Marshal(proto)
	return CompactJson(marshal)
}

func CompactJson(rawJson []byte) bytes.Buffer {
	buffer := bytes.Buffer{}

	errCompact := json.Compact(&buffer, rawJson)
	common.CliExit(errCompact)
	return buffer
}

func PrettifyJson(rawJson bytes.Buffer) bytes.Buffer {
	prettyJson := bytes.Buffer{}

	errIndent := json.Indent(&prettyJson, rawJson.Bytes(), "", "    ")
	common.CliExit(errIndent)

	return prettyJson
}

func MergePrinterMaps(maps ...map[string]Printer) (result map[string]Printer) {
	result = make(map[string]Printer)

	for _, m := range maps {
		for k, v := range m {
			result[k] = v
		}
	}
	return result
}

func RenderPlain(text string) {
	if len(text) == 0 {
		fmt.Println("No entities of this resource type exist.")
	} else {
		fmt.Println(text)
	}
}

func RenderCsv(headers table.Row, rows []table.Row) {
	if len(rows) == 0 {
		fmt.Println("No usage in the provided time period.")
	} else {
		t := table.NewWriter()
		t.SetOutputMirror(os.Stdout)
		t.AppendHeader(headers)
		t.AppendRows(rows)
		t.RenderCSV()
	}
}

func RenderTable(headers table.Row, rows []table.Row) {
	if len(rows) == 0 {
		fmt.Println("No entities of this resource type exist.")
	} else {
		t := table.NewWriter()
		t.SetOutputMirror(os.Stdout)
		t.AppendHeader(headers)
		t.AppendSeparator()
		t.AppendRows(rows)
		t.SetStyle(noBordersStyle)
		t.Render()
	}
}

var noBordersStyle = table.Style{
	Name:    "StyleNoBorders",
	Options: table.OptionsNoBorders,
	Title:   table.TitleOptionsDefault,
	Format:  table.FormatOptionsDefault,
	Box: table.BoxStyle{
		BottomLeft:       " ",
		BottomRight:      " ",
		BottomSeparator:  " ",
		EmptySeparator:   text.RepeatAndTrim(" ", text.RuneCount(" ")),
		Left:             " ",
		LeftSeparator:    " ",
		MiddleHorizontal: " ",
		MiddleSeparator:  " ",
		MiddleVertical:   " ",
		PaddingLeft:      " ",
		PaddingRight:     " ",
		PageSeparator:    "\n",
		Right:            " ",
		RightSeparator:   " ",
		TopLeft:          " ",
		TopRight:         " ",
		TopSeparator:     " ",
		UnfinishedRow:    "...",
	},
}
func ReplaceConsent(consentLevelType string) string {
	return strings.Replace(consentLevelType, "CONSENT", "PURPOSE", -1)
}