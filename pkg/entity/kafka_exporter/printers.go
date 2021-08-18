package kafka_exporter

import (
	"fmt"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/spf13/cobra"
	"github.com/streammachineio/api-definitions-go/api/entities/v1"
	"github.com/streammachineio/api-definitions-go/api/kafka_exporters/v1"
	"google.golang.org/protobuf/proto"
	"streammachine.io/strm/pkg/common"
	"streammachine.io/strm/pkg/constants"
	"streammachine.io/strm/pkg/entity/kafka_cluster"
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

type listTablePrinter struct{}
type getTablePrinter struct{}
type createTablePrinter struct{}

type deletePrinter struct{}

func (p listTablePrinter) Print(data proto.Message) {
	listResponse, _ := (data).(*kafka_exporters.ListKafkaExportersResponse)
	printTable(listResponse.KafkaExporters)
}

func (p getTablePrinter) Print(data proto.Message) {
	getResponse, _ := (data).(*kafka_exporters.GetKafkaExporterResponse)
	printTable([]*entities.KafkaExporter{getResponse.KafkaExporter})
}

func (p createTablePrinter) Print(data proto.Message) {
	createResponse, _ := (data).(*kafka_exporters.CreateKafkaExporterResponse)
	printTable([]*entities.KafkaExporter{createResponse.KafkaExporter})
}

func (p listPlainPrinter) Print(data proto.Message) {
	listResponse, _ := (data).(*kafka_exporters.ListKafkaExportersResponse)
	printPlain(listResponse.KafkaExporters)
}

func (p getPlainPrinter) Print(data proto.Message) {
	getResponse, _ := (data).(*kafka_exporters.GetKafkaExporterResponse)
	printPlain([]*entities.KafkaExporter{getResponse.KafkaExporter})
}

func (p createPlainPrinter) Print(data proto.Message) {
	createResponse, _ := (data).(*kafka_exporters.CreateKafkaExporterResponse)
	printPlain([]*entities.KafkaExporter{createResponse.KafkaExporter})
}

func (p deletePrinter) Print(_ proto.Message) {
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
