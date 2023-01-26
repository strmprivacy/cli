package stream

import (
	"github.com/spf13/cobra"
	"strmprivacy/strm/pkg/common"
	"strmprivacy/strm/pkg/entity/data_contract"
	"strmprivacy/strm/pkg/entity/policy"
	"strmprivacy/strm/pkg/util"
)

func CreateCmd() *cobra.Command {
	stream := &cobra.Command{
		Use:               "stream (name)",
		Short:             "Create a Stream",
		Long:              longDocCreate,
		Example:           createExample,
		DisableAutoGenTag: true,
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
	flags.String(descriptionFlag, "", "description of this stream")
	flags.StringSlice(tagsFlag, []string{}, "a list of tags for this stream")
	flags.Bool(saveFlag, true, "if true, save the credentials in ~/.config/strmprivacy/saved-entities.")
	flags.StringArrayP(maskedFieldsFlag, "M", []string{}, maskedFieldHelp)
	flags.String(maskedFieldsSeed, "", `a seed used for masking`)
	policy.SetupFlags(stream, flags)
	err := stream.RegisterFlagCompletionFunc(linkedStreamFlag, SourceNamesCompletion)
	err = stream.RegisterFlagCompletionFunc(maskedFieldsFlag, completion)
	common.CliExit(err)
	return stream
}

func completion(cmd *cobra.Command, args []string, complete string) ([]string, cobra.ShellCompDirective) {
	s, c := data_contract.RefsCompletion(cmd, args, complete)
	s = util.MapStrings(s, func(j string) string { return j + ":" })
	return s, c

}

func DeleteCmd() *cobra.Command {
	stream := &cobra.Command{
		Use:               "stream (name ...)",
		Short:             "Delete one or more Streams",
		DisableAutoGenTag: true,
		PreRun: func(cmd *cobra.Command, args []string) {
			printer = configurePrinter(cmd)
		},
		Run: func(cmd *cobra.Command, args []string) {
			recursive, _ := cmd.Flags().GetBool(common.RecursiveFlagName)
			for _, arg := range args {
				del(&arg, recursive, cmd)
			}
		},
		Args:              cobra.MinimumNArgs(1),
		ValidArgsFunction: NamesCompletion,
	}

	return stream
}

func GetCmd() *cobra.Command {
	return &cobra.Command{
		Use:               "stream (name)",
		Short:             "Get Stream by name",
		Example:           getExample,
		DisableAutoGenTag: true,
		PreRun: func(cmd *cobra.Command, args []string) {
			printer = configurePrinter(cmd)
		},
		Run: func(cmd *cobra.Command, args []string) {
			recursive, _ := cmd.Flags().GetBool(common.RecursiveFlagName)
			get(&args[0], recursive, cmd)
		},
		Args:              cobra.ExactArgs(1), // the stream name
		ValidArgsFunction: NamesCompletion,
	}
}
func ListCmd() *cobra.Command {
	return &cobra.Command{
		Use:               "streams",
		Short:             "List Streams",
		Example:           listExample,
		DisableAutoGenTag: true,
		PreRun: func(cmd *cobra.Command, args []string) {
			printer = configurePrinter(cmd)
		},
		Run: func(cmd *cobra.Command, args []string) {
			recursive, _ := cmd.Flags().GetBool(common.RecursiveFlagName)
			list(recursive, cmd)
		},
		ValidArgsFunction: common.NoFilesEmptyCompletion,
	}
}
