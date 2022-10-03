package batch_exporter

import (
	"fmt"
	"github.com/spf13/cobra"
	"strmprivacy/strm/pkg/common"
	"strmprivacy/strm/pkg/entity/data_connector"
	"strmprivacy/strm/pkg/entity/stream"
)

const (
	dataConnectorFlag         = "data-connector"
	nameFlag                  = "name"
	intervalFlag              = "interval"
	pathPrefix                = "path-prefix"
	exportKeys                = "export-keys"
	includeExistingEventsFlag = "include-existing-events"
	projectName               = "project"
)

var longDoc = `
A Batch Exporter listens to a stream and writes all events to files using a Data Connector. This happens with a regular interval.

Each file follows the JSON Lines format, which is one full JSON document per line.

A [Data Connector](docs/04-reference/01-cli-reference/` + fmt.Sprint(common.RootCommandName) + `/create/data-connector.md) is a configuration
entity that comprises a location (GCS bucket, AWS S3 bucket, ...) and associated credentials.

A Data Connector must be created *before* you can create a batch exporter that uses it.

### Usage
`

func DeleteCmd() *cobra.Command {
	batchExporter := &cobra.Command{
		Use:   "batch-exporter [name ...]",
		Short: "Delete one or more Batch exporters by name",
		Long:  longDoc,
		PreRun: func(cmd *cobra.Command, args []string) {
			printer = configurePrinter(cmd)
		},
		Run: func(cmd *cobra.Command, args []string) {
			for i, _ := range args {
				del(&args[i])
			}
		},
		Args:              cobra.MinimumNArgs(1), // the stream names
		DisableAutoGenTag: true,
		ValidArgsFunction: namesCompletion,
	}

	return batchExporter
}

func GetCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "batch-exporter [name]",
		Short: "Get Batch exporter by name",
		Long:  longDoc,
		PreRun: func(cmd *cobra.Command, args []string) {
			printer = configurePrinter(cmd)
		},
		Run: func(cmd *cobra.Command, args []string) {
			get(&args[0], cmd)
		},
		Args:              cobra.ExactArgs(1), // the stream name
		DisableAutoGenTag: true,
		ValidArgsFunction: namesCompletion,
	}
}

func ListCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "batch-exporters",
		Short: "List Batch exporters",
		Long:  longDoc,
		PreRun: func(cmd *cobra.Command, args []string) {
			printer = configurePrinter(cmd)
		},
		DisableAutoGenTag: true,
		Run: func(cmd *cobra.Command, args []string) {
			list()
		},
	}
}

func CreateCmd() *cobra.Command {
	batchExporter := &cobra.Command{
		Use:   "batch-exporter [stream-name]",
		Short: "Create batch exporter",
		Long:  longDoc,
		PreRun: func(cmd *cobra.Command, args []string) {
			printer = configurePrinter(cmd)
		},
		DisableAutoGenTag: true,
		Run: func(cmd *cobra.Command, args []string) {
			create(&args[0], cmd)
		},
		Args:              cobra.ExactArgs(1), // the stream name
		ValidArgsFunction: stream.NamesCompletion,
	}

	flags := batchExporter.Flags()
	flags.String(dataConnectorFlag, "", "name of the data connector - optional if you own only one data connector")
	flags.String(nameFlag, "", "optional batch exporter name")
	flags.String(pathPrefix, "", "path prefix on bucket")
	flags.Int64(intervalFlag, 60, "Interval in seconds between batches")
	flags.Bool(exportKeys, false, "Do we want to export the keys stream")
	flags.Bool(includeExistingEventsFlag, false, "Do we want to include all existing events")
	flags.String(projectName, "", `Project name to create resource in`)
	err := batchExporter.RegisterFlagCompletionFunc(dataConnectorFlag, data_connector.NamesCompletion)
	common.CliExit(err)

	return batchExporter
}
