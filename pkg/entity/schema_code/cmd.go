package schema_code

import (
	"strings"
	"strmprivacy/strm/pkg/entity/data_contract"
	"strmprivacy/strm/pkg/util"

	"github.com/spf13/cobra"
)

var longDoc = `In order to simplify sending correctly serialized data to STRM Privacy it is recommended to use generated
source code that defines a class/object structure in a certain programming language, that knows (with help of some
open-source libraries) how to serialize objects.

The result of a ` + "`get schema-code`" + ` is a zip file with some source code files for getting started with sending events in a
certain programming language. Generally this will be code where you’ll have to do some sort of ` + "`build`" + ` step in order to
make this fully operational in your development setting (using a JDK, a Python or a Node.js environment).

### Usage`

var example = util.DedentTrim(`
strm get schema-code strmprivacy/example/1.3.0 --language=python
Saved to python-avro-example-1.3.0.zip
`)
var languages = []string{"java", "typescript", "python", "rust"}
var languagesString = strings.Join(languages, ", ")

func GetCmd() *cobra.Command {
	getCommand := &cobra.Command{
		Use:               "schema-code (data-contract-ref)",
		Short:             "Get schema code archive by data-contract-ref",
		Long:              longDoc,
		Example:           example,
		DisableAutoGenTag: true,
		Run: func(cmd *cobra.Command, args []string) {
			get(cmd, &args[0])
		},
		Args:              cobra.ExactArgs(1), // the schema ref
		ValidArgsFunction: data_contract.RefsCompletion,
	}
	flags := getCommand.Flags()
	flags.String(languageFlag, "", "which programming language. Supported are: "+languagesString)
	flags.String(filenameFlag, "", "destination zip file location")
	flags.Bool(overwriteFlag, false, "do we allow overwriting an existing file")
	_ = getCommand.MarkFlagRequired(languageFlag)
	_ = getCommand.RegisterFlagCompletionFunc(languageFlag, languageCompletion)
	return getCommand
}

func languageCompletion(cmd *cobra.Command, args []string, complete string) ([]string, cobra.ShellCompDirective) {
	return languages, cobra.ShellCompDirectiveDefault
}
