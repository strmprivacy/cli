package kafka_exporter

import (
	"fmt"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/spf13/cobra"
	"github.com/streammachineio/api-definitions-go/api/entities/v1"
	"github.com/streammachineio/api-definitions-go/api/kafka_exporters/v1"
	"streammachine.io/strm/pkg/common"
	"streammachine.io/strm/pkg/entity/kafka_cluster"
	"streammachine.io/strm/pkg/util"
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
	listResponse, _ := (data).(*kafka_exporters.ListKafkaExportersResponse)
	printTable(listResponse.KafkaExporters)
}

func (p getTablePrinter) Print(data interface{}) {
	getResponse, _ := (data).(*kafka_exporters.GetKafkaExporterResponse)
	printTable([]*entities.KafkaExporter{getResponse.KafkaExporter})
}

func (p createTablePrinter) Print(data interface{}) {
	createResponse, _ := (data).(*kafka_exporters.CreateKafkaExporterResponse)
	printTable([]*entities.KafkaExporter{createResponse.KafkaExporter})
}

func (p listPlainPrinter) Print(data interface{}) {
	listResponse, _ := (data).(*kafka_exporters.ListKafkaExportersResponse)
	printPlain(listResponse.KafkaExporters)
}

func (p getPlainPrinter) Print(data interface{}) {
	getResponse, _ := (data).(*kafka_exporters.GetKafkaExporterResponse)
	printPlain([]*entities.KafkaExporter{getResponse.KafkaExporter})
}

func (p createPlainPrinter) Print(data interface{}) {
	createResponse, _ := (data).(*kafka_exporters.CreateKafkaExporterResponse)
	printPlain([]*entities.KafkaExporter{createResponse.KafkaExporter})
}

func (p deletePrinter) Print(data interface{}) {
	fmt.Println("Kafka Exporter has been deleted")
}

func printTable(kafkaExporters []*entities.KafkaExporter) {
	rows := make([]table.Row, 0, len(kafkaExporters))

	for _, exporter := range kafkaExporters {
		rows = append(rows, table.Row{
			exporter.Ref.Name,
			exporter.StreamRef.Name,
			kafka_cluster.RefToString(exporter.Target.ClusterRef),
			exporter.Target.Topic,
		})
	}

	util.RenderTable(
		table.Row{
			"Kafka Exporter",
			"Stream",
			"Target Cluster",
			"Topic",
		},
		rows,
	)
}

func printPlain(kafkaExporters []*entities.KafkaExporter) {
	var names string
	lastIndex := len(kafkaExporters) - 1

	for index, exporter := range kafkaExporters {
		names = names + exporter.Ref.Name

		if index != lastIndex {
			names = names + "\n"
		}
	}

	util.RenderPlain(names)
}
