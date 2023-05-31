package purpose_mapping

import (
	"github.com/spf13/cobra"
	"strconv"
	"strmprivacy/strm/pkg/common"
)

func CreateCmd() *cobra.Command {
	return &cobra.Command{
		Use:               "purpose-mapping (name)",
		Short:             "Create a purpose mapping",
		Long:              longDocCreate,
		Example:           createExample,
		DisableAutoGenTag: true,
		PreRun: func(cmd *cobra.Command, args []string) {
			printer = configurePrinter(cmd)
		},
		Run: func(cmd *cobra.Command, args []string) {
			create(args[0])
		},
		Args:              cobra.ExactArgs(1), // the name of the purpose mapping
		ValidArgsFunction: common.NoFilesEmptyCompletion,
	}
}

func GetCmd() *cobra.Command {
	return &cobra.Command{
		Use:               "purpose-mapping (value)",
		Short:             "Get a purpose mapping by the integer value",
		Example:           getExample,
		DisableAutoGenTag: true,
		PreRun: func(cmd *cobra.Command, args []string) {
			printer = configurePrinter(cmd)
		},
		Run: func(cmd *cobra.Command, args []string) {
			level, err := strconv.Atoi(args[0])
			common.CliExit(err)

			get(int32(level))
		},
		Args:              cobra.ExactArgs(1), // the integer value of the purpose mapping
		ValidArgsFunction: LevelsCompletion,
	}
}
func ListCmd() *cobra.Command {
	return &cobra.Command{
		Use:               "purpose-mappings",
		Short:             "List purpose mappings",
		Example:           listExample,
		DisableAutoGenTag: true,
		PreRun: func(cmd *cobra.Command, args []string) {
			printer = configurePrinter(cmd)
		},
		Run: func(cmd *cobra.Command, args []string) {
			listAndPrint()
		},
		ValidArgsFunction: common.NoFilesEmptyCompletion,
	}
}
