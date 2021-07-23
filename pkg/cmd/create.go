package cmd

import (
	"github.com/spf13/cobra"
	"streammachine.io/strm/pkg/entity/batch_exporter"
	"streammachine.io/strm/pkg/entity/event_contract"
	"streammachine.io/strm/pkg/entity/kafka_exporter"
	"streammachine.io/strm/pkg/entity/kafka_user"
	"streammachine.io/strm/pkg/entity/schema"
	"streammachine.io/strm/pkg/entity/sink"
	"streammachine.io/strm/pkg/entity/stream"
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
	CreateCmd.AddCommand(schema.CreateCmd())
	CreateCmd.AddCommand(event_contract.CreateCmd())
}