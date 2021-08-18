package usage

import (
	"fmt"
	"github.com/bykof/gostradamus"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/spf13/cobra"
	"github.com/streammachineio/api-definitions-go/api/usage/v1"
	"google.golang.org/protobuf/proto"
	"math"
	"streammachine.io/strm/pkg/common"
	"streammachine.io/strm/pkg/constants"
	"streammachine.io/strm/pkg/util"
	"time"
)

var printer util.Printer

func configurePrinter(cmd *cobra.Command) util.Printer {
	outputFormat := util.GetStringAndErr(cmd.Flags(), util.OutputFormatFlag)

	switch outputFormat {
	case constants.OutputFormatJsonRaw:
		return util.GenericRawJsonPrinter{}
	case constants.OutputFormatJson:
		return util.GenericPrettyJsonPrinter{}
	case constants.OutputFormatCsv:
		return getCsvPrinter{}
	default:
		common.CliExit(fmt.Sprintf("Output format '%v' is not supported for usage. Allowed values: %v", outputFormat, constants.UsageOutputFormatFlagAllowedValuesText))
		return nil
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
			isoFormat(window.StartTime.AsTime()),
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

func isoFormat(t time.Time) string {
	n := gostradamus.DateTimeFromTime(t)
	return n.InTimezone(tz).IsoFormatTZ()
}
