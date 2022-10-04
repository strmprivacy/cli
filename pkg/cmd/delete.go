package cmd

import (
	"github.com/spf13/cobra"
	"strmprivacy/strm/pkg/auth"
	"strmprivacy/strm/pkg/common"
	"strmprivacy/strm/pkg/entity/batch_exporter"
	"strmprivacy/strm/pkg/entity/batch_job"
	"strmprivacy/strm/pkg/entity/data_connector"
	"strmprivacy/strm/pkg/entity/data_contract"
	"strmprivacy/strm/pkg/entity/data_subject"
	"strmprivacy/strm/pkg/entity/event_contract"
	"strmprivacy/strm/pkg/entity/kafka_exporter"
	"strmprivacy/strm/pkg/entity/kafka_user"
	"strmprivacy/strm/pkg/entity/project"
	"strmprivacy/strm/pkg/entity/schema"
	"strmprivacy/strm/pkg/entity/stream"
)

// DeleteCmd represents the delete command
var DeleteCmd = &cobra.Command{
	Use:               common.DeleteCommandName,
	PersistentPreRun:  auth.RequireAuthenticationPreRun,
	DisableAutoGenTag: true,
	Short:             "Delete an entity",
}

func init() {
	DeleteCmd.AddCommand(stream.DeleteCmd())
	DeleteCmd.AddCommand(kafka_exporter.DeleteCmd())
	DeleteCmd.AddCommand(batch_exporter.DeleteCmd())
	DeleteCmd.AddCommand(batch_job.DeleteCmd())
	DeleteCmd.AddCommand(data_connector.DeleteCmd())
	DeleteCmd.AddCommand(kafka_user.DeleteCmd())
	DeleteCmd.AddCommand(event_contract.DeleteCmd())
	DeleteCmd.AddCommand(schema.DeleteCmd())
	DeleteCmd.AddCommand(data_subject.DeleteCmd())
	DeleteCmd.AddCommand(data_contract.DeleteCmd())
	DeleteCmd.AddCommand(project.DeleteCmd())

	DeleteCmd.PersistentFlags().BoolP(common.RecursiveFlagName, common.RecursiveFlagShorthand, false, common.RecursiveFlagUsage)
}
