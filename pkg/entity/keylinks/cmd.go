package keylinks

import (
	"github.com/spf13/cobra"
)

func ListCmd() *cobra.Command {
	longDoc := `Retrieve keylinks associated with data subjects`
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
