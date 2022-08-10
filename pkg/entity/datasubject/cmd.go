package datasubject

import (
	"github.com/spf13/cobra"
)

var longDoc = `### Usage`

func ListCmd() *cobra.Command {
	command := &cobra.Command{
		Use:               "data-subjects",
		Short:             "List data subjects",
		Long:              longDoc,
		DisableAutoGenTag: true,
		PreRun: func(cmd *cobra.Command, args []string) {
			printer = configurePrinter(cmd)
		},
		Run: func(cmd *cobra.Command, args []string) {
			list(cmd)
		},
	}
	flags := command.Flags()
	flags.Int32(pageSizeFlag, 0, "maximum number of items to be returned")
	flags.String(pageTokenFlag, "", "page token to be entered for next page")
	return command
}