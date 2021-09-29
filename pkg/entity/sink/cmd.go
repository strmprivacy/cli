package sink

import (
	"github.com/spf13/cobra"
	"streammachine.io/strm/pkg/common"
)

const (
	sinkTypeFlag        = "sink-type"
	credentialsFileFlag = "credentials-file"
)

func GetCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "sink [name]",
		Short: "Get sink by name",
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
		Use:   "sinks",
		Short: "List sinks",
		PreRun: func(cmd *cobra.Command, args []string) {
			printer = configurePrinter(cmd)
		},
		Run: func(cmd *cobra.Command, args []string) {
			recursive, _ := cmd.Flags().GetBool(common.RecursiveFlagName)
			list(recursive)
		},
	}
}
func DeleteCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "sink [name]",
		Short: "Delete sinks",
		PreRun: func(cmd *cobra.Command, args []string) {
			printer = configurePrinter(cmd)
		},
		Run: func(cmd *cobra.Command, args []string) {
			recursive, _ := cmd.Flags().GetBool(common.RecursiveFlagName)
			del(&args[0], recursive)
		},
		Args:              cobra.ExactArgs(1),
		ValidArgsFunction: NamesCompletion,
	}
}
func CreateCmd() *cobra.Command {
	sink := &cobra.Command{
		Use:   "sink [sink-name] [bucket-name]",
		Short: "Create sink",
		PreRun: func(cmd *cobra.Command, args []string) {
			printer = configurePrinter(cmd)
		},
		Run: func(cmd *cobra.Command, args []string) {
			create(&args[0], &args[1], cmd)
		},
		Args: cobra.ExactArgs(2),
	}

	flags := sink.Flags()
	flags.String(sinkTypeFlag, "", "S3 or GCLOUD")
	flags.String(credentialsFileFlag, "", "file with credentials")
	_ = sink.MarkFlagRequired(credentialsFileFlag)
	_ = sink.MarkFlagFilename(credentialsFileFlag, "json")
	err := sink.RegisterFlagCompletionFunc(sinkTypeFlag, sinkTypesCompletion)
	common.CliExit(err)
	return sink
}

func sinkTypesCompletion(cmd *cobra.Command, args []string, complete string) ([]string, cobra.ShellCompDirective) {
	return []string{"S3", "GCLOUD"}, cobra.ShellCompDirectiveNoFileComp
}
