package batch_exporter

import (
	"fmt"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/spf13/cobra"
	"github.com/streammachineio/api-definitions-go/api/batch_exporters/v1"
	v1 "github.com/streammachineio/api-definitions-go/api/entities/v1"
	"google.golang.org/protobuf/proto"
	"streammachine.io/strm/pkg/constants"
	"streammachine.io/strm/pkg/util"
)

var printer util.Printer

func configurePrinter(command *cobra.Command) util.Printer {
	outputFormat := util.GetStringAndErr(command.Flags(), util.OutputFormatFlag)

	switch outputFormat {
	case "json":
		return util.GenericPrettyJsonPrinter{}
	case "json-raw":
		return util.GenericRawJsonPrinter{}
	case "table":
		switch command.Parent().Name() {
		case constants.ListCommandName:
			return listTablePrinter{}
		case constants.GetCommandName:
			return getTablePrinter{}
		case constants.DeleteCommandName:
			return deletePrinter{}
		case constants.CreateCommandName:
			return createTablePrinter{}
		}

		return util.GenericPrettyJsonPrinter{}
	case "plain":
		switch command.Parent().Name() {
		case constants.ListCommandName:
			return listPlainPrinter{}
		case constants.GetCommandName:
			return getPlainPrinter{}
		case constants.DeleteCommandName:
			return deletePrinter{}
		case constants.CreateCommandName:
			return createPlainPrinter{}
		}

		return util.GenericPrettyJsonPrinter{}
	default:
		return util.GenericPrettyJsonPrinter{}
	}
}

type listPlainPrinter struct{}
type getPlainPrinter struct{}
type createPlainPrinter struct{}

type listTablePrinter struct{}
type getTablePrinter struct{}
type createTablePrinter struct{}

type deletePrinter struct{}

func (p listTablePrinter) Print(data proto.Message) {
	listResponse, _ := (data).(*batch_exporters.ListBatchExportersResponse)
	printTable(listResponse.BatchExporters)
}

func (p getTablePrinter) Print(data proto.Message) {
	getResponse, _ := (data).(*batch_exporters.GetBatchExporterResponse)
	printTable([]*v1.BatchExporter{getResponse.BatchExporter})
}

func (p createTablePrinter) Print(data proto.Message) {
	createResponse, _ := (data).(*batch_exporters.CreateBatchExporterResponse)
	printTable([]*v1.BatchExporter{createResponse.BatchExporter})
}

func (p listPlainPrinter) Print(data proto.Message) {
	listResponse, _ := (data).(*batch_exporters.ListBatchExportersResponse)
	printPlain(listResponse.BatchExporters)
}

func (p getPlainPrinter) Print(data proto.Message) {
	getResponse, _ := (data).(*batch_exporters.GetBatchExporterResponse)
	printPlain([]*v1.BatchExporter{getResponse.BatchExporter})
}

func (p createPlainPrinter) Print(data proto.Message) {
	createResponse, _ := (data).(*batch_exporters.CreateBatchExporterResponse)
	printPlain([]*v1.BatchExporter{createResponse.BatchExporter})
}

func (p deletePrinter) Print(_ proto.Message) {
	fmt.Println("Batch Exporter has been deleted")
}

func printTable(batchExporters []*v1.BatchExporter) {
	rows := make([]table.Row, 0, len(batchExporters))

	for _, batchExporter := range batchExporters {

		var streamOrKeyStreamName string

		switch ref := batchExporter.StreamOrKeyStreamRef.(type) {
		case *v1.BatchExporter_KeyStreamRef:
			streamOrKeyStreamName = ref.KeyStreamRef.Name + " (key stream)"
		case *v1.BatchExporter_StreamRef:
			streamOrKeyStreamName = ref.StreamRef.Name
		}

		rows = append(rows, table.Row{
			batchExporter.Ref.Name,
			streamOrKeyStreamName,
			batchExporter.SinkName,
			batchExporter.Interval,
			batchExporter.PathPrefix,
		})
	}

	util.RenderTable(
		table.Row{
			"Batch Exporter",
			"Stream",
			"Sink",
			"Interval",
			"Path Prefix",
		},
		rows,
	)
}

func printPlain(batchExporters []*v1.BatchExporter) {
	var names string
	lastIndex := len(batchExporters) - 1

	for index, batchExporter := range batchExporters {
		names = names + batchExporter.Ref.Name

		if index != lastIndex {
			names = names + "\n"
		}
	}

	fmt.Println(names)
}
