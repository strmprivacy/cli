package schema_code

import (
	"github.com/spf13/cobra"
	"streammachine.io/strm/pkg/entity/schema"
)

func GetCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "schema-code (schema-ref)",
		Short: "Get schema code archive by schema-ref",
		Run: func(cmd *cobra.Command, args []string) {
			get(&args[0])
		},
		Args: cobra.ExactArgs(1), // the schema ref
		ValidArgsFunction: schema.NamesCompletion,
	}
}