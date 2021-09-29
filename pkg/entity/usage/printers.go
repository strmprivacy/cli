package usage

import (
	"fmt"
	"github.com/bykof/gostradamus"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/spf13/cobra"
	"github.com/streammachineio/api-definitions-go/api/usage/v1"
	"math"
	"streammachine.io/strm/pkg/common"
	"streammachine.io/strm/pkg/util"
	"time"
)

var printer util.Printer

func configurePrinter(command *cobra.Command) util.Printer {
	outputFormat := util.GetStringAndErr(command.Flags(), common.OutputFormatFlag)

	p := availablePrinters()[outputFormat]

	if p == nil {
		common.CliExit(fmt.Sprintf("Output format '%v' is not supported for usage. Allowed values: %v", outputFormat, common.UsageOutputFormatFlagAllowedValuesText))
	}

	return p
}

func availablePrinters() map[string]util.Printer {
	return map[string]util.Printer{
		common.OutputFormatJsonRaw: util.ProtoMessageJsonRawPrinter{},
		common.OutputFormatJson:    util.ProtoMessageJsonPrettyPrinter{},
		common.OutputFormatCsv:     getCsvPrinter{},
	}
}

type getCsvPrinter struct{}

func (p getCsvPrinter) Print(data interface{}) {
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
