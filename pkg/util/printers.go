package util

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/jedib0t/go-pretty/v6/table"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
	"os"
	"streammachine.io/strm/pkg/common"
)

const OutputFormatFlag = "output"

type Printer interface {
	Print(proto proto.Message)
}

type GenericJsonPrinter struct{}

func (p GenericJsonPrinter) Print(proto proto.Message) {
	// As protojson.Marshal adds random spaces, we use json.Compact to omit the random spaces in the output.
	// Linked issue in google/protobuf: https://github.com/golang/protobuf/issues/1082
	marshal, _ := protojson.Marshal(proto)
	buffer := bytes.Buffer{}
	err := json.Compact(&buffer, marshal)

	if err != nil {
		common.CliExit(err)
	}

	fmt.Println(string(buffer.Bytes()))
}

func RenderTable(headers table.Row, rows []table.Row) {
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(headers)
	t.AppendSeparator()
	t.AppendRows(rows)
	t.Render()
}
