package kafka_exporter

import (
	"github.com/spf13/cobra"
	"io/ioutil"
	"strmprivacy/strm/pkg/common"
	"strmprivacy/strm/pkg/entity/kafka_cluster"
	"strmprivacy/strm/pkg/entity/stream"
)

const (
	clusterFlag = "cluster"
	saveFlag    = "save"
)

var content, _ = ioutil.ReadFile("pkg/entity/kafka_exporter/docstring.md")

func DeleteCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "kafka-exporter [name]",
		Short: "Delete a Kafka exporter",
		//	Long: `Delete a Kafka exporter.
		//If a kafka-exporter has dependents (like Kafka users), you can use
		//the 'recursive' option to get rid of those also.
		//Returns everything that was deleted. `,
		Long:              string(content),
		DisableAutoGenTag: true,
		PreRun: func(cmd *cobra.Command, args []string) {
			printer = configurePrinter(cmd)
		},
		Run: func(cmd *cobra.Command, args []string) {
			recursive, _ := cmd.Flags().GetBool(common.RecursiveFlagName)
			del(&args[0], recursive)
		},
		Args:              cobra.ExactArgs(1), // the stream name
		ValidArgsFunction: NamesCompletion,
	}

}
func GetCmd() *cobra.Command {
	return &cobra.Command{
		Use:               "kafka-exporter [name]",
		Short:             "Get Kafka exporter by name",
		Long:              string(content),
		DisableAutoGenTag: true,
		PreRun: func(cmd *cobra.Command, args []string) {
			printer = configurePrinter(cmd)
		},
		Run: func(cmd *cobra.Command, args []string) {
			recursive, _ := cmd.Flags().GetBool(common.RecursiveFlagName)
			get(&args[0], recursive)
		},
		Args:              cobra.ExactArgs(1), // the stream name
		ValidArgsFunction: NamesCompletion,
	}
}

func ListCmd() *cobra.Command {
	return &cobra.Command{
		Use:               "kafka-exporters",
		Short:             "List Kafka exporters",
		Long:              string(content),
		DisableAutoGenTag: true,
		PreRun: func(cmd *cobra.Command, args []string) {
			printer = configurePrinter(cmd)
		},
		Run: func(cmd *cobra.Command, args []string) {
			flag, _ := cmd.Root().PersistentFlags().GetBool(common.RecursiveFlagName)
			list(flag)
		},
	}
}

func CreateCmd() *cobra.Command {
	kafkaExporter := &cobra.Command{
		Use:               "kafka-exporter [stream-name]",
		Short:             "Create a Kafka exporter",
		Long:              string(content),
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
