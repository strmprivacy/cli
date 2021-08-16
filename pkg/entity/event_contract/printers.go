package event_contract

import (
	"github.com/spf13/cobra"
	"streammachine.io/strm/pkg/util"
)

var printer util.Printer

func configurePrinter(cmd *cobra.Command) util.Printer {
	outputFormat := util.GetStringAndErr(cmd.Flags(), util.OutputFormatFlag)

	switch outputFormat {
	case "json":
		return util.GenericJsonPrinter{}
	case "table":
		return util.GenericJsonPrinter{}
	default:
		return util.GenericJsonPrinter{}
	}
}
