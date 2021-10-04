package context

import (
	"fmt"
	"github.com/spf13/cobra"
	"path"
	"streammachine.io/strm/pkg/common"
)

const (
	configCommandName     = "config"
	entityInfoCommandName = "info"
)

func Configuration() *cobra.Command {
	configuration := &cobra.Command{
		Use:   configCommandName,
		Short: "Shows the config path and preferences",
		PreRun: func(cmd *cobra.Command, args []string) {
			printer = configurePrinter(cmd)
		},
		Run: func(cmd *cobra.Command, args []string) {
			showConfiguration()
		},
	}

	configuration.Flags().StringP(
		common.OutputFormatFlag,
		common.OutputFormatFlagShort,
		common.OutputFormatPlain,
		fmt.Sprintf("Configuration output format [%v]", common.ConfigOutputFormatFlagAllowedValues),
	)
	err := configuration.RegisterFlagCompletionFunc(common.OutputFormatFlag, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return common.ConfigOutputFormatFlagAllowedValues, cobra.ShellCompDirectiveNoFileComp
	})

	common.CliExit(err)

	return configuration
}

func EntityInfo() *cobra.Command {
	entityInfo := &cobra.Command{
		Use:   entityInfoCommandName,
		Short: "Show the stored information for a saved entity",
		PreRun: func(cmd *cobra.Command, args []string) {
			printer = configurePrinter(cmd)
		},
		Run: func(cmd *cobra.Command, args []string) {
			entityInfo(args)
		},
		Args:              cobra.ExactArgs(1), // the contract reference
		ValidArgsFunction: savedEntitiesCompletion,
	}

	entityInfo.Flags().StringP(
		common.OutputFormatFlag,
		common.OutputFormatFlagShort,
		common.OutputFormatFilepath,
		fmt.Sprintf("Entity information output format [%v]", common.ContextOutputFormatFlagAllowedValues),
	)
	err := entityInfo.RegisterFlagCompletionFunc(common.OutputFormatFlag, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return common.ContextOutputFormatFlagAllowedValues, cobra.ShellCompDirectiveNoFileComp
	})

	common.CliExit(err)

	return entityInfo
}

func savedEntitiesCompletion(cmd *cobra.Command, args []string, complete string) ([]string, cobra.ShellCompDirective) {
	return listSavedEntities(path.Join(common.ConfigPath, common.SavedEntitiesDirectory)), cobra.ShellCompDirectiveNoFileComp
}
