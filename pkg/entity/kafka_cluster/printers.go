package kafka_cluster

import (
	"fmt"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/spf13/cobra"
	"github.com/streammachineio/api-definitions-go/api/entities/v1"
	"github.com/streammachineio/api-definitions-go/api/kafka_clusters/v1"
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
			constants.OutputFormatTable + constants.ListCommandName: listTablePrinter{},
			constants.OutputFormatTable + constants.GetCommandName:  getTablePrinter{},
			constants.OutputFormatPlain + constants.ListCommandName: listPlainPrinter{},
			constants.OutputFormatPlain + constants.GetCommandName:  getPlainPrinter{},
		},
	)
}

type listPlainPrinter struct{}
type getPlainPrinter struct{}

type listTablePrinter struct{}
type getTablePrinter struct{}

func (p listTablePrinter) Print(data interface{}) {
	listResponse, _ := (data).(*kafka_clusters.ListKafkaClustersResponse)
	printTable(listResponse.KafkaClusters)
}

func (p getTablePrinter) Print(data interface{}) {
	getResponse, _ := (data).(*kafka_clusters.GetKafkaClusterResponse)
	printTable([]*entities.KafkaCluster{getResponse.KafkaCluster})
}

func (p listPlainPrinter) Print(data interface{}) {
	listResponse, _ := (data).(*kafka_clusters.ListKafkaClustersResponse)
	printPlain(listResponse.KafkaClusters)
}

func (p getPlainPrinter) Print(data interface{}) {
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

	util.RenderPlain(names)
}
