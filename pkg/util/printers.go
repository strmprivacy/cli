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
	"streammachine.io/strm/pkg/common"
	"streammachine.io/strm/pkg/constants"
)

const OutputFormatFlag = "output"

type Printer interface {
	Print(proto proto.Message)
}

var DefaultPrinters = map[string]Printer{
	constants.OutputFormatJson + constants.ListCommandName:      GenericPrettyJsonPrinter{},
	constants.OutputFormatJson + constants.GetCommandName:       GenericPrettyJsonPrinter{},
	constants.OutputFormatJson + constants.DeleteCommandName:    GenericPrettyJsonPrinter{},
	constants.OutputFormatJson + constants.CreateCommandName:    GenericPrettyJsonPrinter{},
	constants.OutputFormatJsonRaw + constants.ListCommandName:   GenericRawJsonPrinter{},
	constants.OutputFormatJsonRaw + constants.GetCommandName:    GenericRawJsonPrinter{},
	constants.OutputFormatJsonRaw + constants.DeleteCommandName: GenericRawJsonPrinter{},
	constants.OutputFormatJsonRaw + constants.CreateCommandName: GenericRawJsonPrinter{},
}

type GenericRawJsonPrinter struct{}
type GenericPrettyJsonPrinter struct{}

func (p GenericRawJsonPrinter) Print(proto proto.Message) {
	rawJson := protoMessageToRawJson(proto)
	fmt.Println(string(rawJson.Bytes()))
}

func (p GenericPrettyJsonPrinter) Print(proto proto.Message) {
	prettyJson := protoMessageToPrettyJson(proto)
	fmt.Println(string(prettyJson.Bytes()))
}

func protoMessageToPrettyJson(proto proto.Message) bytes.Buffer {
	rawJson := protoMessageToRawJson(proto)
	prettyJson := bytes.Buffer{}

	errIndent := json.Indent(&prettyJson, rawJson.Bytes(), "", "    ")
	common.CliExit(errIndent)

	return prettyJson
}

func protoMessageToRawJson(proto proto.Message) bytes.Buffer {
	// As protojson.Marshal adds random spaces, we use json.Compact to omit the random spaces in the output.
	// Linked issue in google/protobuf: https://github.com/golang/protobuf/issues/1082
	marshal, _ := protojson.Marshal(proto)
	buffer := bytes.Buffer{}

	errCompact := json.Compact(&buffer, marshal)
	common.CliExit(errCompact)

	return buffer
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
