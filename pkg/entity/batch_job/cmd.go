package batch_job

import (
	"github.com/spf13/cobra"
	"strmprivacy/strm/pkg/common"
)

const (
	batch_jobs_file_flag_name = "file"
)

var longDoc = `
A Batch Job outputs all events in file all events to files in a Sink. This happens with a regular interval.

Each file follows the JSON Lines format, which is one full JSON document per line.

A [sink](sink.md) is a configuration item that defines location
(Gcloud, AWS, ..) bucket and associated credentials.

A sink needs to be created *before* you can create a batch job that uses it

### Usage
`

func DeleteCmd() *cobra.Command {
	batchJob := &cobra.Command{
		Use:               "batch-job [id ...]",
		Short:             "Delete on or more Batch Jobs by id",
		Long:              longDoc,
		DisableAutoGenTag: true,
		PreRun: func(cmd *cobra.Command, args []string) {
			printer = configurePrinter(cmd)
		},
		Run: func(cmd *cobra.Command, args []string) {
			for i, _ := range args {
				del(&args[i])
			}
		},
		Args:              cobra.MinimumNArgs(1),
		ValidArgsFunction: namesCompletion,
	}

	return batchJob

}

func GetCmd() *cobra.Command {
	return &cobra.Command{
		Use:               "batch-job [id]",
		Short:             "Get a Batch Job by id",
		Long:              longDoc,
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
		Long:              longDoc,
		DisableAutoGenTag: true,
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
		Use:               "batch-job",
		Short:             "Create a Batch Job",
		Long:              longDoc,
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

	flags.StringP(batch_jobs_file_flag_name, "F", "",
		`The path to the JSON file containing the batch job configuration`)
	err := batchJob.MarkFlagRequired(batch_jobs_file_flag_name)
	common.CliExit(err)

	return batchJob
}
