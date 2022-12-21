package monitor

import (
	"errors"
	"fmt"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/spf13/cobra"
	"github.com/strmprivacy/api-definitions-go/v2/api/monitoring/v1"
	"strmprivacy/strm/pkg/common"
	"strmprivacy/strm/pkg/util"
)

var printer util.Printer

func configurePrinter(command *cobra.Command) util.Printer {
	outputFormat := util.GetStringAndErr(command.Flags(), common.OutputFormatFlag)
	follow := util.GetBool(command.Flags(), followFlag)

	var p util.Printer
	var allowedValues string

	if follow {
		p = availablePrinters()[outputFormat+command.Parent().Name()+followFlag]
		allowedValues = common.MonitorFollowOutputFormatFlagAllowedValuesText
	} else {
		p = availablePrinters()[outputFormat+command.Parent().Name()]
		allowedValues = common.MonitorOutputFormatFlagAllowedValuesText
	}

	if p == nil {
		common.CliExit(errors.New(fmt.Sprintf("Output format '%v' is not supported. Allowed values: %v", outputFormat, allowedValues)))
	}

	return p
}

func availablePrinters() map[string]util.Printer {
	return map[string]util.Printer{
		common.OutputFormatTable + common.MonitorCommandName:   monitorLatestTablePrinter{},
		common.OutputFormatPlain + common.MonitorCommandName:   monitorLatestPlainPrinter{},
		common.OutputFormatJson + common.MonitorCommandName:    util.ProtoMessageJsonPrettyPrinter{},
		common.OutputFormatJsonRaw + common.MonitorCommandName: util.ProtoMessageJsonRawPrinter{},

		common.OutputFormatPlain + common.MonitorCommandName + followFlag:   monitorFollowPlainPrinter{},
		common.OutputFormatJson + common.MonitorCommandName + followFlag:    util.ProtoMessageJsonPrettyPrinter{},
		common.OutputFormatJsonRaw + common.MonitorCommandName + followFlag: util.ProtoMessageJsonRawPrinter{},
	}
}

type monitorLatestTablePrinter struct{}
type monitorLatestPlainPrinter struct{}
type monitorFollowPlainPrinter struct{}

func (p monitorLatestTablePrinter) Print(data interface{}) {
	response, _ := (data).(*monitoring.GetLatestEntityStatesResponse)
	printTable(response.State)
}
func (p monitorLatestPlainPrinter) Print(data interface{}) {
	response, _ := (data).(*monitoring.GetLatestEntityStatesResponse)

	for _, state := range response.State {
		fmt.Printf("%s %v %s %s %s %s\n",
			util.IsoFormat(tz, state.StateTime),
			state.Ref.Type,
			state.Ref.Name,
			state.Status,
			state.ResourceType,
			state.Message,
		)
	}
}

func (p monitorFollowPlainPrinter) Print(data interface{}) {
	resp, _ := (data).(*monitoring.GetEntityStateResponse)
	fmt.Printf("%s %v %s %s %s %s\n",
		util.IsoFormat(tz, resp.State.StateTime),
		resp.State.Ref.Type,
		resp.State.Ref.Name,
		resp.State.Status,
		resp.State.ResourceType,
		resp.State.Message,
	)
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
