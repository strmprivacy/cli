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
		Example:           "strm create data-contract my-handle/my-contract/1.0.0 --contract-definition my-def.json",
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
	flags.String(contractDefinitionFlag, "",
		`filename of the JSON contract definition, containing the keyField, and possibly field metadata, validations and data subject field. Example JSON definition file:
{
  "keyField": "sessionId",
  "fieldMetadata": [
    {
      "fieldName": "userName",
      "personalDataConfig": {
        "isPii": true,
        "isQuasiId": false,
        "purposeLevel": 1
      }
    },
    {
      "fieldName": "userAgeGroup",
      "personalDataConfig": {
        "isPii": false,
        "isQuasiId": true
      },
      "statisticalDataType": "ORDINAL",
      "ordinalValues": ["child","teenager","adult","senior"],
      "nullHandlingConfig": {
        "type": "DEFAULT_VALUE",
        "defaultValue": "adult"
      }
    }
  ],
  "validations": [
    {
      "field": "referrerUrl",
      "type": "regex",
      "value": "^.*strmprivacy.*$"
    }
  ],
  "dataSubjectField": "userId"
}
`)
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
		Long:  longDoc,
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
