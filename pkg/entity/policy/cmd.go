package policy

import (
	"fmt"
	"github.com/spf13/cobra"
	"strings"
	"strmprivacy/strm/pkg/common"
	"strmprivacy/strm/pkg/util"
)

func ListCmd() *cobra.Command {
	outputFormatFlagAllowedValues := []string{common.OutputFormatPlain, common.OutputFormatTable,
		common.OutputFormatJson, common.OutputFormatJsonRaw}
	outputFormatFlagAllowedValuesText := strings.Join(outputFormatFlagAllowedValues, ", ")
	command := &cobra.Command{
		Use:               "policies",
		Short:             "list all policies belonging to this organization",
		Example:           listExample,
		Long:              `List the Policies owned by this organization`,
		DisableAutoGenTag: true,
		PreRun: func(cmd *cobra.Command, args []string) {
			printer = configurePrinter(cmd)
		},
		Run: func(cmd *cobra.Command, args []string) {
			list()
		},
	}
	flags := command.Flags()
	flags.StringP(common.OutputFormatFlag, common.OutputFormatFlagShort, common.OutputFormatTable,
		fmt.Sprintf("output format [%v]", outputFormatFlagAllowedValuesText))
	return command
}

func GetCmd() *cobra.Command {
	outputFormatFlagAllowedValues := []string{common.OutputFormatPlain,
		common.OutputFormatJson, common.OutputFormatJsonRaw}
	outputFormatFlagAllowedValuesText := strings.Join(outputFormatFlagAllowedValues, ", ")
	command := &cobra.Command{
		Use:               "policy [name]",
		Short:             "Get Policy by name or id",
		Long:              "Get a Policy by name or by --id=policy-id",
		DisableAutoGenTag: true,
		Example:           policyExample,
		PreRun: func(cmd *cobra.Command, args []string) {
			printer = configurePrinter(cmd)
		},
		Run: func(cmd *cobra.Command, args []string) {
			get(cmd, args)
		},
		Args:              cobra.MaximumNArgs(1), // the policy name
		ValidArgsFunction: namesCompletion,
	}
	flags := command.Flags()
	flags.StringP(common.OutputFormatFlag, common.OutputFormatFlagShort, common.OutputFormatPlain,
		fmt.Sprintf("output format [%v]", outputFormatFlagAllowedValuesText))
	flags.String(idFlag, "", "policy id")
	flags.Bool(defaultPolicyFlag, false, "get the no-name/no-id default policy")
	_ = command.RegisterFlagCompletionFunc(idFlag, idsCompletion)
	return command
}

func DeleteCmd() *cobra.Command {
	outputFormatFlagAllowedValues := []string{common.OutputFormatPlain,
		common.OutputFormatJson, common.OutputFormatJsonRaw}
	outputFormatFlagAllowedValuesText := strings.Join(outputFormatFlagAllowedValues, ", ")
	command := &cobra.Command{
		Use:   "policy [name]",
		Short: "Delete Policy by name or id",
		Long: util.LongDocsUsage(`
			Delete a Policy by name or id.
			Policies can only be deleted in DRAFT state.
		`),
		Example:           `strm delete policy "1 year" or strm delete policy --id=34c4709e-b8bc-4b45-aa5a-883f471869e3`,
		DisableAutoGenTag: true,
		PreRun: func(cmd *cobra.Command, args []string) {
			printer = configurePrinter(cmd)
		},
		Run: func(cmd *cobra.Command, args []string) {
			del(cmd, args)
		},
		Args:              cobra.MaximumNArgs(1), // the policy name
		ValidArgsFunction: namesCompletion,
	}
	flags := command.Flags()
	flags.StringP(common.OutputFormatFlag, common.OutputFormatFlagShort, common.OutputFormatPlain,
		fmt.Sprintf("output format [%v]", outputFormatFlagAllowedValuesText))
	flags.String(idFlag, "", "policy id")
	_ = command.RegisterFlagCompletionFunc(idFlag, idsCompletion)
	return command
}

func ActivateCmd() *cobra.Command {
	outputFormatFlagAllowedValues := []string{common.OutputFormatPlain,
		common.OutputFormatJson, common.OutputFormatJsonRaw}
	outputFormatFlagAllowedValuesText := strings.Join(outputFormatFlagAllowedValues, ", ")
	command := &cobra.Command{
		Use:     "policy [policy]",
		Short:   "Set the state of a Policy to ACTIVATED",
		Example: `strm activate policy "1 year" or strm activate policy --id=34c4709e-b8bc-4b45-aa5a-883f471869e3`,
		PreRun: func(cmd *cobra.Command, args []string) {
			printer = configurePrinter(cmd)
		},
		Run: func(cmd *cobra.Command, args []string) {
			activate(cmd, args)
		},
		Args:              cobra.MaximumNArgs(1), // the policy name
		DisableAutoGenTag: true,
		ValidArgsFunction: namesCompletion,
	}
	flags := command.Flags()
	flags.StringP(common.OutputFormatFlag, common.OutputFormatFlagShort, common.OutputFormatPlain,
		fmt.Sprintf("output format [%v]", outputFormatFlagAllowedValuesText))
	flags.StringP(idFlag, "", "", "policy id")
	_ = command.RegisterFlagCompletionFunc(idFlag, idsCompletion)
	return command
}

func ArchiveCmd() *cobra.Command {
	outputFormatFlagAllowedValues := []string{common.OutputFormatPlain,
		common.OutputFormatJson, common.OutputFormatJsonRaw}
	outputFormatFlagAllowedValuesText := strings.Join(outputFormatFlagAllowedValues, ", ")
	command := &cobra.Command{
		Use:     "policy [policy]",
		Short:   "Set the state of a Policy to ARCHIVED",
		Example: `strm archive policy "1 year" or strm archive policy --id=34c4709e-b8bc-4b45-aa5a-883f471869e3`,
		PreRun: func(cmd *cobra.Command, args []string) {
			printer = configurePrinter(cmd)
		},
		Run: func(cmd *cobra.Command, args []string) {
			archive(cmd, args)
		},
		Args:              cobra.MaximumNArgs(1), // the policy name
		DisableAutoGenTag: true,
		ValidArgsFunction: namesCompletion,
	}
	flags := command.Flags()
	flags.StringP(common.OutputFormatFlag, common.OutputFormatFlagShort, common.OutputFormatPlain,
		fmt.Sprintf("output format [%v]", outputFormatFlagAllowedValuesText))
	flags.StringP(idFlag, "", "", "policy id")
	_ = command.RegisterFlagCompletionFunc(idFlag, idsCompletion)

	return command
}

func CreateCmd() *cobra.Command {
	outputFormatFlagAllowedValues := []string{common.OutputFormatPlain,
		common.OutputFormatJson, common.OutputFormatJsonRaw}
	outputFormatFlagAllowedValuesText := strings.Join(outputFormatFlagAllowedValues, ", ")
	command := &cobra.Command{
		Use:               "policy",
		Short:             "Create a Policy",
		Long:              longCreateDoc,
		Example:           `strm create policy --name="1 year" --retention 365 --description "1 year for marketing"`,
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
	flags.String(descriptionFlag, "", "description of the policy")
	flags.String(legalGroundsFlag, "", "legal grounds of this policy")
	flags.Int32(retentionFlag, 365, "retention in days of this policy")
	_ = command.MarkFlagRequired(nameFlag)
	_ = command.MarkFlagRequired(retentionFlag)
	return command
}

func UpdateCmd() *cobra.Command {
	outputFormatFlagAllowedValues := []string{common.OutputFormatPlain,
		common.OutputFormatJson, common.OutputFormatJsonRaw}
	outputFormatFlagAllowedValuesText := strings.Join(outputFormatFlagAllowedValues, ", ")
	command := &cobra.Command{
		Use:     "policy policy-id",
		Short:   "Update a Policy",
		Example: `strm update policy 34c4709e-b8bc-4b45-aa5a-883f471869e3 --legal-grounds "EU law x.y.z"`,
		Long: util.LongDocsUsage(`
		Update the attributes of a Policy

		Policies can only be updated while in draft state!
		The policy to be updated must be referenced by its id.
		You can change all other attributes of a policy.

		In order to make a policy active for pipeline processing, you must first 'activate' it.
		`),
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
	flags.String(descriptionFlag, "", "description of the policy")
	flags.String(legalGroundsFlag, "", "legal grounds of this policy")
	flags.Int32(retentionFlag, 365, "retention in days of this policy")
	return command
}
