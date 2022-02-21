package event_contract

import (
	"github.com/spf13/cobra"
	"io/ioutil"
	"strmprivacy/strm/pkg/common"
	"strmprivacy/strm/pkg/entity/schema"
)

const (
	schemaRefFlag  = "schema-ref"
	isPublicFlag   = "public"
	definitionFile = "definition-file"
)

var content, _ = ioutil.ReadFile("pkg/entity/batch_exporter/docstring.md")

func GetCmd() *cobra.Command {
	return &cobra.Command{
		Use:               "event-contract (reference)",
		Short:             "Get Event Contract by reference",
		Long:              string(content),
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

func ListCmd() *cobra.Command {
	return &cobra.Command{
		Use:               "event-contracts",
		Short:             "List Event Contracts",
		Long:              string(content),
		DisableAutoGenTag: true,
		PreRun: func(cmd *cobra.Command, args []string) {
			printer = configurePrinter(cmd)
		},
		Run: func(cmd *cobra.Command, args []string) {
			list()
		},
	}
}

func CreateCmd() *cobra.Command {
	contract := &cobra.Command{
		Use:               "event-contract (handle/name/version)",
		Short:             "Create an event-contract with reference 'handle/name/version'",
		Long:              string(content),
		DisableAutoGenTag: true,
		PreRun: func(cmd *cobra.Command, args []string) {
			printer = configurePrinter(cmd)
		},
		Run: func(cmd *cobra.Command, args []string) {
			create(cmd, &args[0])
		},
		Args: cobra.ExactArgs(1), // the contract reference
	}

	flags := contract.Flags()
	flags.StringP(schemaRefFlag, "S", "", "The Serialization Schema to which this Event Contract is linked")
	flags.BoolP(isPublicFlag, "P", false, "Public visibility of the Event Contract (allow others to use this contract)")
	flags.StringP(definitionFile, "F", "",
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

	common.MarkRequiredFlags(contract, schemaRefFlag, definitionFile)
	_ = contract.RegisterFlagCompletionFunc(schemaRefFlag, schema.NamesCompletion)

	return contract
}
