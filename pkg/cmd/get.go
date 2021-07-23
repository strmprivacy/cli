package cmd

import (
	"github.com/spf13/cobra"
	"streammachine.io/strm/pkg/entity/batch_exporter"
	"streammachine.io/strm/pkg/entity/event_contract"
	"streammachine.io/strm/pkg/entity/kafka_cluster"
	"streammachine.io/strm/pkg/entity/kafka_exporter"
	"streammachine.io/strm/pkg/entity/kafka_user"
	"streammachine.io/strm/pkg/entity/key_stream"
	"streammachine.io/strm/pkg/entity/schema"
	"streammachine.io/strm/pkg/entity/sink"
	"streammachine.io/strm/pkg/entity/stream"
	"streammachine.io/strm/pkg/entity/usage"
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
	GetCmd.AddCommand(usage.GetCmd())
	GetCmd.PersistentFlags().BoolP("recursive", "r", false, "recursive")
}