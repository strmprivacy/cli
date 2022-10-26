package policy

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/strmprivacy/api-definitions-go/v2/api/entities/v1"
	"strings"
	"strmprivacy/strm/pkg/common"
)

func ListCmd() *cobra.Command {
	longDoc := `Query Prime for a list of policies owned by this organization
`
	outputFormatFlagAllowedValues := []string{common.OutputFormatPlain,
		common.OutputFormatJson, common.OutputFormatJsonRaw}
	outputFormatFlagAllowedValuesText := strings.Join(outputFormatFlagAllowedValues, ", ")
	command := &cobra.Command{
		Use:               "policies",
		Short:             "List all policies belonging to this organization",
		Long:              longDoc,
		DisableAutoGenTag: true,
		PreRun: func(cmd *cobra.Command, args []string) {
			printer = configurePrinter(cmd)
		},
		Run: func(cmd *cobra.Command, args []string) {
			list()
		},
	}
	flags := command.Flags()
	flags.StringP(common.OutputFormatFlag, common.OutputFormatFlagShort, common.OutputFormatPlain,
		fmt.Sprintf("output format [%v]", outputFormatFlagAllowedValuesText))
	return command
}

func GetCmd() *cobra.Command {
	outputFormatFlagAllowedValues := []string{common.OutputFormatPlain,
		common.OutputFormatJson, common.OutputFormatJsonRaw}
	outputFormatFlagAllowedValuesText := strings.Join(outputFormatFlagAllowedValues, ", ")
	command := &cobra.Command{
		Use:               "policy (name)",
		Short:             "Get Policy by name",
		DisableAutoGenTag: true,
		PreRun: func(cmd *cobra.Command, args []string) {
			printer = configurePrinter(cmd)
		},
		Run: func(cmd *cobra.Command, args []string) {
			get(cmd, args)
		},
		Args:              cobra.MaximumNArgs(1), // the contract reference
		ValidArgsFunction: namesCompletion,
	}
	flags := command.Flags()
	flags.StringP(common.OutputFormatFlag, common.OutputFormatFlagShort, common.OutputFormatPlain,
		fmt.Sprintf("output format [%v]", outputFormatFlagAllowedValuesText))
	flags.String(idFlag, "", "policy id")
	_ = command.RegisterFlagCompletionFunc(idFlag, idsCompletion)
	return command
}

func DeleteCmd() *cobra.Command {
	outputFormatFlagAllowedValues := []string{common.OutputFormatPlain,
		common.OutputFormatJson, common.OutputFormatJsonRaw}
	outputFormatFlagAllowedValuesText := strings.Join(outputFormatFlagAllowedValues, ", ")
	command := &cobra.Command{
		Use:               "policy (name)",
		Short:             "Delete Policy by name",
		DisableAutoGenTag: true,
		PreRun: func(cmd *cobra.Command, args []string) {
			printer = configurePrinter(cmd)
		},
		Run: func(cmd *cobra.Command, args []string) {
			del(cmd, args)
		},
		Args:              cobra.MaximumNArgs(1), // the contract reference
		ValidArgsFunction: namesCompletion,
	}
	flags := command.Flags()
	flags.StringP(common.OutputFormatFlag, common.OutputFormatFlagShort, common.OutputFormatPlain,
		fmt.Sprintf("output format [%v]", outputFormatFlagAllowedValuesText))
	flags.String(idFlag, "", "policy id")
	_ = command.RegisterFlagCompletionFunc(idFlag, idsCompletion)
	return command
}

func CreateCmd() *cobra.Command {
	outputFormatFlagAllowedValues := []string{common.OutputFormatPlain,
		common.OutputFormatJson, common.OutputFormatJsonRaw}
	outputFormatFlagAllowedValuesText := strings.Join(outputFormatFlagAllowedValues, ", ")
	command := &cobra.Command{
		Use:               "policy (name)",
		Short:             "Create Policy",
		DisableAutoGenTag: true,
		PreRun: func(cmd *cobra.Command, args []string) {
			printer = configurePrinter(cmd)
		},
		Run: func(cmd *cobra.Command, args []string) {
			create(cmd)
		},
		Args: cobra.NoArgs,
	}
	flags := command.Flags()
	flags.StringP(common.OutputFormatFlag, common.OutputFormatFlagShort, common.OutputFormatPlain,
		fmt.Sprintf("output format [%v]", outputFormatFlagAllowedValuesText))
	flags.String(nameFlag, "", "name")
	flags.String(descriptionFlag, "", "descriptive text")
	flags.String(legalGroundsFlag, "", "legal grounds of this policy")
	flags.Int32(retentionFlag, 365, "retention in days of this policy")
	draft := entities.Policy_State_name[1]
	flags.String(stateFlag, draft, "State of the policy")
	return command
}

func UpdateCmd() *cobra.Command {
	outputFormatFlagAllowedValues := []string{common.OutputFormatPlain,
		common.OutputFormatJson, common.OutputFormatJsonRaw}
	outputFormatFlagAllowedValuesText := strings.Join(outputFormatFlagAllowedValues, ", ")
	command := &cobra.Command{
		Use:               "policy (policy-id)",
		Short:             "Update Policy",
		DisableAutoGenTag: true,
		PreRun: func(cmd *cobra.Command, args []string) {
			printer = configurePrinter(cmd)
		},
		Run: func(cmd *cobra.Command, args []string) {
			update(cmd, args[0])
		},
		Args:              cobra.ExactArgs(1),
		ValidArgsFunction: idsCompletion,
	}
	flags := command.Flags()
	flags.StringP(common.OutputFormatFlag, common.OutputFormatFlagShort, common.OutputFormatPlain,
		fmt.Sprintf("output format [%v]", outputFormatFlagAllowedValuesText))
	flags.String(nameFlag, "", "name")
	flags.String(descriptionFlag, "", "descriptive text")
	flags.String(legalGroundsFlag, "", "legal grounds of this policy")
	flags.Int32(retentionFlag, 365, "retention in days of this policy")
	flags.StringSliceP(updateMaskFlag, "u", []string{}, "list of fields to update")
	_ = command.RegisterFlagCompletionFunc(updateMaskFlag, fieldMaskCompletion)
	flags.String(stateFlag, "", "State of the policy")
	_ = command.RegisterFlagCompletionFunc(stateFlag, stateCompletion)
	_ = command.RegisterFlagCompletionFunc(updateMaskFlag, fieldMaskCompletion)
	return command
}

func stateCompletion(_ *cobra.Command, _ []string, _ string) ([]string, cobra.ShellCompDirective) {
	return []string{"draft", "active", "archived"}, cobra.ShellCompDirectiveNoFileComp
}

func fieldMaskCompletion(_ *cobra.Command, _ []string, _ string) ([]string, cobra.ShellCompDirective) {
	return []string{
		"name", "description", "legal_grounds", "retention_days", "state",
	}, cobra.ShellCompDirectiveNoFileComp
}
