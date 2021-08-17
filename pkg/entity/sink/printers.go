package sink

import (
	"fmt"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/spf13/cobra"
	"github.com/streammachineio/api-definitions-go/api/entities/v1"
	"github.com/streammachineio/api-definitions-go/api/sinks/v1"
	"google.golang.org/protobuf/proto"
	"streammachine.io/strm/pkg/common"
	"streammachine.io/strm/pkg/constants"
	"streammachine.io/strm/pkg/util"
)

var printer util.Printer

func configurePrinter(command *cobra.Command) util.Printer {
	outputFormat := util.GetStringAndErr(command.Flags(), util.OutputFormatFlag)

	var recursive bool

	if command.Parent().Name() != constants.CreateCommandName {
		recursive = util.GetBoolAndErr(command.Flags(), constants.RecursiveFlagName)
	}

	switch outputFormat {
	case constants.OutputFormatJson:
		return util.GenericPrettyJsonPrinter{}
	case constants.OutputFormatJsonRaw:
		return util.GenericRawJsonPrinter{}
	case constants.OutputFormatTable:
		switch command.Parent().Name() {
		case constants.ListCommandName:
			return listTablePrinter{recursive}
		case constants.GetCommandName:
			return getTablePrinter{recursive}
		case constants.DeleteCommandName:
			return deletePrinter{recursive}
		case constants.CreateCommandName:
			return createTablePrinter{}
		}

		return util.GenericPrettyJsonPrinter{}
	case constants.OutputFormatPlain:
		switch command.Parent().Name() {
		case constants.ListCommandName:
			return listPlainPrinter{}
		case constants.GetCommandName:
			return getPlainPrinter{}
		case constants.DeleteCommandName:
			return deletePrinter{recursive}
		case constants.CreateCommandName:
			return createPlainPrinter{}
		}

		return util.GenericPrettyJsonPrinter{}
	default:
		common.CliExit(fmt.Sprintf("Output format '%v' is not supported. Allowed values: %v", outputFormat, constants.OutputFormatFlagAllowedValuesText))
		return nil
	}
}

type listPlainPrinter struct{}
type getPlainPrinter struct{}
type createPlainPrinter struct{}

type listTablePrinter struct {
	recursive bool
}
type getTablePrinter struct {
	recursive bool
}
type createTablePrinter struct{}

type deletePrinter struct {
	recursive bool
}

func (p listTablePrinter) Print(data proto.Message) {
	listResponse, _ := (data).(*sinks.ListSinksResponse)
	printTable(listResponse.Sinks, p.recursive)
}

func (p getTablePrinter) Print(data proto.Message) {
	getResponse, _ := (data).(*sinks.GetSinkResponse)
	printTable([]*entities.SinkTree{getResponse.SinkTree}, p.recursive)
}

func (p createTablePrinter) Print(data proto.Message) {
	createResponse, _ := (data).(*sinks.CreateSinkResponse)
	printTable([]*entities.SinkTree{{Sink: createResponse.Sink}}, false)
}

func (p listPlainPrinter) Print(data proto.Message) {
	listResponse, _ := (data).(*sinks.ListSinksResponse)
	printPlain(listResponse.Sinks)
}

func (p getPlainPrinter) Print(data proto.Message) {
	getResponse, _ := (data).(*sinks.GetSinkResponse)
	printPlain([]*entities.SinkTree{getResponse.SinkTree})
}

func (p createPlainPrinter) Print(data proto.Message) {
	createResponse, _ := (data).(*sinks.CreateSinkResponse)
	printPlain([]*entities.SinkTree{{Sink: createResponse.Sink}})
}

func (p deletePrinter) Print(_ proto.Message) {
	if p.recursive {
		fmt.Println("Sink and linked resources have been deleted")
	} else {
		fmt.Println("Sink has been deleted")
	}
}

func printTable(sinkTreeArray []*entities.SinkTree, recursive bool) {
	rows := make([]table.Row, 0, len(sinkTreeArray))

	for _, sink := range sinkTreeArray {
		var sinkBucketName string

		switch config := sink.Sink.Config.(type) {
		case *entities.Sink_Bucket:
			sinkBucketName = config.Bucket.BucketName
		}

		var row table.Row

		if recursive {
			row = table.Row{
				sink.Sink.Ref.Name,
				sink.Sink.SinkType.String(),
				sinkBucketName,
				len(sink.BatchExporters),
			}
		} else {
			row = table.Row{
				sink.Sink.Ref.Name,
				sink.Sink.SinkType.String(),
				sinkBucketName,
			}
		}

		rows = append(rows, row)
	}

	var header table.Row

	if recursive {
		header = table.Row{
			"Sink",
			"Type",
			"Bucket",
			"# Batch Exporters",
		}
	} else {
		header = table.Row{
			"Sink",
			"Type",
			"Bucket",
		}
	}

	util.RenderTable(header, rows)
}

func printPlain(sinkTreeArray []*entities.SinkTree) {
	var names string
	lastIndex := len(sinkTreeArray) - 1

	for index, sink := range sinkTreeArray {
		names = names + sink.Sink.Ref.Name

		if index != lastIndex {
			names = names + "\n"
		}
	}

	util.RenderPlain(names)
}
