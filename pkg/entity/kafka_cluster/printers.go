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
			return ListStreamsTablePrinter{}
		case constants.GetCommandName:
			return GetStreamTablePrinter{}
		}

		return util.GenericPrettyJsonPrinter{}
	case "plain":
		switch command.Parent().Name() {
		case constants.ListCommandName:
			return ListStreamsPlainPrinter{}
		case constants.GetCommandName:
			return GetStreamPlainPrinter{}
		}

		return util.GenericPrettyJsonPrinter{}
	default:
		return util.GenericPrettyJsonPrinter{}
	}
}

type ListStreamsPlainPrinter struct{}
type GetStreamPlainPrinter struct{}

type ListStreamsTablePrinter struct{}
type GetStreamTablePrinter struct{}

func (p ListStreamsTablePrinter) Print(data proto.Message) {
	listResponse, _ := (data).(*kafka_clusters.ListKafkaClustersResponse)
	printTable(listResponse.KafkaClusters)
}

func (p GetStreamTablePrinter) Print(data proto.Message) {
	getResponse, _ := (data).(*kafka_clusters.GetKafkaClusterResponse)
	printTable([]*entities.KafkaCluster{getResponse.KafkaCluster})
}

func (p ListStreamsPlainPrinter) Print(data proto.Message) {
	listResponse, _ := (data).(*kafka_clusters.ListKafkaClustersResponse)
	printPlain(listResponse.KafkaClusters)
}

func (p GetStreamPlainPrinter) Print(data proto.Message) {
	getResponse, _ := (data).(*kafka_clusters.GetKafkaClusterResponse)
	printPlain([]*entities.KafkaCluster{getResponse.KafkaCluster})
}

func printTable(kafkaClusters []*entities.KafkaCluster) {
	rows := make([]table.Row, 0, len(kafkaClusters))

	for _, cluster := range kafkaClusters {

		rows = append(rows, table.Row{
			refToString(cluster.Ref),
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
		names = names + refToString(cluster.Ref)

		if index != lastIndex {
			names = names + "\n"
		}
	}

	fmt.Println(names)
}

func refToString(clusterRef *entities.KafkaClusterRef) string {
	return fmt.Sprintf("%v/%v", clusterRef.BillingId, clusterRef.Name)
}
