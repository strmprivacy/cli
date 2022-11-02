package cmd

import (
	"github.com/spf13/cobra"
	"strmprivacy/strm/pkg/auth"
	"strmprivacy/strm/pkg/common"
	"strmprivacy/strm/pkg/entity/batch_exporter"
	"strmprivacy/strm/pkg/entity/batch_job"
	"strmprivacy/strm/pkg/entity/data_connector"
	"strmprivacy/strm/pkg/entity/data_contract"
	"strmprivacy/strm/pkg/entity/installation"
	"strmprivacy/strm/pkg/entity/kafka_cluster"
	"strmprivacy/strm/pkg/entity/kafka_exporter"
	"strmprivacy/strm/pkg/entity/kafka_user"
	"strmprivacy/strm/pkg/entity/key_stream"
	"strmprivacy/strm/pkg/entity/member"
	"strmprivacy/strm/pkg/entity/policy"
	"strmprivacy/strm/pkg/entity/project"
	"strmprivacy/strm/pkg/entity/schema_code"
	"strmprivacy/strm/pkg/entity/stream"
	"strmprivacy/strm/pkg/entity/usage"
)

var GetCmd = &cobra.Command{
	Use:               common.GetCommandName,
	PersistentPreRun:  auth.RequireAuthenticationPreRun,
	DisableAutoGenTag: true,
	Short:             "Get an entity",
}

func init() {
	GetCmd.AddCommand(stream.GetCmd())
	GetCmd.AddCommand(kafka_exporter.GetCmd())
	GetCmd.AddCommand(batch_exporter.GetCmd())
	GetCmd.AddCommand(batch_job.GetCmd())
	GetCmd.AddCommand(data_connector.GetCmd())
	GetCmd.AddCommand(kafka_cluster.GetCmd())
	GetCmd.AddCommand(kafka_user.GetCmd())
	GetCmd.AddCommand(key_stream.GetCmd())
	GetCmd.AddCommand(schema_code.GetCmd())
	GetCmd.AddCommand(usage.GetCmd())
	GetCmd.AddCommand(installation.GetCmd())
	GetCmd.AddCommand(member.GetCmd())
	GetCmd.AddCommand(data_contract.GetCmd())
	GetCmd.AddCommand(project.GetCmd())
	GetCmd.AddCommand(policy.GetCmd())

	GetCmd.PersistentFlags().BoolP(common.RecursiveFlagName, common.RecursiveFlagShorthand, false, common.RecursiveFlagUsage)
}
