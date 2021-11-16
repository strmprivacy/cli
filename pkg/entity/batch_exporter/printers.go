package batch_exporter

import (
	"fmt"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/spf13/cobra"
	"github.com/strmprivacy/api-definitions-go/v2/api/batch_exporters/v1"
	v1 "github.com/strmprivacy/api-definitions-go/v2/api/entities/v1"
	"strmprivacy/strm/pkg/common"
	"strmprivacy/strm/pkg/util"
)

var printer util.Printer

func configurePrinter(command *cobra.Command) util.Printer {
	outputFormat := util.GetStringAndErr(command.Flags(), common.OutputFormatFlag)

	p := availablePrinters()[outputFormat+command.Parent().Name()]

	if p == nil {
		common.CliExit(fmt.Sprintf("Output format '%v' is not supported. Allowed values: %v", outputFormat, common.OutputFormatFlagAllowedValuesText))
	}

	return p
}

func availablePrinters() map[string]util.Printer {
	return util.MergePrinterMaps(
		util.DefaultPrinters,
		map[string]util.Printer{
			common.OutputFormatTable + common.ListCommandName:   listTablePrinter{},
			common.OutputFormatTable + common.GetCommandName:    getTablePrinter{},
			common.OutputFormatTable + common.DeleteCommandName: deletePrinter{},
			common.OutputFormatTable + common.CreateCommandName: createTablePrinter{},
			common.OutputFormatPlain + common.ListCommandName:   listPlainPrinter{},
			common.OutputFormatPlain + common.GetCommandName:    getPlainPrinter{},
			common.OutputFormatPlain + common.DeleteCommandName: deletePrinter{},
			common.OutputFormatPlain + common.CreateCommandName: createPlainPrinter{},
		},
	)
}

type listPlainPrinter struct{}
type getPlainPrinter struct{}
type createPlainPrinter struct{}

type listTablePrinter struct{}
type getTablePrinter struct{}
type createTablePrinter struct{}

type deletePrinter struct{}

func (p listTablePrinter) Print(data interface{}) {
	listResponse, _ := (data).(*batch_exporters.ListBatchExportersResponse)
	printTable(listResponse.BatchExporters)
}

func (p getTablePrinter) Print(data interface{}) {
	getResponse, _ := (data).(*batch_exporters.GetBatchExporterResponse)
	printTable([]*v1.BatchExporter{getResponse.BatchExporter})
}

func (p createTablePrinter) Print(data interface{}) {
	createResponse, _ := (data).(*batch_exporters.CreateBatchExporterResponse)
	printTable([]*v1.BatchExporter{createResponse.BatchExporter})
}

func (p listPlainPrinter) Print(data interface{}) {
	listResponse, _ := (data).(*batch_exporters.ListBatchExportersResponse)
	printPlain(listResponse.BatchExporters)
}

func (p getPlainPrinter) Print(data interface{}) {
	getResponse, _ := (data).(*batch_exporters.GetBatchExporterResponse)
	printPlain([]*v1.BatchExporter{getResponse.BatchExporter})
}

func (p createPlainPrinter) Print(data interface{}) {
	createResponse, _ := (data).(*batch_exporters.CreateBatchExporterResponse)
	printPlain([]*v1.BatchExporter{createResponse.BatchExporter})
}

func (p deletePrinter) Print(data interface{}) {
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

	util.RenderPlain(names)
}
