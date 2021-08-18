package sink

import (
	"fmt"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/spf13/cobra"
	"github.com/streammachineio/api-definitions-go/api/entities/v1"
	"github.com/streammachineio/api-definitions-go/api/sinks/v1"
	"streammachine.io/strm/pkg/common"
	"streammachine.io/strm/pkg/constants"
	"streammachine.io/strm/pkg/util"
)

var printer util.Printer

func configurePrinter(command *cobra.Command) util.Printer {
	outputFormat := util.GetStringAndErr(command.Flags(), constants.OutputFormatFlag)

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
			constants.OutputFormatTable + constants.ListCommandName:   listTablePrinter{},
			constants.OutputFormatTable + constants.GetCommandName:    getTablePrinter{},
			constants.OutputFormatTable + constants.DeleteCommandName: deletePrinter{},
			constants.OutputFormatTable + constants.CreateCommandName: createTablePrinter{},
			constants.OutputFormatPlain + constants.ListCommandName:   listPlainPrinter{},
			constants.OutputFormatPlain + constants.GetCommandName:    getPlainPrinter{},
			constants.OutputFormatPlain + constants.DeleteCommandName: deletePrinter{},
			constants.OutputFormatPlain + constants.CreateCommandName: createPlainPrinter{},
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
