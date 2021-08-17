package kafka_cluster

import (
	"fmt"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/spf13/cobra"
	"github.com/streammachineio/api-definitions-go/api/entities/v1"
	"github.com/streammachineio/api-definitions-go/api/kafka_clusters/v1"
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
		}

		return util.GenericPrettyJsonPrinter{}
	case "plain":
		switch command.Parent().Name() {
		case constants.ListCommandName:
			return listPlainPrinter{}
		case constants.GetCommandName:
			return getPlainPrinter{}
		}

		return util.GenericPrettyJsonPrinter{}
	default:
		return util.GenericPrettyJsonPrinter{}
	}
}

type listPlainPrinter struct{}
type getPlainPrinter struct{}

type listTablePrinter struct{}
type getTablePrinter struct{}

func (p listTablePrinter) Print(data proto.Message) {
	listResponse, _ := (data).(*kafka_clusters.ListKafkaClustersResponse)
	printTable(listResponse.KafkaClusters)
}

func (p getTablePrinter) Print(data proto.Message) {
	getResponse, _ := (data).(*kafka_clusters.GetKafkaClusterResponse)
	printTable([]*entities.KafkaCluster{getResponse.KafkaCluster})
}

func (p listPlainPrinter) Print(data proto.Message) {
	listResponse, _ := (data).(*kafka_clusters.ListKafkaClustersResponse)
	printPlain(listResponse.KafkaClusters)
}

func (p getPlainPrinter) Print(data proto.Message) {
	getResponse, _ := (data).(*kafka_clusters.GetKafkaClusterResponse)
	printPlain([]*entities.KafkaCluster{getResponse.KafkaCluster})
}

func printTable(kafkaClusters []*entities.KafkaCluster) {
	rows := make([]table.Row, 0, len(kafkaClusters))

	for _, cluster := range kafkaClusters {

		rows = append(rows, table.Row{
			RefToString(cluster.Ref),
			cluster.BootstrapServers,
			cluster.AuthMechanism,
			cluster.TokenUri,
		})
	}

	util.RenderTable(
		table.Row{
			"Name",
			"Bootstrap Servers",
			"Auth Mechanism",
			"Token URI",
		},
		rows,
	)
}

func printPlain(kafkaClusters []*entities.KafkaCluster) {
	var names string
	lastIndex := len(kafkaClusters) - 1

	for index, cluster := range kafkaClusters {
		names = names + RefToString(cluster.Ref)

		if index != lastIndex {
			names = names + "\n"
		}
	}

	fmt.Println(names)
}
