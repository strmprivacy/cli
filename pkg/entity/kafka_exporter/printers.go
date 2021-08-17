package kafka_exporter

import (
	"fmt"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/spf13/cobra"
	"github.com/streammachineio/api-definitions-go/api/entities/v1"
	"github.com/streammachineio/api-definitions-go/api/kafka_exporters/v1"
	"google.golang.org/protobuf/proto"
	"streammachine.io/strm/pkg/constants"
	"streammachine.io/strm/pkg/entity/kafka_cluster"
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
			return ListKafkaExportersTablePrinter{}
		case constants.GetCommandName:
			return GetKafkaExportersTablePrinter{}
		case constants.DeleteCommandName:
			return DeleteKafkaExportersPrinter{}
		case constants.CreateCommandName:
			return CreateKafkaExportersTablePrinter{}
		}

		return util.GenericPrettyJsonPrinter{}
	case "plain":
		switch command.Parent().Name() {
		case constants.ListCommandName:
			return ListKafkaExportersPlainPrinter{}
		case constants.GetCommandName:
			return GetKafkaExportersPlainPrinter{}
		case constants.DeleteCommandName:
			return DeleteKafkaExportersPrinter{}
		case constants.CreateCommandName:
			return CreateKafkaExportersPlainPrinter{}
		}

		return util.GenericPrettyJsonPrinter{}
	default:
		return util.GenericPrettyJsonPrinter{}
	}
}

type ListKafkaExportersPlainPrinter struct{}
type GetKafkaExportersPlainPrinter struct{}
type CreateKafkaExportersPlainPrinter struct{}

type ListKafkaExportersTablePrinter struct{}
type GetKafkaExportersTablePrinter struct{}
type CreateKafkaExportersTablePrinter struct{}

type DeleteKafkaExportersPrinter struct{}

func (p ListKafkaExportersTablePrinter) Print(data proto.Message) {
	listResponse, _ := (data).(*kafka_exporters.ListKafkaExportersResponse)
	printTable(listResponse.KafkaExporters)
}

func (p GetKafkaExportersTablePrinter) Print(data proto.Message) {
	getResponse, _ := (data).(*kafka_exporters.GetKafkaExporterResponse)
	printTable([]*entities.KafkaExporter{getResponse.KafkaExporter})
}

func (p CreateKafkaExportersTablePrinter) Print(data proto.Message) {
	createResponse, _ := (data).(*kafka_exporters.CreateKafkaExporterResponse)
	printTable([]*entities.KafkaExporter{createResponse.KafkaExporter})
}

func (p ListKafkaExportersPlainPrinter) Print(data proto.Message) {
	listResponse, _ := (data).(*kafka_exporters.ListKafkaExportersResponse)
	printPlain(listResponse.KafkaExporters)
}

func (p GetKafkaExportersPlainPrinter) Print(data proto.Message) {
	getResponse, _ := (data).(*kafka_exporters.GetKafkaExporterResponse)
	printPlain([]*entities.KafkaExporter{getResponse.KafkaExporter})
}

func (p CreateKafkaExportersPlainPrinter) Print(data proto.Message) {
	createResponse, _ := (data).(*kafka_exporters.CreateKafkaExporterResponse)
	printPlain([]*entities.KafkaExporter{createResponse.KafkaExporter})
}

func (p DeleteKafkaExportersPrinter) Print(_ proto.Message) {
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

	fmt.Println(names)
}
