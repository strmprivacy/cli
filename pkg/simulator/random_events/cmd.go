package random_events

import (
	"github.com/spf13/cobra"
	"strmprivacy/strm/pkg/common"
	"strmprivacy/strm/pkg/entity/stream"
	"strmprivacy/strm/pkg/simulator"
)

var longDoc = `The global ` + "`simulate`" + ` command runs a simulation with random events against a given Source Stream, using the ClickStream
schema.

### Usage`

func RunCmd() (cmd *cobra.Command) {
	simCmd := &cobra.Command{
		Use:               "random-events [stream-name]",
		Short:             "Run a simulator that will send random events to a stream",
		Long:              longDoc,
		DisableAutoGenTag: true,
		Run:               func(cmd *cobra.Command, args []string) { run(cmd, &args[0]) },
		Args:              cobra.ExactArgs(1), // the stream name
		ValidArgsFunction: stream.SourceNamesCompletion,
	}
	flags := simCmd.Flags()
	flags.String(sim.EventsApiUrlFlag, "https://events.strmprivacy.io/event", "Endpoint to send events to")
	flags.Int(sim.IntervalFlag, 1000, "Interval in ms. between simulated events")
	flags.Int(sim.SessionRangeFlag, 1000, "Number of different sessions being generated")
	flags.String(sim.SessionPrefixFlag, "session", "Prefix string for sessions")
	flags.String(common.ClientIdFlag, "", "Client id to be used for sending data")
	flags.String(common.ClientSecretFlag, "", "Client secret to be used for sending data")
	flags.Bool(sim.QuietFlag, false, "don't spam stderr")
	flags.StringSlice(sim.ConsentLevelsFlag, []string{"", "0", "0/1", "0/1/2", "0/1/2/3"}, "consent levels to be simulated")
	return simCmd
}
