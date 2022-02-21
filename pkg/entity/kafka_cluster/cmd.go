package kafka_cluster

import (
	"github.com/spf13/cobra"
	"io/ioutil"
)

var content, _ = ioutil.ReadFile("pkg/entity/kafka_cluster/docstring.md")

func GetCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "kafka-cluster [name]",
		Short: "Get Kafka cluster by name",
		Long:  string(content),
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
		Short: "List Kafka clusters",
		Long:  string(content),
		PreRun: func(cmd *cobra.Command, args []string) {
			printer = configurePrinter(cmd)
		},
		DisableAutoGenTag: true,
		Run: func(cmd *cobra.Command, args []string) {
			list()
		},
	}
}
