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
	flags.String(EventsApiUrlFlag, "https://events.strmprivacy.io/event", "endpoint to send events to")
	flags.Int(IntervalFlag, 1000, "interval in ms. between simulated events")
	flags.Int(SessionRangeFlag, 1000, "number of different sessions being generated")
	flags.String(SessionPrefixFlag, "session", "prefix string for sessions")
	flags.Bool(QuietFlag, false, "do not print to stderr")
	flags.StringSlice(purposesFlag, []string{"", "0", "0/1", "0/1/2", "0/1/2/3"}, "purpose consent to be simulated")

	return simCmd
}
