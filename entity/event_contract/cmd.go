package event_contract

import "github.com/spf13/cobra"

func GetCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "event-contract [name]",
		Short: "get event-contract by name",
		Run: func(cmd *cobra.Command, args []string) {
			get(&args[0])
		},
		Args:              cobra.ExactArgs(1), // the stream name
		ValidArgsFunction: ExistingNamesCompletion,
	}
}
func ListCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "event-contracts",
		Short: "List event-contracts",
		Run: func(cmd *cobra.Command, args []string) {
			list()
		},
	}
}
