package monitor

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/strmprivacy/api-definitions-go/v2/api/monitoring/v1"
	"strmprivacy/strm/pkg/common"
	"strmprivacy/strm/pkg/entity/batch_exporter"
	"strmprivacy/strm/pkg/entity/batch_job"
	"strmprivacy/strm/pkg/entity/stream"
	"strmprivacy/strm/pkg/util"
)

var longDoc = util.LongDocsUsage(`
View states for an entity. Can be used to determine why certain entities are not behaving as expected.
`)

const followFlag = "follow"
const followFlagWatchAlias = "watch"

func Command(entityType monitoring.EntityState_EntityType) *cobra.Command {
	normalizedType := "all"
	maxArgs := 0
	short := "monitor all entity types"
	var aliases []string

	if entityType != 0 {
		normalizedType = util.NormalizeEntityStateTypeName(entityType)
		aliases = []string{fmt.Sprintf("%ss", normalizedType)}
		maxArgs = 1
		short = "monitor entities of type " + normalizedType
	}

	cmd := &cobra.Command{
		Use:               normalizedType,
		Short:             short,
		DisableAutoGenTag: true,
		Long:              longDoc,
		Aliases:           aliases,
		PreRun: func(cmd *cobra.Command, args []string) {
			setDefaultOutputFormat(cmd)
			printer = configurePrinter(cmd)
		},
		Run: func(cmd *cobra.Command, args []string) {
			run(cmd, entityType, args)
		},
		Args:              cobra.MaximumNArgs(maxArgs), // the optional followFlag of the entity
		ValidArgsFunction: namesCompletion,
	}

	flags := cmd.Flags()
	flags.StringP(
		common.OutputFormatFlag,
		common.OutputFormatFlagShort,
		common.OutputFormatTable,
		fmt.Sprintf("monitor output format, follow specified=[%v], default=[%v]", common.MonitorFollowOutputFormatFlagAllowedValuesText, common.MonitorOutputFormatFlagAllowedValuesText),
	)

	err := cmd.RegisterFlagCompletionFunc(common.OutputFormatFlag, func(command *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		log.Traceln(fmt.Sprintf("Registering flag completion for: %v", cmd.CommandPath()))

		follow := util.GetBool(command.Flags(), followFlag)
		log.Traceln(fmt.Sprintf("%v should follow: %v", cmd.CommandPath(), follow))
		var allowedValues []string

		if follow {
			allowedValues = common.MonitorFollowOutputFormatFlagAllowedValues
		} else {
			allowedValues = common.MonitorOutputFormatFlagAllowedValues
		}

		log.Traceln(fmt.Sprintf("%v allowed values: %v", cmd.CommandPath(), allowedValues))

		return allowedValues, cobra.ShellCompDirectiveNoFileComp
	})
	common.CliExit(err)

	flags.Bool(followFlag, false, "continuously monitor these events")
	flags.Bool(followFlagWatchAlias, false, "continuously monitor these events")
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

func setDefaultOutputFormat(cmd *cobra.Command) {
	f := cmd.Flags()
	follow := util.GetBool(f, followFlag)
	log.Traceln(fmt.Sprintf("Monitor normalize flags, follow = %t", follow))
	outputFormat := util.GetStringAndErr(f, common.OutputFormatFlag)
	if follow && outputFormat == common.OutputFormatTable {
		err := f.Set(common.OutputFormatFlag, common.OutputFormatJson)
		common.CliExit(err)
	}
}
