package random_events

import (
	"github.com/spf13/cobra"
	"strmprivacy/strm/pkg/entity/stream"
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
		Run:               func(cmd *cobra.Command, args []string) { randomEvents(cmd, &args[0]) },
		Args:              cobra.ExactArgs(1), // the stream name
		ValidArgsFunction: stream.SourceNamesCompletion,
	}

	flags := simCmd.Flags()
	flags.String(EventsApiUrlFlag, "https://events.strmprivacy.io/event", "Endpoint to send events to")
	flags.Int(IntervalFlag, 1000, "Interval in ms. between simulated events")
	flags.Int(SessionRangeFlag, 1000, "Number of different sessions being generated")
	flags.String(SessionPrefixFlag, "session", "Prefix string for sessions")
	flags.Bool(QuietFlag, false, "Do not print to stderr")
	flags.StringSlice(ConsentLevelsFlag, []string{"", "0", "0/1", "0/1/2", "0/1/2/3"}, "consent levels to be simulated")

	return simCmd
}
