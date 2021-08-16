package stream

import (
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/spf13/cobra"
	"github.com/streammachineio/api-definitions-go/api/streams/v1"
	"google.golang.org/protobuf/proto"
	"streammachine.io/strm/pkg/constants"
	"streammachine.io/strm/pkg/util"
)

var printer util.Printer

func configurePrinter(command *cobra.Command) util.Printer {
	outputFormat := util.GetStringAndErr(command.Flags(), util.OutputFormatFlag)

	switch outputFormat {
	case "json":
		return util.GenericJsonPrinter{}
	case "table":
		switch command.Parent().Name() {
		case constants.ListCommandName:
			return ListStreamsTablePrinter{}
		}
		return util.GenericJsonPrinter{}
	default:
		return util.GenericJsonPrinter{}
	}
}

type ListStreamsTablePrinter struct{}

func (p ListStreamsTablePrinter) Print(data proto.Message) {
	streamsResponse, _ := (data).(*streams.ListStreamsResponse)

	rows := make([]table.Row, 0, len(streamsResponse.Streams))

	for _, stream := range streamsResponse.Streams {
		rows = append(rows, table.Row{stream.Stream.Ref.Name, stream.Stream.Enabled, stream.KeyStream != nil})
	}

	util.RenderTable(
		table.Row{"Name", "Enabled", "Has Key Stream"},
		rows,
	)
}
