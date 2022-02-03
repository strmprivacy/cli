package cmd

import (
	"github.com/spf13/cobra"
	"strmprivacy/strm/pkg/common"
	"strmprivacy/strm/pkg/entity/batch_exporter"
	"strmprivacy/strm/pkg/entity/batch_job"
	"strmprivacy/strm/pkg/entity/event_contract"
	"strmprivacy/strm/pkg/entity/kafka_exporter"
	"strmprivacy/strm/pkg/entity/kafka_user"
	"strmprivacy/strm/pkg/entity/schema"
	"strmprivacy/strm/pkg/entity/sink"
	"strmprivacy/strm/pkg/entity/stream"
)

// CreateCmd represents the create command
var CreateCmd = &cobra.Command{
	Use:   common.CreateCommandName,
	Short: "Create an entity",
}

func init() {
	CreateCmd.AddCommand(stream.CreateCmd())
	CreateCmd.AddCommand(sink.CreateCmd())
	CreateCmd.AddCommand(batch_exporter.CreateCmd())
	CreateCmd.AddCommand(batch_job.CreateCmd())
	CreateCmd.AddCommand(kafka_exporter.CreateCmd())
	CreateCmd.AddCommand(kafka_user.CreateCmd())
	CreateCmd.AddCommand(schema.CreateCmd())
	CreateCmd.AddCommand(event_contract.CreateCmd())
}
