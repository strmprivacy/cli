package usage

import (
	"github.com/spf13/cobra"
	"streammachine.io/strm/pkg/entity/stream"
)

const (
	jsonFlag    = "json"
	fromFlag    = "from"  // iso or .. ago
	untilFlag   = "until" // iso or .. ago
	aggregateBy = "by"    // aggregation seconds
)

func GetCmd() *cobra.Command {
	usage := &cobra.Command{
		Use:   "usage [name]",
		Short: "Get Usage for a certain stream name",
		Run: func(cmd *cobra.Command, args []string) {
			get(cmd, &args[0])
		},
		Args:              cobra.ExactArgs(1), // the stream name
		ValidArgsFunction: stream.SourceNamesCompletion,
	}
	flags := usage.Flags()

	flags.Bool(jsonFlag, false, "json output")
	flags.String(fromFlag, "", "from")
	flags.String(untilFlag, "", "from")
	flags.String(aggregateBy, "", "aggregate (s)")

	return usage
}