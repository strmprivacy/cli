package usage

import (
	"fmt"
	"github.com/spf13/cobra"
	"strmprivacy/strm/pkg/common"
	"strmprivacy/strm/pkg/entity/stream"
)

const (
	jsonFlag        = "json"
	fromFlag        = "from"  // iso or .. ago
	untilFlag       = "until" // iso or .. ago
	aggregateByFlag = "by"    // aggregation seconds
)

func GetCmd() *cobra.Command {
	usage := &cobra.Command{
		Use:               "usage (stream-name)",
		Short:             "Get Usage for a certain stream name",
		Long:              longGetDoc,
		Example:           getExample,
		DisableAutoGenTag: true,
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

	flags.StringP(common.OutputFormatFlag, common.OutputFormatFlagShort, "csv", fmt.Sprintf("Usage output format [%v]", common.UsageOutputFormatFlagAllowedValuesText))
	err := usage.RegisterFlagCompletionFunc(common.OutputFormatFlag, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return common.UsageOutputFormatFlagAllowedValues, cobra.ShellCompDirectiveNoFileComp
	})

	common.CliExit(err)

	return usage
}
