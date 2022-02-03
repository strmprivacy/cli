package batch_job

import (
	"github.com/spf13/cobra"
	"strmprivacy/strm/pkg/common"
)

const (
	batch_jobs_file_flag_name = "file"
)

func DeleteCmd() *cobra.Command {
	batchJob := &cobra.Command{
		Use:   "batch-job [id]",
		Short: "Delete a Batch Job by id",
		PreRun: func(cmd *cobra.Command, args []string) {
			printer = configurePrinter(cmd)
		},
		Run: func(cmd *cobra.Command, args []string) {
			del(&args[0])
		},
		Args:              cobra.ExactArgs(1),
		ValidArgsFunction: namesCompletion,
	}
	return batchJob

}

func GetCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "batch-job [id]",
		Short: "Get a Batch Job by id",
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
		Use:   "batch-jobs",
		Short: "List Batch Jobs",
		PreRun: func(cmd *cobra.Command, args []string) {
			printer = configurePrinter(cmd)
		},
		Run: func(cmd *cobra.Command, args []string) {
			list()
		},
	}
}

func CreateCmd() *cobra.Command {
	batchJob := &cobra.Command{
		Use:   "batch-job",
		Short: "Create a Batch Job",
		PreRun: func(cmd *cobra.Command, args []string) {
			printer = configurePrinter(cmd)
		},
		Run: func(cmd *cobra.Command, args []string) {
			create(cmd)
		},
		Args: cobra.ExactArgs(0),
	}

	flags := batchJob.Flags()

	flags.StringP(batch_jobs_file_flag_name, "F", "",
		`The path to the JSON file containing the batch job configuration`)
	err := batchJob.MarkFlagRequired(batch_jobs_file_flag_name)
	common.CliExit(err)

	return batchJob
}
