package stream

import (
	"github.com/spf13/cobra"
	"strmprivacy/strm/pkg/common"
	"strmprivacy/strm/pkg/entity/event_contract"
	"strmprivacy/strm/pkg/entity/policy"
	"strmprivacy/strm/pkg/util"
)

var longDoc = `A stream is the central resource in STRM Privacy. Clients can connect to a stream to send and to receive events. A
stream can be either an "input stream" or a "derived stream".

Events are always sent to an input stream. Sending events to a derived stream is not possible. After validation and
encryption of all PII fields, STRM Privacy sends all events to the input stream. Clients consuming from the input stream
will see all events, but with all PII fields encrypted.

Derived streams can be made on top of a input stream. A derived stream is configured with one or more consent levels and
it only receives events with matching consent levels (see details about this matching process here). The PII fields with
matching consent levels are decrypted and sent to the derived stream. Clients connecting to the derived stream will only
receive the events on this stream.

Every stream has its own set of access tokens. Connecting to an input stream requires different credentials than when
connecting to a derived stream.

### Usage`

func CreateCmd() *cobra.Command {
	stream := &cobra.Command{
		Use:               "stream [name]",
		Short:             "Create a stream",
		Long:              longDoc,
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
	flags.String(descriptionFlag, "", "description")
	flags.StringSlice(tagsFlag, []string{}, "tags")
	flags.Bool(saveFlag, true, "if true, save the result in the config directory (~/.config/strmprivacy/saved-entities). (default is true)")
	flags.StringArrayP(maskedFieldsFlag, "M", []string{}, maskedFieldHelp)
	flags.String(maskedFieldsSeed, "", `A seed used for masking`)
	policy.SetupFlags(stream, flags)
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
		Use:               "stream [name ...]",
		Short:             "Delete one or more streams",
		Long:              longDoc,
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
		Use:               "stream [name]",
		Short:             "Get stream by name",
		Long:              longDoc,
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
		Short:             "List streams",
		Long:              longDoc,
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
