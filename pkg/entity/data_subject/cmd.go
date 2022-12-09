package data_subject

import (
	"fmt"
	"github.com/spf13/cobra"
	"strings"
	"strmprivacy/strm/pkg/common"
	"strmprivacy/strm/pkg/util"
)

func ListCmd() *cobra.Command {
	longDoc := util.LongDocsUsage(`
Query the Data Subjects service for a list of data-subjects.

Returns paginated data. If one page of data has following pages, its °next_page_token°
field must be added to the following call via the °page-token° flag.
`)
	outputFormatFlagAllowedValues := []string{common.OutputFormatPlain, common.OutputFormatPlain + "0",
		common.OutputFormatJson, common.OutputFormatJsonRaw}
	outputFormatFlagAllowedValuesText := strings.Join(outputFormatFlagAllowedValues, ", ")
	command := &cobra.Command{
		Use:               "data-subjects",
		Short:             "List a page of data subjects",
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
	flags.String(pageTokenFlag, "", `page token to be entered for next page.
Use the nextPageToken (if any) returned from the previous call`)
	flags.StringP(common.OutputFormatFlag, common.OutputFormatFlagShort, common.OutputFormatPlain,
		fmt.Sprintf("output format [%v]", outputFormatFlagAllowedValuesText))
	return command
}

func DeleteCmd() *cobra.Command {
	longDoc := util.LongDocsUsage(`
		Deletes 1 or more data subjects from the Data Subjects Service.
		Returns the number of deleted key links that were associated with these data subjects.
`)

	outputFormatFlagAllowedValues := []string{common.OutputFormatPlain,
		common.OutputFormatJson, common.OutputFormatJsonRaw}
	outputFormatFlagAllowedValuesText := strings.Join(outputFormatFlagAllowedValues, ", ")
	command := &cobra.Command{
		Use:               "data-subjects (data-subject-id...)",
		Short:             "Delete Data Subjects",
		Long:              longDoc,
		DisableAutoGenTag: true,
		PreRun: func(cmd *cobra.Command, args []string) {
			printer = configurePrinter(cmd)
		},
		Run: func(cmd *cobra.Command, args []string) {
			del(args, cmd)
		},
	}
	flags := command.Flags()
	flags.StringP(common.OutputFormatFlag, common.OutputFormatFlagShort, common.OutputFormatPlain,
		fmt.Sprintf("output format [%v]", outputFormatFlagAllowedValuesText))
	return command
}
