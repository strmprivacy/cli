package monitor

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/strmprivacy/api-definitions-go/v2/api/monitoring/v1"
	"strings"
	"strmprivacy/strm/pkg/common"
	"strmprivacy/strm/pkg/util"
)

var longDoc = util.LongDocsUsage(``)

const followFlag = "follow"
const followFlagWatchAlias = "watch"

func Command(entityType monitoring.EntityState_EntityType) *cobra.Command {
	typeLowercase := "all"
	maxArgs := 0
	short := "monitor all entity types"
	if entityType != 0 {
		typeLowercase = strings.ReplaceAll(strings.ToLower(entityType.String()), "_", "-")
		maxArgs = 1
		short = "monitor entities of type " + typeLowercase
	}

	cmd := &cobra.Command{
		Use:               typeLowercase,
		Short:             short,
		DisableAutoGenTag: true,
		Long:              longDoc,
		PreRun: func(cmd *cobra.Command, args []string) {
			printer = configurePrinter(cmd)
		},
		Run: func(cmd *cobra.Command, args []string) {
			run(cmd, entityType, args)
		},
		Args: cobra.MaximumNArgs(maxArgs), // the optional followFlag of the entity
	}
	log.Infoln("Hello")

	flags := cmd.Flags()
	//flags.StringP(
	//	common.OutputFormatFlag,
	//	common.OutputFormatFlagShort,
	//	common.OutputFormatTable,
	//	fmt.Sprintf("configuration output format [%v]", common.ConfigOutputFormatFlagAllowedValues),
	//)

	//err := cmd.RegisterFlagCompletionFunc(common.OutputFormatFlag, func(command *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	//	//log.Infoln(fmt.Sprintf("Registering flag completion for: %v", cmd.CommandPath()))
	//	//
	//	//follow := util.GetBool(command.Flags(), followFlag)
	//	//log.Traceln(fmt.Sprintf("%v should follow: %v", cmd.CommandPath(), follow))
	//	//var allowedValues []string
	//	//
	//	//if follow {
	//	//	allowedValues = common.MonitorFollowOutputFormatFlagAllowedValues
	//	//} else {
	//	//	allowedValues = common.MonitorOutputFormatFlagAllowedValues
	//	//}
	//	//
	//	//log.Traceln(fmt.Sprintf("%v allowed values: %v", cmd.CommandPath(), allowedValues))
	//
	//	return common.MonitorOutputFormatFlagAllowedValues, cobra.ShellCompDirectiveNoFileComp
	//})
	//common.CliExit(err)

	flags.Bool(followFlag, false, "continuously monitor these events")
	//flags.Bool(followFlagWatchAlias, false, "continuously monitor these events")
	//cmd.Flags().SetNormalizeFunc(normalizeWatch)
	return cmd
}

func normalizeWatch(f *pflag.FlagSet, name string) pflag.NormalizedName {
	switch name {
	case followFlagWatchAlias:
		name = followFlag
	case followFlag:
		follow := util.GetBool(f, followFlag)
		outputFormat := util.GetStringAndErr(f, common.OutputFormatFlag)
		if follow && outputFormat == common.OutputFormatTable {
			err := f.Set(common.OutputFormatFlag, common.OutputFormatJson)
			common.CliExit(err)
		}
		break
	}
	return pflag.NormalizedName(name)
}
