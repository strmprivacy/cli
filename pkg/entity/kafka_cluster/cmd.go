package kafka_cluster

import (
	"github.com/spf13/cobra"
)

var longDoc = `A Kafka Cluster can be used for exporting directly from STRM Privacy to a Kafka Cluster owned by the client, or to the
shared Kafka Export Cluster, hosted by STRM Privacy. This gives all the performance, scalability and reliability
benefits offered by Kafka.

The Kafka Cluster is only a configuration object, it does not create the actual cluster infrastructure. It only points
to an existing Kafka Cluster.

At the moment, it’s not possible to create your own Kafka Cluster. All Kafka Exporters use the STRM Privacy Shared
Cluster.

### Usage
`

func GetCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "kafka-cluster (name)",
		Short: "Get Kafka Cluster by name",
		Long:  longDoc,
		PreRun: func(cmd *cobra.Command, args []string) {
			printer = configurePrinter(cmd)
		},
		DisableAutoGenTag: true,
		Run: func(cmd *cobra.Command, args []string) {
			get(&args[0])
		},
		Args:              cobra.ExactArgs(1), // the stream name
		ValidArgsFunction: NamesCompletion,
	}
}
func ListCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "kafka-clusters",
		Short: "List Kafka Clusters",
		Long:  longDoc,
		PreRun: func(cmd *cobra.Command, args []string) {
			printer = configurePrinter(cmd)
		},
		DisableAutoGenTag: true,
		Run: func(cmd *cobra.Command, args []string) {
			list()
		},
	}
}
