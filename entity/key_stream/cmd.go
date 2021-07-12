package key_stream

import "github.com/spf13/cobra"

func GetCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "key-stream [name]",
		Short: "Get key stream by name",
		Run: func(cmd *cobra.Command, args []string) {
			get(&args[0])
		},
		Args:              cobra.ExactArgs(1), // the stream name
		ValidArgsFunction: ExistingNamesCompletion,
	}
}
func ListCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "key-streams",
		Short: "List key streams",
		Run: func(cmd *cobra.Command, args []string) {
			list()
		},
	}
}
