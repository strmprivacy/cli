package context

import (
	"fmt"
	"github.com/spf13/cobra"
	"path"
	"strings"
	"strmprivacy/strm/pkg/auth"
	"strmprivacy/strm/pkg/common"
	"strmprivacy/strm/pkg/entity/project"
)

const (
	configCommandName     = "config"
	entityInfoCommandName = "info"
	accountCommandName    = "account"
	projectCommandName    = "project"
)

func Configuration() *cobra.Command {
	configuration := &cobra.Command{
		Use:               configCommandName,
		Short:             "Shows the config path and preferences",
		DisableAutoGenTag: true,
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
		fmt.Sprintf("configuration output format [%v]", common.ConfigOutputFormatFlagAllowedValues),
	)
	err := configuration.RegisterFlagCompletionFunc(common.OutputFormatFlag, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return common.ConfigOutputFormatFlagAllowedValues, cobra.ShellCompDirectiveNoFileComp
	})

	common.CliExit(err)

	return configuration
}

func Account() *cobra.Command {
	cmd := &cobra.Command{
		Use:               accountCommandName,
		Short:             "Show the handle of this account",
		DisableAutoGenTag: true,
		PersistentPreRun:  auth.RequireAuthenticationPreRun,
		PreRun: func(cmd *cobra.Command, args []string) {
			printer = configurePrinter(cmd)
		},
		Run: func(cmd *cobra.Command, args []string) {
			showAccountDetails()
		},
	}
	cmd.Flags().StringP(
		common.OutputFormatFlag,
		common.OutputFormatFlagShort,
		common.OutputFormatPlain,
		fmt.Sprintf("configuration output format [%v]", common.AccountOutputFormatFlagAllowedValuesText),
	)
	err := cmd.RegisterFlagCompletionFunc(common.OutputFormatFlag, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return common.AccountOutputFormatFlagAllowedValues, cobra.ShellCompDirectiveNoFileComp
	})
	common.CliExit(err)
	return cmd
}

func EntityInfo() *cobra.Command {
	entityInfo := &cobra.Command{
		Use:     entityInfoCommandName + " (entity-reference)",
		Short:   "Show the stored information for a saved entity",
		Example: `strm context info Stream/demo1`,
		Long: `Shows information of entities that have been saved with the --save option
use tab-completion for easy access.

One can see the file-system path, or the entity json contents depending on the output format flag.
Note that these entities do not necessarily still exist on the STRM service. They don't get automatically
removed if for instance you deleted an entity via another client.
`,
		DisableAutoGenTag: true,
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
		fmt.Sprintf("entity information output format [%v]", strings.Join(common.ContextOutputFormatFlagAllowedValues, ", ")),
	)
	err := entityInfo.RegisterFlagCompletionFunc(common.OutputFormatFlag, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return common.ContextOutputFormatFlagAllowedValues, cobra.ShellCompDirectiveNoFileComp
	})

	common.CliExit(err)

	return entityInfo
}

func Project() *cobra.Command {
	cmd := &cobra.Command{
		Use:               projectCommandName + " [name]",
		Short:             "Show or set the active project",
		Args:              cobra.MaximumNArgs(1),
		DisableAutoGenTag: true,
		PersistentPreRun:  auth.RequireAuthenticationPreRun,
		PreRun: func(cmd *cobra.Command, args []string) {
			printer = configurePrinter(cmd)
		},
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) > 0 {
				SetActiveProject(args[0])
			} else {
				showActiveProject()
			}
		},
		ValidArgsFunction: project.NamesCompletion,
	}
	cmd.Flags().StringP(
		common.OutputFormatFlag,
		common.OutputFormatFlagShort,
		common.OutputFormatPlain,
		common.ProjectOutputFormatFlagAllowedValuesText,
	)
	err := cmd.RegisterFlagCompletionFunc(common.OutputFormatFlag, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return common.ProjectOutputFormatFlagAllowedValues, cobra.ShellCompDirectiveNoFileComp
	})

	common.CliExit(err)
	return cmd
}

func savedEntitiesCompletion(cmd *cobra.Command, args []string, complete string) ([]string, cobra.ShellCompDirective) {
	return listSavedEntities(path.Join(common.ConfigPath, common.SavedEntitiesDirectory)), cobra.ShellCompDirectiveNoFileComp
}
