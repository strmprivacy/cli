package schema_code

import (
	"github.com/spf13/cobra"
	"streammachine.io/strm/pkg/entity/schema"
)

func GetCmd() *cobra.Command {
	getCommand := &cobra.Command{
		Use:   "schema-code (schema-ref)",
		Short: "Get schema code archive by schema-ref",
		Run: func(cmd *cobra.Command, args []string) {
			get(cmd, &args[0])
		},
		Args:              cobra.ExactArgs(1), // the schema ref
		ValidArgsFunction: schema.NamesCompletion,
	}
	flags := getCommand.Flags()
	flags.String(languageFlag, "", "which programming language")
	flags.String(filenameFlag, "", "Destination zip file location")
	flags.Bool(overwriteFlag, false, "do we allow overwriting an existing file")
	_ = getCommand.MarkFlagRequired(languageFlag)
	_ = getCommand.RegisterFlagCompletionFunc(languageFlag, languageCompletion)
	return getCommand
}

func languageCompletion(cmd *cobra.Command, args []string, complete string) ([]string, cobra.ShellCompDirective) {
	return []string{"java", "typescript", "python"}, cobra.ShellCompDirectiveDefault
}
