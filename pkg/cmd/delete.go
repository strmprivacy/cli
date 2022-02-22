package cmd

import (
	"github.com/spf13/cobra"
	"strmprivacy/strm/pkg/common"
	"strmprivacy/strm/pkg/entity/batch_exporter"
	"strmprivacy/strm/pkg/entity/batch_job"
	"strmprivacy/strm/pkg/entity/kafka_exporter"
	"strmprivacy/strm/pkg/entity/kafka_user"
	"strmprivacy/strm/pkg/entity/sink"
	"strmprivacy/strm/pkg/entity/stream"
)

// DeleteCmd represents the delete command
var DeleteCmd = &cobra.Command{
	Use:               common.DeleteCommandName,
	DisableAutoGenTag: true,
	Short:             "Delete an entity",
}

func init() {
	DeleteCmd.AddCommand(stream.DeleteCmd())
	DeleteCmd.AddCommand(kafka_exporter.DeleteCmd())
	DeleteCmd.AddCommand(batch_exporter.DeleteCmd())
	DeleteCmd.AddCommand(batch_job.DeleteCmd())
	DeleteCmd.AddCommand(sink.DeleteCmd())
	DeleteCmd.AddCommand(kafka_user.DeleteCmd())

	DeleteCmd.PersistentFlags().BoolP(common.RecursiveFlagName, common.RecursiveFlagShorthand, false, common.RecursiveFlagUsage)
}
