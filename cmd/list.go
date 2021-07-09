package cmd

import (
	"github.com/spf13/cobra"
	"streammachine.io/strm/entity/batch_exporter"
	"streammachine.io/strm/entity/event_contract"
	"streammachine.io/strm/entity/kafka_cluster"
	"streammachine.io/strm/entity/kafka_exporter"
	"streammachine.io/strm/entity/kafka_user"
	"streammachine.io/strm/entity/key_stream"
	"streammachine.io/strm/entity/schema"
	"streammachine.io/strm/entity/sink"
	"streammachine.io/strm/entity/stream"
)

var ListCmd = &cobra.Command{
	Use:   "list",
	Short: "List entities",
}

func init() {
	ListCmd.AddCommand(stream.ListCmd())
	ListCmd.AddCommand(kafka_exporter.ListCmd())
	ListCmd.AddCommand(batch_exporter.ListCmd())
	ListCmd.AddCommand(sink.ListCmd())
	ListCmd.AddCommand(kafka_cluster.ListCmd())
	ListCmd.AddCommand(kafka_user.ListCmd())
	ListCmd.AddCommand(key_stream.ListCmd())
	ListCmd.AddCommand(schema.ListCmd())
	ListCmd.AddCommand(event_contract.ListCmd())
	ListCmd.PersistentFlags().BoolP("recursive", "r", false, "recursive")
}
