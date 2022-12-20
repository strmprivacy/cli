package monitor

import (
	"errors"
	"fmt"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/spf13/cobra"
	v1 "github.com/strmprivacy/api-definitions-go/v2/api/entities/v1"
	"github.com/strmprivacy/api-definitions-go/v2/api/monitoring/v1"
	"strmprivacy/strm/pkg/common"
	"strmprivacy/strm/pkg/util"
)

var printer util.Printer

func configurePrinter(command *cobra.Command) util.Printer {
	outputFormat := util.GetStringAndErr(command.Flags(), common.OutputFormatFlag)
	follow := util.GetBool(command.Flags(), followFlag)

	var p util.Printer

	if follow {
		p = availablePrinters()[outputFormat+command.Parent().Name()+followFlag]
	} else {
		p = availablePrinters()[outputFormat+command.Parent().Name()]
	}

	if p == nil {
		common.CliExit(errors.New(fmt.Sprintf("Output format '%v' is not supported. Allowed values: %v", outputFormat, common.OutputFormatFlagAllowedValuesText)))
	}

	return p
}

func availablePrinters() map[string]util.Printer {
	return util.MergePrinterMaps(
		map[string]util.Printer{
			common.OutputFormatTable + common.MonitorCommandName:              monitorLatestTablePrinter{},
			common.OutputFormatTable + common.MonitorCommandName + followFlag: monitorFollowTablePrinter{},
		},
	)
}

type monitorLatestTablePrinter struct{}
type monitorFollowTablePrinter struct{}

func (p monitorLatestTablePrinter) Print(data interface{}) {
	response, _ := (data).(*monitoring.GetLatestEntityStatesResponse)
	printTable(response.State)
}

func (p monitorFollowTablePrinter) Print(data interface{}) {
	response, _ := (data).(*monitoring.GetLatestEntityStatesResponse)
	printTable(response.State)
}

func printTable(entityStates []*monitoring.EntityState) {
	rows := make([]table.Row, 0, len(entityStates))

	for _, state := range entityStates {
		row := table.Row{
			state.StateTime.AsTime(),
			state.Ref.Type.String(),
			state.ResourceType.String(),
			state.Ref.Name,
			state.Status.String(),
			state.Message,
		}
		rows = append(rows, row)
	}

	headers := table.Row{
		"Timestamp",
		"Entity Type",
		"Resource Type",
		"Name",
		"Status",
		"Message",
	}
	util.RenderTable(headers, rows)
}

func printPlain(streamTreeArray []*v1.StreamTree) {
	var names string
	lastIndex := len(streamTreeArray) - 1

	for index, stream := range streamTreeArray {
		names = names + stream.Stream.Ref.Name

		if index != lastIndex {
			names = names + "\n"
		}
	}

	util.RenderPlain(names)
}
