package cmd

import (
	"github.com/spf13/cobra"
	"streammachine.io/strm/entity/batch_exporter"
	"streammachine.io/strm/entity/kafka_exporter"
	"streammachine.io/strm/entity/kafka_user"
	"streammachine.io/strm/entity/sink"
	"streammachine.io/strm/entity/stream"
)

// CreateCmd represents the create command
var CreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create an entity",
}

func init() {
	CreateCmd.AddCommand(stream.CreateCmd())
	CreateCmd.AddCommand(sink.CreateCmd())
	CreateCmd.AddCommand(batch_exporter.CreateCmd())
	CreateCmd.AddCommand(kafka_exporter.CreateCmd())
	CreateCmd.AddCommand(kafka_user.CreateCmd())
}
