package datasubject

import (
	"fmt"
	"github.com/spf13/cobra"
	"strings"
	"strmprivacy/strm/pkg/common"
)

var longDoc = `### Usage`

var outputFormatFlagAllowedValues = []string{common.OutputFormatPlain, common.OutputFormatPlain + "0",
	common.OutputFormatJson, common.OutputFormatJsonRaw}
var outputFormatFlagAllowedValuesText = strings.Join(outputFormatFlagAllowedValues, ", ")

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
	flags.StringP(common.OutputFormatFlag, common.OutputFormatFlagShort, common.OutputFormatPlain,
		fmt.Sprintf("output format [%v]", outputFormatFlagAllowedValuesText))
	return command
}