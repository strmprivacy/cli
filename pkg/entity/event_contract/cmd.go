package event_contract

import (
	"github.com/spf13/cobra"
	"streammachine.io/strm/pkg/common"
	"streammachine.io/strm/pkg/entity/schema"
)

const (
	schemaRefFlag  = "schema-ref"
	isPublicFlag   = "public"
	definitionFile = "definition-file"
)

func GetCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "event-contract (reference)",
		Short: "Get Event Contract by reference",
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
		Use:   "event-contracts",
		Short: "List Event Contracts",
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
		Use:   "event-contract (reference)",
		Short: "Create an event-contract",
		Long:  `Create an event contract from a JSON definition file`,
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
            "value": "^.*streammachine.*$"
        }
    ]
}`)

	common.MarkRequiredFlags(contract, schemaRefFlag, definitionFile)
	_ = contract.RegisterFlagCompletionFunc(schemaRefFlag, schema.NamesCompletion)

	return contract
}
