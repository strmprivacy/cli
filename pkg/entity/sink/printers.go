package sink

import (
	"fmt"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/spf13/cobra"
	"github.com/strmprivacy/api-definitions-go/v2/api/entities/v1"
	"github.com/strmprivacy/api-definitions-go/v2/api/sinks/v1"
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

func (p listTablePrinter) Print(data interface{}) {
	listResponse, _ := (data).(*sinks.ListSinksResponse)
	printTable(listResponse.Sinks, p.recursive)
}

func (p getTablePrinter) Print(data interface{}) {
	getResponse, _ := (data).(*sinks.GetSinkResponse)
	printTable([]*entities.SinkTree{getResponse.SinkTree}, p.recursive)
}

func (p createTablePrinter) Print(data interface{}) {
	createResponse, _ := (data).(*sinks.CreateSinkResponse)
	printTable([]*entities.SinkTree{{Sink: createResponse.Sink}}, false)
}

func (p listPlainPrinter) Print(data interface{}) {
	listResponse, _ := (data).(*sinks.ListSinksResponse)
	printPlain(listResponse.Sinks)
}

func (p getPlainPrinter) Print(data interface{}) {
	getResponse, _ := (data).(*sinks.GetSinkResponse)
	printPlain([]*entities.SinkTree{getResponse.SinkTree})
}

func (p createPlainPrinter) Print(data interface{}) {
	createResponse, _ := (data).(*sinks.CreateSinkResponse)
	printPlain([]*entities.SinkTree{{Sink: createResponse.Sink}})
}

func (p deletePrinter) Print(data interface{}) {
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
