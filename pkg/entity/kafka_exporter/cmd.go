package kafka_exporter

import (
	"github.com/spf13/cobra"
	"strmprivacy/strm/pkg/common"
	"strmprivacy/strm/pkg/entity/kafka_cluster"
	"strmprivacy/strm/pkg/entity/stream"
)

const (
	clusterFlag = "cluster"
	saveFlag    = "save"
)

var longDoc = `A Kafka Exporter, like a Batch Exporter, can be used to export events from Stream Machine to somewhere outside of STRM
Privacy. But in contrast to a Batch Exporter, a Kafka Exporter does not work in batches, but processes the events in
real time.

After creation, the CLI exposes the authentication information that is needed to connect to it with your own Kafka
Consumer.

In case your data are Avro encoded, the Kafka exporter provides a *json format* conversion of your data for easier
downstream processing. See the [exporting Kafka](docs/03-quickstart/01-streaming/04-receiving-data/03-exporting-kafka.md) page for how to consume from the
exporter.

If a kafka-exporter has dependents (like Kafka users), you can use
the 'recursive' option to get rid of those also.
Returns everything that was deleted.

### Usage`

func DeleteCmd() *cobra.Command {
	return &cobra.Command{
		Use:               "kafka-exporter [name]",
		Short:             "Delete a Kafka exporter",
		Long:              longDoc,
		DisableAutoGenTag: true,
		PreRun: func(cmd *cobra.Command, args []string) {
			printer = configurePrinter(cmd)
		},
		Run: func(cmd *cobra.Command, args []string) {
			recursive, _ := cmd.Flags().GetBool(common.RecursiveFlagName)
			del(&args[0], recursive, cmd)
		},
		Args:              cobra.ExactArgs(1), // the stream name
		ValidArgsFunction: NamesCompletion,
	}

}
func GetCmd() *cobra.Command {
	return &cobra.Command{
		Use:               "kafka-exporter [name]",
		Short:             "Get Kafka exporter by name",
		Long:              longDoc,
		DisableAutoGenTag: true,
		PreRun: func(cmd *cobra.Command, args []string) {
			printer = configurePrinter(cmd)
		},
		Run: func(cmd *cobra.Command, args []string) {
			recursive, _ := cmd.Flags().GetBool(common.RecursiveFlagName)
			get(&args[0], cmd, recursive)
		},
		Args:              cobra.ExactArgs(1), // the stream name
		ValidArgsFunction: NamesCompletion,
	}
}

func ListCmd() *cobra.Command {
	return &cobra.Command{
		Use:               "kafka-exporters",
		Short:             "List Kafka exporters",
		Long:              longDoc,
		DisableAutoGenTag: true,
		PreRun: func(cmd *cobra.Command, args []string) {
			printer = configurePrinter(cmd)
		},
		Run: func(cmd *cobra.Command, args []string) {
			flag, _ := cmd.Root().PersistentFlags().GetBool(common.RecursiveFlagName)
			list(flag, cmd)
		},
	}
}

func CreateCmd() *cobra.Command {
	kafkaExporter := &cobra.Command{
		Use:               "kafka-exporter [stream-name]",
		Short:             "Create a Kafka exporter",
		Long:              longDoc,
		DisableAutoGenTag: true,
		PreRun: func(cmd *cobra.Command, args []string) {
			printer = configurePrinter(cmd)
		},
		Run: func(cmd *cobra.Command, args []string) {
			streamName := &args[0]
			create(streamName, cmd)
		},
		Args:              cobra.ExactArgs(1), // the stream name
		ValidArgsFunction: stream.NamesCompletion,
	}

	flags := kafkaExporter.Flags()
	flags.String(clusterFlag, "", "name of the kafka cluster")
	flags.Bool(saveFlag, false, "save the result in the config directory")

	// not yet handling the external cluster flags

	err := kafkaExporter.RegisterFlagCompletionFunc(clusterFlag, kafka_cluster.NamesCompletion)
	common.CliExit(err)
	return kafkaExporter
}
