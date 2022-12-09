package data_contract

import (
	"github.com/lithammer/dedent"
	"github.com/spf13/cobra"
	"strings"
	"strmprivacy/strm/pkg/common"
	"strmprivacy/strm/pkg/util"
)

var longDoc = util.LongDocsUsage(`
Data Contracts are the core of STRM Privacy.
See [here](https://docs.strmprivacy.io/docs/latest/concepts/data-contracts/) for details
`)

func CreateCmd() *cobra.Command {
	dataContract := &cobra.Command{
		Use:               "data-contract (handle/name/version)",
		Short:             "Create a Data Contract",
		Long:              longDoc,
		Example:           createExample,
		DisableAutoGenTag: true,
		PreRun: func(cmd *cobra.Command, args []string) {
			printer = configurePrinter(cmd)
		},
		Run: func(cmd *cobra.Command, args []string) {
			create(cmd, &args[0])
		},
		Args: cobra.ExactArgs(1),
	}
	flags := dataContract.Flags()
	flags.String(schemaDefinitionFlag, "", "filename of the schema definition (yaml or JSON) - either a Simple Schema, Avro Schema or JSON Schema")
	flags.Bool(publicFlag, false, "whether the data contract should be made public (accessible to other STRM Privacy customers)")
	flags.String(contractDefinitionFlag, "", dedent.Dedent(strings.TrimSpace(`the path to the file with the keyField, and possibly piiFields and validations. See example.`)))
	common.MarkRequiredFlags(dataContract, schemaDefinitionFlag, contractDefinitionFlag)
	return dataContract
}

func ListCmd() *cobra.Command {
	return &cobra.Command{
		Use:               "data-contracts",
		Short:             "List Data Contracts",
		Long:              longDoc,
		DisableAutoGenTag: true,
		PreRun: func(cmd *cobra.Command, args []string) {
			printer = configurePrinter(cmd)
		},
		Run: func(cmd *cobra.Command, args []string) {
			list(cmd)
		},
	}
}

func GetCmd() *cobra.Command {
	return &cobra.Command{
		Use:               "data-contract (handle/name/version)",
		Short:             "Get Event Contract by reference",
		Long:              longDoc,
		DisableAutoGenTag: true,
		PreRun: func(cmd *cobra.Command, args []string) {
			printer = configurePrinter(cmd)
		},
		Run: func(cmd *cobra.Command, args []string) {
			get(&args[0])
		},
		Args:              cobra.ExactArgs(1), // the contract reference
		ValidArgsFunction: RefsCompletion,
	}
}

func ActivateCmd() *cobra.Command {
	dataContract := &cobra.Command{
		Use:   "data-contract (handle/name/version)",
		Short: "Set the state of a Data Contract to ACTIVATED",
		Long:  longDoc,
		Run: func(cmd *cobra.Command, args []string) {
			activate(&args[0])
		},
		Args:              cobra.ExactArgs(1), // the contract reference
		DisableAutoGenTag: true,
		ValidArgsFunction: RefsCompletion,
	}

	return dataContract
}

func ArchiveCmd() *cobra.Command {
	dataContract := &cobra.Command{
		Use:   "data-contract (handle/name/version)",
		Short: "Set the state of an Data Contract to ARCHIVED",
		Long:  longDoc,
		Run: func(cmd *cobra.Command, args []string) {
			archive(&args[0])
		},
		Args:              cobra.ExactArgs(1), // the contract reference
		DisableAutoGenTag: true,
		ValidArgsFunction: RefsCompletion,
	}

	return dataContract
}

func DeleteCmd() *cobra.Command {
	dataContract := &cobra.Command{
		Use:   "data-contract (handle/name/version)",
		Short: "Delete Data Contract by reference",
		PreRun: func(cmd *cobra.Command, args []string) {
			printer = configurePrinter(cmd)
		},
		Run: func(cmd *cobra.Command, args []string) {
			del(&args[0], cmd)
		},
		Args:              cobra.ExactArgs(1), // the contract reference
		DisableAutoGenTag: true,
		ValidArgsFunction: RefsCompletion,
	}

	return dataContract
}
