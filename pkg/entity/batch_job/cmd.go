package batch_job

import (
	"fmt"
	"github.com/spf13/cobra"
	"strmprivacy/strm/pkg/common"
)

const (
	batchJobsFileFlagName = "file"
)

var longDoc = `
A Batch Job reads all events from a Data Connector and writes them to one or more Data Connectors,
applying our privacy algorithm as defined by the job's configuration file.

A [Data Connector](/cli-reference/` + fmt.Sprint(common.RootCommandName) + `/create/data-connector.md) is a configuration
entity that comprises a location (GCS bucket, AWS S3 bucket, ...) and associated credentials.

A Data Connector must be created *before* you can create a batch job that uses it.

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

	flags.StringP(batchJobsFileFlagName, "F", "",
		`The path to the JSON file containing the batch job configuration`)
	err := batchJob.MarkFlagRequired(batchJobsFileFlagName)
	common.CliExit(err)

	return batchJob
}
