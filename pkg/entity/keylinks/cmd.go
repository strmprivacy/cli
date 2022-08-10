package keylinks

import (
	"github.com/spf13/cobra"
)

var longDoc = `### Usage`

func ListCmd() *cobra.Command {
	command := &cobra.Command{
		Use:               "data-subject-keylinks <data-subject-id>...",
		Short:             "List data subjects keylinks",
		Long:              longDoc,
		DisableAutoGenTag: true,
		PreRun: func(cmd *cobra.Command, args []string) {
			printer = configurePrinter(cmd)
		},
		Run: func(cmd *cobra.Command, args []string) {
			list(args, cmd)
		},
	}
	return command
}