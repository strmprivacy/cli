package randomsim

import (
	"github.com/spf13/cobra"
	"streammachine.io/strm/entity/stream"
	"streammachine.io/strm/sims"
)

func RunCmd() (cmd *cobra.Command) {
	simCmd := &cobra.Command{
		Use:   "run-random [stream-name]",
		Short: "Run a simulator that will send random events to a stream",
		Long: `Run a random simulator using the clickstream schema

Uses a saved stream definition if available, otherwise, client id and secret are required`,
		Run:               func(cmd *cobra.Command, args []string) { run(cmd, &args[0]) },
		Args:              cobra.ExactArgs(1), // the stream name
		ValidArgsFunction: stream.SourceNamesCompletion,
	}
	flags := simCmd.Flags()
	flags.String(sims.EventGatewayFlag, "https://in.strm.services/event", "Endpoint to send events to")
	flags.Int(sims.IntervalFlag, 1000, "Interval in ms. between simulated events")
	flags.Int(sims.SessionRangeFlag, 1000, "Number of different sessions being generated")
	flags.String(sims.SessionPrefixFlag, "session", "Prefix string for sessions")
	flags.String(sims.ClientIdFlag, "", "Client id to be used for sending data")
	flags.String(sims.ClientSecretFlag, "", "Client secret to be used for sending data")
	flags.Bool(sims.QuietFlag, false, "don't spam stderr")
	flags.StringSlice(sims.ConsentLevelsFlag, []string{"0", "0/1", "0/1/2", "0/1/2/3"}, "consent levels to be simulated")
	return simCmd
}
