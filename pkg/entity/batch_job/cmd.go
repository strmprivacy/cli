package batch_job

import (
	"github.com/spf13/cobra"
	"strmprivacy/strm/pkg/entity/policy"
)

const (
	batchJobFileFlagName = "file"
	batchJobTypeFlagName = "type"
	encryptionType       = "encryption"
	microAggregationType = "micro-aggregation"
)

func DeleteCmd() *cobra.Command {
	batchJob := &cobra.Command{
		Use:               "batch-job (id ...)",
		Short:             "Delete one or more Batch Jobs by id",
		DisableAutoGenTag: true,
		PreRun: func(cmd *cobra.Command, args []string) {
			printer = configurePrinter(cmd)
		},
		Run: func(cmd *cobra.Command, args []string) {
			for _, arg := range args {
				del(&arg, cmd)
			}
		},
		Args:              cobra.MinimumNArgs(1),
		ValidArgsFunction: namesCompletion,
	}

	return batchJob

}

func GetCmd() *cobra.Command {
	return &cobra.Command{
		Use:               "batch-job (id)",
		Short:             "Get a Batch Job by id",
		DisableAutoGenTag: true,
		PreRun: func(cmd *cobra.Command, args []string) {
			printer = configurePrinter(cmd)
		},
		Run: func(cmd *cobra.Command, args []string) {
			get(&args[0], cmd)
		},
		Args:              cobra.ExactArgs(1),
		ValidArgsFunction: namesCompletion,
	}
}
func ListCmd() *cobra.Command {
	return &cobra.Command{
		Use:               "batch-jobs",
		Short:             "List Batch Jobs",
		DisableAutoGenTag: true,
		PreRun: func(cmd *cobra.Command, args []string) {
			printer = configurePrinter(cmd)
		},
		Run: func(cmd *cobra.Command, args []string) {
			list(cmd)
		},
	}
}

func CreateCmd() *cobra.Command {
	batchJob := &cobra.Command{
		Use:               "batch-job",
		Short:             "Create a Batch Job",
		Long:              longDoc,
		Example:           example,
		DisableAutoGenTag: true,
		PreRun: func(cmd *cobra.Command, args []string) {
			printer = configurePrinter(cmd)
		},
		Run: func(cmd *cobra.Command, args []string) {
			create(cmd)
		},
		Args: cobra.ExactArgs(0),
	}

	flags := batchJob.Flags()

	flags.StringP(batchJobFileFlagName, "F", "",
		`the path to the JSON file containing the batch job configuration`)
	flags.StringP(batchJobTypeFlagName, "T", "encryption",
		`the type of batch job (encryption, micro-aggregation), defaults to encryption`)
	policy.SetupFlags(batchJob, flags)
	_ = batchJob.MarkFlagRequired(batchJobFileFlagName)

	return batchJob
}
