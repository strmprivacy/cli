package sink

import (
	"github.com/spf13/cobra"
	"strmprivacy/strm/pkg/common"
)

const (
	sinkTypeFlag        = "sink-type"
	credentialsFileFlag = "credentials-file"
	assumeRoleArnFlag   = "assume-role-arn"
)

var longDoc = `A Sink is a STRM Privacy configuration object for a remote file storage. For now, AWS S3 and Google Cloud Storage
Buckets are supported. By itself, a Sink does nothing. A Batch Exporter needs to be connected to a Sink and a Stream to
start outputting events.

Upon creation, STRM Privacy validates whether or not the Bucket exists and if it is accessible with the given
credentials.

### Usage`

func GetCmd() *cobra.Command {
	return &cobra.Command{
		Use:               "sink [name]",
		Short:             "Get sink by name",
		Long:              longDoc,
		DisableAutoGenTag: true,
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
		Use:               "sinks",
		Short:             "List sinks",
		Long:              longDoc,
		DisableAutoGenTag: true,
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
		Use:               "sink [name ...]",
		Short:             "Delete sinks",
		Long:              longDoc,
		DisableAutoGenTag: true,
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
}
func CreateCmd() *cobra.Command {
	sink := &cobra.Command{
		Use:               "sink [sink-name] [bucket-name]",
		Short:             "Create sink",
		Long:              longDoc,
		DisableAutoGenTag: true,
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
	flags.String(assumeRoleArnFlag, "", "ARN of the role to assume")
	_ = sink.MarkFlagRequired(credentialsFileFlag)
	_ = sink.MarkFlagFilename(credentialsFileFlag, "json")
	err := sink.RegisterFlagCompletionFunc(sinkTypeFlag, sinkTypesCompletion)
	common.CliExit(err)
	return sink
}

func sinkTypesCompletion(cmd *cobra.Command, args []string, complete string) ([]string, cobra.ShellCompDirective) {
	return []string{"S3", "GCLOUD"}, cobra.ShellCompDirectiveNoFileComp
}
