package usage

import (
	"fmt"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/spf13/cobra"
	"github.com/streammachineio/api-definitions-go/api/usage/v1"
	"google.golang.org/protobuf/proto"
	"math"
	"streammachine.io/strm/pkg/common"
	"streammachine.io/strm/pkg/util"
)

var printer util.Printer

func configurePrinter(cmd *cobra.Command) util.Printer {
	outputFormat := util.GetStringAndErr(cmd.Flags(), util.OutputFormatFlag)

	switch outputFormat {
	case "json-raw":
		return util.GenericRawJsonPrinter{}
	case "json":
		return util.GenericPrettyJsonPrinter{}
	case "table":
		common.CliExit("Output format 'table' is not supported for usage.")
		return nil
	case "csv":
		return getCsvPrinter{}
	default:
		return util.GenericPrettyJsonPrinter{}
	}
}

type getCsvPrinter struct{}

func (p getCsvPrinter) Print(data proto.Message) {
	streamUsage, _ := (data).(*usage.GetStreamEventUsageResponse)

	rows := make([]table.Row, 0, len(streamUsage.Windows))

	windowCount := int64(-1)
	change := math.NaN()
	for _, window := range streamUsage.Windows {
		if windowCount != -1 {
			change = float64(window.EventCount - windowCount)
		}
		windowCount = window.EventCount

		windowDuration := window.EndTime.AsTime().Sub(window.StartTime.AsTime())
		rate := change / windowDuration.Seconds()

		rows = append(rows, table.Row{
			fmt.Sprintf("%d", windowCount),
			fmt.Sprintf("%.0f", windowDuration.Seconds()),
			fmt.Sprintf("%v", change),
			fmt.Sprintf("%.2f", rate),
		})
	}

	util.RenderCsv(
		table.Row{
			"From",
			"Count",
			"Duration",
			"Change",
			"Rate",
		},
		rows,
	)
}
