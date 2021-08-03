package event_contract

import "github.com/spf13/cobra"

func GetCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "event-contract (name)",
		Short: "Get Event Contract by name",
		Run: func(cmd *cobra.Command, args []string) {
			get(&args[0])
		},
		Args:              cobra.ExactArgs(1), // the stream name
		ValidArgsFunction: refsCompletion,
	}
}
func ListCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "event-contracts",
		Short: "List Event Contracts",
		Run: func(cmd *cobra.Command, args []string) {
			list()
		},
	}
}
func CreateCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "event-contract (file)",
		Short: "create an event-contract",
		Long:  `create an event contract from a json definition file`,
		Run: func(cmd *cobra.Command, args []string) {
			create(cmd, &args[0])
		},
		Args: cobra.ExactArgs(1), // the file name
	}

}