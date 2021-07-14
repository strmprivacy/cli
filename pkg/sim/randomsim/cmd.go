package randomsim

import (
	"github.com/spf13/cobra"
	"streammachine.io/strm/pkg/entity/stream"
	"streammachine.io/strm/pkg/sim"
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
	flags.String(sim.EventGatewayFlag, "https://in.strm.services/event", "Endpoint to send events to")
	flags.Int(sim.IntervalFlag, 1000, "Interval in ms. between simulated events")
	flags.Int(sim.SessionRangeFlag, 1000, "Number of different sessions being generated")
	flags.String(sim.SessionPrefixFlag, "session", "Prefix string for sessions")
	flags.String(sim.ClientIdFlag, "", "Client id to be used for sending data")
	flags.String(sim.ClientSecretFlag, "", "Client secret to be used for sending data")
	flags.Bool(sim.QuietFlag, false, "don't spam stderr")
	flags.StringSlice(sim.ConsentLevelsFlag, []string{"0", "0/1", "0/1/2", "0/1/2/3"}, "consent levels to be simulated")
	return simCmd
}
