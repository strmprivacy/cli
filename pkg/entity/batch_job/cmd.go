package batch_job

import (
	"fmt"
	"github.com/spf13/cobra"
	"strmprivacy/strm/pkg/common"
	"strmprivacy/strm/pkg/entity/policy"
)

const (
	batchJobFileFlagName = "file"
	batchJobTypeFlagName = "type"
	encryptionType = "encryption"
	microAggregationType = "micro-aggregation"
)

// Todo: point to detailed docs/quickstarts on Batch Jobs
var longDoc = `
A Batch Job reads all events from a Data Connector and writes them to one or more Data Connectors,
applying one of our privacy algorithms as defined by the job's configuration file. An encryption batch job
encrypts sensitive data, while a micro-aggregation batch job applies k-member clustering and replaces
the values of quasi identifier fields with an aggregated value (e.g. mean value of a cluster). 

A [Data Connector](docs/04-reference/01-cli-reference/` + fmt.Sprint(common.RootCommandName) + `/create/data-connector.md) is a configuration
entity that comprises a location (GCS bucket, AWS S3 bucket, ...) and associated credentials.

A Data Connector must be created in the same project *before* you can create a batch job that uses it.

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
			list(cmd)
		},
	}
}

func CreateCmd() *cobra.Command {
	batchJob := &cobra.Command{
		Use:               "batch-job",
		Short:             "Create a Batch Job",
		Long:              longDoc,
		DisableAutoGenTag: true,
		Example: "strm create batch-job --type encryption --file my_config.json",
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
		`The path to the JSON file containing the batch job configuration`)
	flags.StringP(batchJobTypeFlagName, "T", "encryption",
		`The type of batch job (encryption, micro-aggregation), defaults to encryption`)
	policy.SetupFlags(batchJob, flags)
	_ = batchJob.MarkFlagRequired(batchJobFileFlagName)

	return batchJob
}
