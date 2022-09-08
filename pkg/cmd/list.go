package cmd

import (
	"github.com/spf13/cobra"
	"strmprivacy/strm/pkg/common"
	"strmprivacy/strm/pkg/entity/batch_exporter"
	"strmprivacy/strm/pkg/entity/batch_job"
	"strmprivacy/strm/pkg/entity/data_connector"
	"strmprivacy/strm/pkg/entity/data_contract"
	"strmprivacy/strm/pkg/entity/data_subject"
	"strmprivacy/strm/pkg/entity/event_contract"
	"strmprivacy/strm/pkg/entity/installation"
	"strmprivacy/strm/pkg/entity/kafka_cluster"
	"strmprivacy/strm/pkg/entity/kafka_exporter"
	"strmprivacy/strm/pkg/entity/kafka_user"
	"strmprivacy/strm/pkg/entity/key_stream"
	"strmprivacy/strm/pkg/entity/keylinks"
	"strmprivacy/strm/pkg/entity/member"
	"strmprivacy/strm/pkg/entity/project"
	"strmprivacy/strm/pkg/entity/schema"
	"strmprivacy/strm/pkg/entity/stream"
)

var ListCmd = &cobra.Command{
	Use:               common.ListCommandName,
	DisableAutoGenTag: true,
	Short:             "List entities",
}

func init() {
	ListCmd.AddCommand(stream.ListCmd())
	ListCmd.AddCommand(kafka_exporter.ListCmd())
	ListCmd.AddCommand(batch_exporter.ListCmd())
	ListCmd.AddCommand(batch_job.ListCmd())
	ListCmd.AddCommand(data_connector.ListCmd())
	ListCmd.AddCommand(kafka_cluster.ListCmd())
	ListCmd.AddCommand(kafka_user.ListCmd())
	ListCmd.AddCommand(key_stream.ListCmd())
	ListCmd.AddCommand(schema.ListCmd())
	ListCmd.AddCommand(event_contract.ListCmd())
	ListCmd.AddCommand(installation.ListCmd())
	ListCmd.AddCommand(project.ListCmd())
	ListCmd.AddCommand(member.ListCmd())
	ListCmd.AddCommand(data_subject.ListCmd())
	ListCmd.AddCommand(keylinks.ListCmd())
	ListCmd.AddCommand(data_contract.ListCmd())

	ListCmd.PersistentFlags().BoolP(common.RecursiveFlagName, common.RecursiveFlagShorthand, false, common.RecursiveFlagUsage)
}
