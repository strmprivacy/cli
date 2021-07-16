package schema

import "github.com/spf13/cobra"

func GetCmd() *cobra.Command {
	getSchema := &cobra.Command{
		Use:   "schema [name]",
		Short: "Get schema by name",
		Run: func(cmd *cobra.Command, args []string) {
			get(&args[0], cmd)
		},
		Args:              cobra.ExactArgs(1), // the stream name
		ValidArgsFunction: namesCompletion,
	}
	flags := getSchema.Flags()
	flags.String(kafkaClusterFlag, "", "if present, find the corresponding Confluent ID for the given Kafka cluster")

	return getSchema
}
func ListCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "schemas",
		Short: "List schemas",
		Run: func(cmd *cobra.Command, args []string) {
			list()
		},
	}
}
