package batch_exporter

import (
	"github.com/spf13/cobra"
	"strmprivacy/strm/pkg/common"
	"strmprivacy/strm/pkg/entity/data_connector"
	"strmprivacy/strm/pkg/entity/stream"
	"strmprivacy/strm/pkg/util"
)

const (
	dataConnectorFlag         = "data-connector"
	nameFlag                  = "name"
	intervalFlag              = "interval"
	pathPrefix                = "path-prefix"
	exportKeys                = "export-keys"
	includeExistingEventsFlag = "include-existing-events"
)

var longDoc = util.LongDocsUsage(`
A Batch Exporter listens to a stream and writes all events to files using a Data Connector. This happens with a regular
interval.

When exporting events each file follows the JSON Lines format, which is one full JSON document per line.
When exporting encryption keys, each file is a CSV file.

A [Data Connector](docs/04-reference/01-cli-reference/!strm/create/data-connector.md) is a configuration
entity that comprises a location (GCS bucket, AWS S3 bucket, ...) and associated credentials.

A Data Connector must be created *before* you can create a batch exporter that uses it.



`)

func DeleteCmd() *cobra.Command {
	batchExporter := &cobra.Command{
		Use:   "batch-exporter (name ...)",
		Short: "Delete one or more Batch Exporters by name",
		PreRun: func(cmd *cobra.Command, args []string) {
			printer = configurePrinter(cmd)
		},
		Run: func(cmd *cobra.Command, args []string) {
			for _, arg := range args {
				del(&arg, cmd)
			}
		},
		Args:              cobra.MinimumNArgs(1), // the stream names
		DisableAutoGenTag: true,
		ValidArgsFunction: NamesCompletion,
	}

	return batchExporter
}

func GetCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "batch-exporter (name)",
		Short: "Get Batch Exporter by name",
		PreRun: func(cmd *cobra.Command, args []string) {
			printer = configurePrinter(cmd)
		},
		Run: func(cmd *cobra.Command, args []string) {
			get(&args[0], cmd)
		},
		Args:              cobra.ExactArgs(1), // the stream name
		DisableAutoGenTag: true,
		ValidArgsFunction: NamesCompletion,
	}
}

func ListCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "batch-exporters",
		Short: "List Batch Exporters",
		PreRun: func(cmd *cobra.Command, args []string) {
			printer = configurePrinter(cmd)
		},
		DisableAutoGenTag: true,
		Run: func(cmd *cobra.Command, args []string) {
			list(cmd)
		},
	}
}

func CreateCmd() *cobra.Command {
	batchExporter := &cobra.Command{
		Use:   "batch-exporter (stream-name)",
		Short: "Create a Batch Exporter",
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
	flags.String(nameFlag, "", "optional batch exporter name (default <stream>-<dataconnector>).")
	flags.String(pathPrefix, "", "path prefix on bucket")
	flags.Int64(intervalFlag, 60, "interval in seconds between batches")
	flags.Bool(exportKeys, false, "do we want to export the keys stream")
	flags.Bool(includeExistingEventsFlag, false, "do we want to include all existing events")
	err := batchExporter.RegisterFlagCompletionFunc(dataConnectorFlag, data_connector.NamesCompletion)
	common.CliExit(err)

	return batchExporter
}
