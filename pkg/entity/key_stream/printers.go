package key_stream

import (
	"fmt"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/spf13/cobra"
	"github.com/streammachineio/api-definitions-go/api/entities/v1"
	"github.com/streammachineio/api-definitions-go/api/key_streams/v1"
	"google.golang.org/protobuf/proto"
	"streammachine.io/strm/pkg/common"
	"streammachine.io/strm/pkg/constants"
	"streammachine.io/strm/pkg/util"
)

var printer util.Printer

func configurePrinter(command *cobra.Command) util.Printer {
	outputFormat := util.GetStringAndErr(command.Flags(), util.OutputFormatFlag)

	p := availablePrinters()[outputFormat+command.Parent().Name()]

	if p == nil {
		common.CliExit(fmt.Sprintf("Output format '%v' is not supported. Allowed values: %v", outputFormat, constants.OutputFormatFlagAllowedValuesText))
	}

	return p
}

func availablePrinters() map[string]util.Printer {
	return util.MergePrinterMaps(
		util.DefaultPrinters,
		map[string]util.Printer{
			constants.OutputFormatTable + constants.ListCommandName: listTablePrinter{},
			constants.OutputFormatTable + constants.GetCommandName:  getTablePrinter{},
			constants.OutputFormatPlain + constants.ListCommandName: listPlainPrinter{},
			constants.OutputFormatPlain + constants.GetCommandName:  getPlainPrinter{},
		},
	)
}

type listPlainPrinter struct{}
type getPlainPrinter struct{}

type listTablePrinter struct{}
type getTablePrinter struct{}

func (p listTablePrinter) Print(data proto.Message) {
	listResponse, _ := (data).(*key_streams.ListKeyStreamsResponse)
	printTable(listResponse.KeyStreams)
}

func (p getTablePrinter) Print(data proto.Message) {
	getResponse, _ := (data).(*key_streams.GetKeyStreamResponse)
	printTable([]*entities.KeyStream{getResponse.KeyStream})
}

func (p listPlainPrinter) Print(data proto.Message) {
	listResponse, _ := (data).(*key_streams.ListKeyStreamsResponse)
	printPlain(listResponse.KeyStreams)
}

func (p getPlainPrinter) Print(data proto.Message) {
	getResponse, _ := (data).(*key_streams.GetKeyStreamResponse)
	printPlain([]*entities.KeyStream{getResponse.KeyStream})
}

func printTable(keyStreams []*entities.KeyStream) {
	rows := make([]table.Row, 0, len(keyStreams))

	for _, keyStream := range keyStreams {
		rows = append(rows, table.Row{
			keyStream.Ref.Name,
			keyStream.Status,
		})
	}

	util.RenderTable(
		table.Row{
			"Key Stream",
			"Status",
		},
		rows,
	)
}

func printPlain(keyStreams []*entities.KeyStream) {
	var names string
	lastIndex := len(keyStreams) - 1

	for index, stream := range keyStreams {
		names = names + stream.Ref.Name

		if index != lastIndex {
			names = names + "\n"
		}
	}

	util.RenderPlain(names)
}
