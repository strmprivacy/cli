package batch_exporter

import (
	"github.com/spf13/cobra"
	"io/ioutil"
	"strmprivacy/strm/pkg/common"
	"strmprivacy/strm/pkg/entity/sink"
	"strmprivacy/strm/pkg/entity/stream"
)

const (
	sinkFlag                  = "sink"
	nameFlag                  = "name"
	intervalFlag              = "interval"
	pathPrefix                = "path-prefix"
	exportKeys                = "export-keys"
	includeExistingEventsFlag = "include-existing-events"
)

var content, _ = ioutil.ReadFile("pkg/entity/batch_exporter/docstring.md")

func DeleteCmd() *cobra.Command {
	batchExporter := &cobra.Command{
		Use:   "batch-exporter [name ...]",
		Short: "Delete one or more Batch exporters by name",
		Long:  string(content),
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
		Long:  string(content),
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
		Long:  string(content),
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
	//content, _ := ioutil.("pkg/entity/batch_exporter/docstring.md")
	batchExporter := &cobra.Command{
		Use:   "batch-exporter [stream-name]",
		Short: "Create batch exporter",
		Long:  string(content),
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
	flags.String(sinkFlag, "", "name of the sink. Optional if you have only one defined sink.")
	flags.String(nameFlag, "", "optional batch exporter name")
	flags.String(pathPrefix, "", "path prefix on bucket")
	flags.Int64(intervalFlag, 60, "Interval in seconds between batches")
	flags.Bool(exportKeys, false, "Do we want to export the keys stream")
	flags.Bool(includeExistingEventsFlag, false, "Do we want to include all existing events")
	err := batchExporter.RegisterFlagCompletionFunc(sinkFlag, sink.NamesCompletion)
	common.CliExit(err)

	return batchExporter
}
