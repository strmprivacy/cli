package usage

import (
	"fmt"
	"github.com/spf13/cobra"
	"streammachine.io/strm/pkg/common"
	"streammachine.io/strm/pkg/constants"
	"streammachine.io/strm/pkg/entity/stream"
	"streammachine.io/strm/pkg/util"
)

const (
	jsonFlag        = "json"
	fromFlag        = "from"  // iso or .. ago
	untilFlag       = "until" // iso or .. ago
	aggregateByFlag = "by"    // aggregation seconds
)

func GetCmd() *cobra.Command {
	usage := &cobra.Command{
		Use:   "usage (stream-name)",
		Short: "Get Usage for a certain stream name",
		Long: `Returns the usage for a certain stream over a certain period.

The values are interpolated from cumulative event accounts, and sampled over intervals
(the --by option). The default output is csv, but json is also available.

The default range is over the last 24 hours, with a default interval of 15 minutes.

Example:

strm get usage bart --by 15m --from 2021/7/27-10:00  --until 2021/7/27-12:00

from,count,duration,change,rate
2021-07-27T10:00:00.000000+0200,173478,900,NaN,NaN
2021-07-27T10:15:00.000000+0200,182422,900,8944,9.94
2021-07-27T10:30:00.000000+0200,191363,900,8941,9.93
2021-07-27T10:45:00.000000+0200,200305,900,8942,9.94
2021-07-27T11:00:00.000000+0200,209248,900,8943,9.94
2021-07-27T11:15:00.000000+0200,218192,900,8944,9.94
2021-07-27T11:30:00.000000+0200,227134,900,8942,9.94
2021-07-27T11:45:00.000000+0200,236078,900,8944,9.94
2021-07-27T12:00:00.000000+0200,245023,900,8945,9.94
`,
		PreRun: func(cmd *cobra.Command, args []string) {
			printer = configurePrinter(cmd)
		},
		Run: func(cmd *cobra.Command, args []string) {
			get(cmd, &args[0])
		},
		Args:              cobra.ExactArgs(1), // the stream name
		ValidArgsFunction: stream.SourceNamesCompletion,
	}
	flags := usage.Flags()

	flags.String(fromFlag, "", fmt.Sprintf("from %s", dateTimeParseFormat))
	flags.String(untilFlag, "", fmt.Sprintf("until %s", dateTimeParseFormat))
	flags.String(aggregateByFlag, "", "aggregate by (seconds|..m|..h|..d)")
	_ = usage.RegisterFlagCompletionFunc(aggregateByFlag, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"5m", "15m", "1h", "6h", "24h"}, cobra.ShellCompDirectiveDefault
	})
	_ = usage.RegisterFlagCompletionFunc(fromFlag, dateCompletion)
	_ = usage.RegisterFlagCompletionFunc(untilFlag, dateCompletion)

	flags.StringP(util.OutputFormatFlag, "o", "csv", fmt.Sprintf("Usage output format [%v]", constants.UsageOutputFormatFlagAllowedValuesText))
	err := usage.RegisterFlagCompletionFunc(util.OutputFormatFlag, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return constants.UsageOutputFormatFlagAllowedValues, cobra.ShellCompDirectiveNoFileComp
	})

	common.CliExit(err)

	return usage
}
