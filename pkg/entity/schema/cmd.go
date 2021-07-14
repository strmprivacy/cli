package schema

import "github.com/spf13/cobra"

func GetCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "schema [name]",
		Short: "Get schema by name",
		Run: func(cmd *cobra.Command, args []string) {
			get(&args[0])
		},
		Args:              cobra.ExactArgs(1), // the stream name
		ValidArgsFunction: namesCompletion,
	}
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
