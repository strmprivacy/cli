package data_contract

import (
	"github.com/spf13/cobra"
	"strmprivacy/strm/pkg/common"
)

var longDoc = `### Usage`

func CreateCmd() *cobra.Command {
	dataContract := &cobra.Command{
		Use:               "data-contract (handle/name/version)",
		Short:             "create a data contract",
		Long:              longDoc,
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
	flags.String(schemaDefinitionFlag, "", "filename of the schema definition (yaml or json) - either a Simple Schema, Avro Schema or Json Schema")
	flags.Bool(publicFlag, false, "whether the data contract should be made public (accessible to other STRM Privacy customers)")
	flags.String(projectName, "", `Project name to create resource in`)
	flags.String(contractDefinitionFlag, "",
		`The path to the file with the keyField, and possibly piiFields and validations. Example JSON definition file:
{
    "keyField": "sessionId",
    "piiFields": {
        "sessionId": 2,
        "referrerUrl": 1
    },
    "validations": [
        {
            "field": "referrerUrl",
            "type": "regex",
            "value": "^.*strmprivacy.*$"
        }
    ]
}`)
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
			list()
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
		Long:  longDoc,
		PreRun: func(cmd *cobra.Command, args []string) {
			printer = configurePrinter(cmd)
		},
		Run: func(cmd *cobra.Command, args []string) {
			del(&args[0])
		},
		Args:              cobra.ExactArgs(1), // the contract reference
		DisableAutoGenTag: true,
		ValidArgsFunction: RefsCompletion,
	}

	return dataContract
}
