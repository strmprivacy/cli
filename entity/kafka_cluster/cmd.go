package kafka_cluster

import "github.com/spf13/cobra"

func GetCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "kafka-cluster [name]",
		Short: "get kafka-cluster by name",
		Run: func(cmd *cobra.Command, args []string) {
			get(&args[0])
		},
		Args:              cobra.ExactArgs(1), // the stream name
		ValidArgsFunction: ExistingNamesCompletion,
	}
}
func ListCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "kafka-clusters",
		Short: "List kafka-clusters",
		Run: func(cmd *cobra.Command, args []string) {
			list()
		},
	}
}
