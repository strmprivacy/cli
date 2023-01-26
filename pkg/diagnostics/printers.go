package diagnostics

import (
	"errors"
	"fmt"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/spf13/cobra"
	"math"
	"strconv"
	"strmprivacy/strm/pkg/common"
	"strmprivacy/strm/pkg/util"
)

var printer util.Printer

func configurePrinter(command *cobra.Command) util.Printer {
	outputFormat := util.GetStringAndErr(command.Flags(), common.OutputFormatFlag)

	p := availablePrinters()[outputFormat+command.Parent().Name()]

	if p == nil {
		common.CliExit(errors.New(fmt.Sprintf("Output format '%v' is not supported. Allowed values: %v", outputFormat, common.OutputFormatFlagAllowedValuesText)))
	}

	return p
}

func availablePrinters() map[string]util.Printer {
	return util.MergePrinterMaps(
		util.DefaultPrinters,
		map[string]util.Printer{
			common.OutputFormatTable + common.EvaluateCommandName: diagnosticsTablePrinter{},
			common.OutputFormatPlain + common.EvaluateCommandName: diagnosticsPlainPrinter{},
		},
	)
}

type diagnosticsTablePrinter struct{}
type diagnosticsPlainPrinter struct{}

func (p diagnosticsTablePrinter) Print(data interface{}) {
	metrics, _ := (data).(Metrics)
	printTable(metrics)
}

func (p diagnosticsPlainPrinter) Print(data interface{}) {
	metrics, _ := (data).(Metrics)
	printPlain(metrics)
}

func printTable(metricsResponse Metrics) {
	nRows := metricsResponse.KAnonymity/1 + len(metricsResponse.LDiversity) + int(math.Ceil(metricsResponse.TCloseness/1))
	rows := make([]table.Row, 0, nRows)
	rows = append(rows, table.Row{"k-Anonymity", metricsResponse.KAnonymity, ""})
	ix := 0
	for column, ld := range metricsResponse.LDiversity {
		if ix == 0 {
			rows = append(rows, table.Row{"l-Diversity", ld, column})
		} else {
			rows = append(rows, table.Row{"", ld, column})
		}
		ix += 1
	}
	rows = append(rows, table.Row{"t-Closeness", metricsResponse.TCloseness, ""})

	util.RenderTable(
		table.Row{
			"Metric",
			"value",
			"Sensitive Attribute",
		},
		rows,
	)
}

func printPlain(metricsResponse Metrics) {
	var metrics string
	metrics = metrics + "k-Anonymity: " + strconv.Itoa(metricsResponse.KAnonymity) + "\n"
	metrics = metrics + "l-Diversity: "
	for k, v := range metricsResponse.LDiversity {
		metrics = metrics + k + ": " + strconv.Itoa(v) + "\t"
	}
	metrics = metrics + "\n"
	metrics = metrics + "t-Closeness: " + strconv.FormatFloat(metricsResponse.TCloseness, byte('e'), 5, 32) + "\n"
	util.RenderPlain(metrics)
}
