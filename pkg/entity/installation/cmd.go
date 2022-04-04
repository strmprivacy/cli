package installation

import (
	"github.com/spf13/cobra"
)

var longDoc = `
### Usage
`

func GetCmd() *cobra.Command {
	return &cobra.Command{
		Use:               "installation",
		Short:             "Get your installation",
		Long:              longDoc,
		DisableAutoGenTag: true,
		PreRun: func(cmd *cobra.Command, args []string) {
			printer = configurePrinter(cmd)
		},
		Run: func(cmd *cobra.Command, args []string) {
			get(cmd)
		},
		Args:              cobra.ExactArgs(0),
		ValidArgsFunction: namesCompletion,
	}
}
