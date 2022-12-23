package logs

import (
	"errors"
	"fmt"
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

	if follow {
		p = availablePrinters()[outputFormat+command.Parent().Name()+followFlag]
	} else {
		p = availablePrinters()[outputFormat+command.Parent().Name()]
	}

	if p == nil {
		common.CliExit(errors.New(fmt.Sprintf("Output format '%v' is not supported. Allowed values: %v", outputFormat, common.LogsOutputFormatFlagAllowedValuesText)))
	}

	return p
}

func availablePrinters() map[string]util.Printer {
	return map[string]util.Printer{
		common.OutputFormatPlain + common.LogsCommandName:              logsLatestPlainPrinter{},
		common.OutputFormatPlain + common.LogsCommandName + followFlag: logsFollowPlainPrinter{},
	}
}

type logsLatestPlainPrinter struct{}
type logsFollowPlainPrinter struct{}

func (p logsLatestPlainPrinter) Print(data interface{}) {
	response, _ := (data).(*monitoring.GetLatestEntityStatesResponse)

	for _, state := range response.State {
		printLogEntries(state)
	}
}

func (p logsFollowPlainPrinter) Print(data interface{}) {
	resp, _ := (data).(*monitoring.GetEntityStateResponse)
	printLogEntries(resp.State)
}

func printLogEntries(state *monitoring.EntityState) {
	for _, logLine := range state.Logs {
		if len(logLine) > 0 {
			println(logLine)
		}
	}
}
