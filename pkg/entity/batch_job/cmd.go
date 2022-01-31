package batch_job

import (
	"github.com/spf13/cobra"
	"strmprivacy/strm/pkg/entity/stream"
)

const (
	file = "derived-data-file"
)

func DeleteCmd() *cobra.Command {
	batchJob := &cobra.Command{
		Use:   "batch-job [id]",
		Short: "Delete a Batch job by id",
		PreRun: func(cmd *cobra.Command, args []string) {
			printer = configurePrinter(cmd)
		},
		Run: func(cmd *cobra.Command, args []string) {
			del(&args[0])
		},
		Args:              cobra.ExactArgs(1), // the stream name
		ValidArgsFunction: namesCompletion,
	}
	return batchJob

}

func GetCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "batch-job [id]",
		Short: "Get Batch job by name",
		PreRun: func(cmd *cobra.Command, args []string) {
			printer = configurePrinter(cmd)
		},
		Run: func(cmd *cobra.Command, args []string) {
			get(&args[0], cmd)
		},
		Args:              cobra.ExactArgs(1), // the stream name
		ValidArgsFunction: namesCompletion,
	}
}
func ListCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "batch-jobs",
		Short: "List Batch jobs",
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
		Short: "Create batch job",
		PreRun: func(cmd *cobra.Command, args []string) {
			printer = configurePrinter(cmd)
		},
		Run: func(cmd *cobra.Command, args []string) {
			create(cmd)
		},
		Args:              cobra.ExactArgs(0), // the stream name
		ValidArgsFunction: stream.NamesCompletion,
	}

	flags := batchJob.Flags()
	flags.StringP(file, "F", "",
		`The path to the json file with the batch job properties`)

	return batchJob
}
