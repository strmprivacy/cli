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

var GetCmd = &cobra.Command{
	Use:   "get",
	Short: "Get an entity",
}

func init() {
	GetCmd.AddCommand(stream.GetCmd())
	GetCmd.AddCommand(kafka_exporter.GetCmd())
	GetCmd.AddCommand(batch_exporter.GetCmd())
	GetCmd.AddCommand(sink.GetCmd())
	GetCmd.AddCommand(kafka_cluster.GetCmd())
	GetCmd.AddCommand(kafka_user.GetCmd())
	GetCmd.AddCommand(key_stream.GetCmd())
	GetCmd.AddCommand(schema.GetCmd())
	GetCmd.AddCommand(event_contract.GetCmd())
	GetCmd.PersistentFlags().BoolP("recursive", "r", false, "recursive")
}
