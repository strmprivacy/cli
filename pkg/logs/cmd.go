package logs

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/strmprivacy/api-definitions-go/v2/api/monitoring/v1"
	"strmprivacy/strm/pkg/common"
	"strmprivacy/strm/pkg/entity/batch_exporter"
	"strmprivacy/strm/pkg/entity/batch_job"
	"strmprivacy/strm/pkg/entity/stream"
	"strmprivacy/strm/pkg/util"
)

// TODO long docs
var longDoc = util.LongDocsUsage(``)

const followFlag = "follow"
const followFlagWatchAlias = "watch"

func Command(entityType monitoring.EntityState_EntityType) *cobra.Command {
	normalizedType := util.NormalizeEntityStateTypeName(entityType)
	short := "show logs for entity of type " + normalizedType

	cmd := &cobra.Command{
		Use:               fmt.Sprintf("%v (name)", normalizedType),
		Short:             short,
		DisableAutoGenTag: true,
		Long:              longDoc,
		PreRun: func(cmd *cobra.Command, args []string) {
			printer = configurePrinter(cmd)
		},
		Run: func(cmd *cobra.Command, args []string) {
			run(cmd, entityType, args)
		},
		Args:              cobra.ExactArgs(1),
		ValidArgsFunction: namesCompletion,
	}

	flags := cmd.Flags()
	flags.StringP(
		common.OutputFormatFlag,
		common.OutputFormatFlagShort,
		common.OutputFormatPlain,
		fmt.Sprintf("logs output format [%v]", common.LogsOutputFormatFlagAllowedValuesText),
	)

	err := cmd.RegisterFlagCompletionFunc(common.OutputFormatFlag, func(command *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return common.LogsOutputFormatFlagAllowedValues, cobra.ShellCompDirectiveNoFileComp
	})
	common.CliExit(err)

	flags.Bool(followFlag, false, "continuously show new logs for this entity")
	flags.Bool(followFlagWatchAlias, false, "continuously show new logs for this entity")
	cmd.Flags().SetNormalizeFunc(normalizeWatch)
	return cmd
}

func namesCompletion(cmd *cobra.Command, args []string, complete string) ([]string, cobra.ShellCompDirective) {
	switch cmd.Name() {
	case util.NormalizeEntityStateTypeName(monitoring.EntityState_BATCH_EXPORTER):
		return batch_exporter.NamesCompletion(cmd, args, complete)
	case util.NormalizeEntityStateTypeName(monitoring.EntityState_BATCH_JOB):
		return batch_job.NamesCompletion(cmd, args, complete)
	case util.NormalizeEntityStateTypeName(monitoring.EntityState_STREAM):
		return stream.NamesCompletion(cmd, args, complete)
	}

	return nil, cobra.ShellCompDirectiveNoFileComp
}

func normalizeWatch(_ *pflag.FlagSet, name string) pflag.NormalizedName {
	switch name {
	case followFlagWatchAlias:
		name = followFlag
		break
	}
	return pflag.NormalizedName(name)
}
