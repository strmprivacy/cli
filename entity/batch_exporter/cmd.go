package batch_exporter

import (
	"github.com/spf13/cobra"
	"streammachine.io/strm/entity/sink"
	"streammachine.io/strm/entity/stream"
)

const (
	sinkFlag     = "sink"
	nameFlag     = "name"
	intervalFlag = "interval"
	pathPrefix   = "path-prefix"
	exportKeys   = "export-keys"
)

func DeleteCmd() *cobra.Command {
	batchExporter := &cobra.Command{
		Use:   "batch-exporter [name]",
		Short: "delete batch-exporter by name",
		Run: func(cmd *cobra.Command, args []string) {
			del(&args[0])
		},
		Args:              cobra.ExactArgs(1), // the stream name
		ValidArgsFunction: existingNamesCompletion,
	}
	return batchExporter

}

func GetCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "batch-exporter [name]",
		Short: "get batch-exporter by name",
		Run: func(cmd *cobra.Command, args []string) {
			get(&args[0], cmd)
		},
		Args:              cobra.ExactArgs(1), // the stream name
		ValidArgsFunction: existingNamesCompletion,
	}
}
func ListCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "batch-exporters",
		Short: "List batch-exporters",
		Run: func(cmd *cobra.Command, args []string) {
			list()
		},
	}
}

func CreateCmd() *cobra.Command {
	batchExporter := &cobra.Command{
		Use:   "batch-exporter [stream-name]",
		Short: "Create batch exporter",
		Run: func(cmd *cobra.Command, args []string) {
			create(&args[0], cmd)
		},
		Args:              cobra.ExactArgs(1), // the stream name
		ValidArgsFunction: stream.ExistingNamesCompletion,
	}
	flags := batchExporter.Flags()
	flags.String(sinkFlag, "", "name of the sink. Optional if you have only one defined sink.")
	flags.String(nameFlag, "", "optional batch exporter name")
	flags.String(pathPrefix, "", "path prefix on bucket")
	flags.Int64(intervalFlag, 60, "Interval in seconds between batches")
	flags.Bool(exportKeys, false, "Do we want to export the keys stream")
	err := batchExporter.RegisterFlagCompletionFunc(sinkFlag, sink.ExistingNamesCompletion)
	cobra.CheckErr(err)
	return batchExporter

}
