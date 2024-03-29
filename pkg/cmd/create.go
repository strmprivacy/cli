package cmd

import (
	"github.com/spf13/cobra"
	"strmprivacy/strm/pkg/auth"
	"strmprivacy/strm/pkg/common"
	"strmprivacy/strm/pkg/entity/batch_exporter"
	"strmprivacy/strm/pkg/entity/batch_job"
	"strmprivacy/strm/pkg/entity/data_connector"
	"strmprivacy/strm/pkg/entity/data_contract"
	"strmprivacy/strm/pkg/entity/kafka_exporter"
	"strmprivacy/strm/pkg/entity/kafka_user"
	"strmprivacy/strm/pkg/entity/policy"
	"strmprivacy/strm/pkg/entity/project"
	"strmprivacy/strm/pkg/entity/purpose_mapping"
	"strmprivacy/strm/pkg/entity/stream"
)

// CreateCmd represents the create command
var CreateCmd = &cobra.Command{
	Use:               common.CreateCommandName,
	PersistentPreRun:  auth.RequireAuthenticationPreRun,
	DisableAutoGenTag: true,
	Short:             "Create an entity",
}

func init() {
	CreateCmd.AddCommand(stream.CreateCmd())
	CreateCmd.AddCommand(data_connector.CreateCmd())
	CreateCmd.AddCommand(batch_exporter.CreateCmd())
	CreateCmd.AddCommand(batch_job.CreateCmd())
	CreateCmd.AddCommand(data_contract.CreateCmd())
	CreateCmd.AddCommand(kafka_exporter.CreateCmd())
	CreateCmd.AddCommand(kafka_user.CreateCmd())
	CreateCmd.AddCommand(project.CreateCmd())
	CreateCmd.AddCommand(policy.CreateCmd())
	CreateCmd.AddCommand(purpose_mapping.CreateCmd())
}
