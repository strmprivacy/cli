package randomsim

import (
	"github.com/spf13/cobra"
	"streammachine.io/strm/sims"
	"streammachine.io/strm/entity/stream"
)

func RunCmd() (cmd *cobra.Command) {
	simCmd := &cobra.Command{
		Use:   "run-random [stream-name]",
		Short: "run random sim",
		Long: `Run a random sim using the clickstream schema

Uses the saved stream definition if available. `,
		Run:               func(cmd *cobra.Command, args []string) { run(cmd, &args[0]) },
		Args:              cobra.ExactArgs(1), // the stream name
		ValidArgsFunction: stream.ExistingSourceStreamNames,
	}
	flags := simCmd.Flags()
	flags.String(sims.EventGatewayFlag, "https://in.strm.services/event", "where to send the events")
	flags.Int(sims.IntervalFlag, 1000, "interval in ms. between simulated events")
	flags.Int(sims.SessionRangeFlag, 1000, "number of different sessions being generated")
	flags.String(sims.SessionPrefixFlag, "session", "prefix string for sessions")
	flags.String(sims.ClientIdFlag, "", "client id to be used for sending data")
	flags.String(sims.ClientSecretFlag, "", "client secret to be used for sending data")
	flags.Bool(sims.QuietFlag, false, "don't spam stderr")
	flags.StringSlice(sims.ConsentLevelsFlag, []string{"0", "0/1", "0/1/2", "0/1/2/3"}, "consent levels to be simulated")
	return simCmd
}
