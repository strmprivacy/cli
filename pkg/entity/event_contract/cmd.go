package event_contract

import (
	"github.com/spf13/cobra"
	"strmprivacy/strm/pkg/common"
	"strmprivacy/strm/pkg/entity/schema"
)

const (
	schemaRefFlag  = "schema-ref"
	isPublicFlag   = "public"
	definitionFile = "definition-file"
)

var longDoc = `An Event Contract defines the rules that are to be applied to events.

The Event Contract defines:

- the Schema to use via a full Schema reference (handle/name/version)

- the key field

- the PII fields

- any validations on fields (e.g. a regex to validate an email address)

Like Schemas, Event Contracts can be private or public, allowing them to be found and used by others than the owning
client. Be careful, public Event Contracts cannot be deleted.

Also like Schemas, Event Contracts are versioned using a versioning scheme that can be fully determined by the client.
The only restrictions are that version numbers:

- MUST follow the semantic version format (major/minor/patch),

- MUST always be ascending

An Event Contract is uniquely identified by its Event Contract reference, in the format (organization handle/event
contract name/version).

### Usage
`

func GetCmd() *cobra.Command {
	return &cobra.Command{
		Use:               "event-contract (reference)",
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

func ListCmd() *cobra.Command {
	return &cobra.Command{
		Use:               "event-contracts",
		Short:             "List Event Contracts",
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

func CreateCmd() *cobra.Command {
	contract := &cobra.Command{
		Use:               "event-contract (handle/name/version)",
		Short:             "Create an event-contract with reference 'handle/name/version'",
		Long:              longDoc,
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
