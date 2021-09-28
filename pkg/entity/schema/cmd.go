package schema

import "github.com/spf13/cobra"

func GetCmd() *cobra.Command {
	getSchema := &cobra.Command{
		Use:   "schema [name]",
		Short: "Get schema by name",
		PreRun: func(cmd *cobra.Command, args []string) {
			printer = configurePrinter(cmd)
		},
		Run: func(cmd *cobra.Command, args []string) {
			get(&args[0], cmd)
		},
		Args:              cobra.ExactArgs(1), // the stream name
		ValidArgsFunction: NamesCompletion,
	}
	flags := getSchema.Flags()
	flags.String(kafkaClusterFlag, "", "if present, find the corresponding Confluent ID for the given Kafka cluster")

	return getSchema
}
func ListCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "schemas",
		Short: "List schemas",
		PreRun: func(cmd *cobra.Command, args []string) {
			printer = configurePrinter(cmd)
		},
		Run: func(cmd *cobra.Command, args []string) {
			list()
		},
	}
}

func CreateCmd() *cobra.Command {
	createCmd := &cobra.Command{
		Use:   "schema (handle/name/version)",
		Short: "create a schema",
		PreRun: func(cmd *cobra.Command, args []string) {
			printer = configurePrinter(cmd)
		},
		Run: func(cmd *cobra.Command, args []string) {
			create(cmd, &args[0])
		},
		Args: cobra.ExactArgs(1),
	}
	flags := createCmd.Flags()
	flags.String(definitionFlag, "", "filename of the definition")
	_ = createCmd.MarkFlagRequired(definitionFlag)
	flags.Bool(publicFlag, false, "should the schema become public")
	return createCmd
}