package cmd

import (
	"github.com/spf13/cobra"
	"streammachine.io/strm/pkg/entity/batch_exporter"
	"streammachine.io/strm/pkg/entity/kafka_exporter"
	"streammachine.io/strm/pkg/entity/kafka_user"
	"streammachine.io/strm/pkg/entity/sink"
	"streammachine.io/strm/pkg/entity/stream"
)

// DeleteCmd represents the delete command
var DeleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete an entity",
}

func init() {
	DeleteCmd.PersistentFlags().BoolP("recursive", "r", false, "recursive")
	DeleteCmd.AddCommand(stream.DeleteCmd())
	DeleteCmd.AddCommand(kafka_exporter.DeleteCmd())
	DeleteCmd.AddCommand(batch_exporter.DeleteCmd())
	DeleteCmd.AddCommand(sink.DeleteCmd())
	DeleteCmd.AddCommand(kafka_user.DeleteCmd())
}
