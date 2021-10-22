package stream

import (
	"github.com/spf13/cobra"
	"streammachine.io/strm/pkg/common"
	"streammachine.io/strm/pkg/entity/event_contract"
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
	flags.Bool(saveFlag, false, "save the result in the config directory")
	flags.StringArrayP(maskedFieldsFlag, "M", []string{}, maskedFieldHelp)

	err := stream.RegisterFlagCompletionFunc(linkedStreamFlag, SourceNamesCompletion)
	err = stream.RegisterFlagCompletionFunc(maskedFieldsFlag, completion)
	common.CliExit(err)
	return stream
}

func completion(cmd *cobra.Command, args []string, complete string) ([]string, cobra.ShellCompDirective) {
	s, c := event_contract.RefsCompletion(cmd, args, complete)
	// add : to each of the completions
	return s, c

}

func DeleteCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "stream [name]",
		Short: "Delete a stream",
		Long: `Delete a stream.

	If a stream has dependents (like derived streams or exporters), you can use
	the 'recursive' option to get rid of those also.
	Returns everything that was deleted. `,
		PreRun: func(cmd *cobra.Command, args []string) {
			printer = configurePrinter(cmd)
		},
		Run: func(cmd *cobra.Command, args []string) {
			recursive, _ := cmd.Flags().GetBool(common.RecursiveFlagName)
			del(&args[0], recursive)
		},
		Args:              cobra.ExactArgs(1), // the stream name
		ValidArgsFunction: NamesCompletion,
	}
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