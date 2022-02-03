package cmd

import (
	"github.com/spf13/cobra"
	"strmprivacy/strm/pkg/common"
	"strmprivacy/strm/pkg/entity/batch_exporter"
	"strmprivacy/strm/pkg/entity/batch_job"
	"strmprivacy/strm/pkg/entity/event_contract"
	"strmprivacy/strm/pkg/entity/kafka_cluster"
	"strmprivacy/strm/pkg/entity/kafka_exporter"
	"strmprivacy/strm/pkg/entity/kafka_user"
	"strmprivacy/strm/pkg/entity/key_stream"
	"strmprivacy/strm/pkg/entity/schema"
	"strmprivacy/strm/pkg/entity/schema_code"
	"strmprivacy/strm/pkg/entity/sink"
	"strmprivacy/strm/pkg/entity/stream"
	"strmprivacy/strm/pkg/entity/usage"
)

var GetCmd = &cobra.Command{
	Use:   common.GetCommandName,
	Short: "Get an entity",
}

func init() {
	GetCmd.AddCommand(stream.GetCmd())
	GetCmd.AddCommand(kafka_exporter.GetCmd())
	GetCmd.AddCommand(batch_exporter.GetCmd())
	GetCmd.AddCommand(batch_job.GetCmd())
	GetCmd.AddCommand(sink.GetCmd())
	GetCmd.AddCommand(kafka_cluster.GetCmd())
	GetCmd.AddCommand(kafka_user.GetCmd())
	GetCmd.AddCommand(key_stream.GetCmd())
	GetCmd.AddCommand(schema.GetCmd())
	GetCmd.AddCommand(schema_code.GetCmd())
	GetCmd.AddCommand(event_contract.GetCmd())
	GetCmd.AddCommand(usage.GetCmd())

	GetCmd.PersistentFlags().BoolP(common.RecursiveFlagName, common.RecursiveFlagShorthand, false, common.RecursiveFlagUsage)
}
