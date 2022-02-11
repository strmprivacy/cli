package stream

import (
	"github.com/spf13/cobra"
	"strmprivacy/strm/pkg/common"
	"strmprivacy/strm/pkg/entity/event_contract"
	"strmprivacy/strm/pkg/util"
)

func CreateCmd() *cobra.Command {
	stream := &cobra.Command{
		Use:   "stream [name]",
		Short: "Create a stream",
		PreRun: func(cmd *cobra.Command, args []string) {
			printer = configurePrinter(cmd)
		},
		Run: func(cmd *cobra.Command, args []string) {
			create(args, cmd)
		},
		Args: cobra.MaximumNArgs(1), // the stream name
	}
	flags := stream.Flags()
	flags.StringP(linkedStreamFlag, "D", "",
		"name of stream that this stream is derived from")

	// TODO github.com/thediveo/enumflag might be nicer!
	flags.String(consentLevelTypeFlag, "CUMULATIVE",
		"CUMULATIVE or GRANULAR")
	flags.Int32SliceP(consentLevelsFlag, "L", []int32{},
		"comma separated list of integers for derived streams")
	flags.String(descriptionFlag, "", "description")
	flags.StringSlice(tagsFlag, []string{}, "tags")
	flags.Bool(saveFlag, true, "if true, save the result in the config directory (~/.config/strmprivacy/saved-entities). (default is true)")
	flags.StringArrayP(maskedFieldsFlag, "M", []string{}, maskedFieldHelp)
	flags.String(maskedFieldsSeed, "", `A seed used for masking`)

	err := stream.RegisterFlagCompletionFunc(linkedStreamFlag, SourceNamesCompletion)
	err = stream.RegisterFlagCompletionFunc(maskedFieldsFlag, completion)
	common.CliExit(err)
	return stream
}

func completion(cmd *cobra.Command, args []string, complete string) ([]string, cobra.ShellCompDirective) {
	s, c := event_contract.RefsCompletion(cmd, args, complete)
	s = util.MapStrings(s, func(j string) string { return j + ":" })
	return s, c

}

func DeleteCmd() *cobra.Command {
	stream := &cobra.Command{
		Use:   "stream [name ...]",
		Short: "Delete one or more streams",
		Long: `Delete one or more streams.

	If a stream has dependents (like derived streams or exporters), you can use
	the 'recursive' option to get rid of those also.
	Returns everything that was deleted. `,
		PreRun: func(cmd *cobra.Command, args []string) {
			printer = configurePrinter(cmd)
		},
		Run: func(cmd *cobra.Command, args []string) {
			recursive, _ := cmd.Flags().GetBool(common.RecursiveFlagName)
			for i, _ := range args {
				del(&args[i], recursive)
			}
		},
		Args:              cobra.MinimumNArgs(1),
		ValidArgsFunction: NamesCompletion,
	}

	return stream
}

func GetCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "stream [name]",
		Short: "Get stream by name",
		PreRun: func(cmd *cobra.Command, args []string) {
			printer = configurePrinter(cmd)
		},
		Run: func(cmd *cobra.Command, args []string) {
			recursive, _ := cmd.Flags().GetBool(common.RecursiveFlagName)
			get(&args[0], recursive)
		},
		Args:              cobra.ExactArgs(1), // the stream name
		ValidArgsFunction: NamesCompletion,
	}
}
func ListCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "streams",
		Short: "List streams",
		PreRun: func(cmd *cobra.Command, args []string) {
			printer = configurePrinter(cmd)
		},
		Run: func(cmd *cobra.Command, args []string) {
			recursive, _ := cmd.Flags().GetBool(common.RecursiveFlagName)
			list(recursive)
		},
		ValidArgsFunction: common.NoFilesEmptyCompletion,
	}
}
