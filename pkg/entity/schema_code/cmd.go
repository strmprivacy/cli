package schema_code

import (
	"github.com/spf13/cobra"
	"strmprivacy/strm/pkg/entity/schema"
)

var longDoc = `In order to simplify sending correctly serialized data to STRM Privacy it is recommended to use generated source code
that defines a class/object structure in a certain programming language, that knows
(with help of some open-source libraries) how to serialize objects.

The result of a ` + "`get schema-code`" + ` is a zip file with some source code files for getting started with sending events in a
certain programming language. Generally this will be code where you’ll have to do some sort of ` + "`build`" + ` step in order to
make this fully operational in your development setting (using a JDK, a Python or a Node.js environment).

### Usage`

func GetCmd() *cobra.Command {
	getCommand := &cobra.Command{
		Use:               "schema-code (schema-ref)",
		Short:             "Get schema code archive by schema-ref",
		Long:              longDoc,
		DisableAutoGenTag: true,
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
