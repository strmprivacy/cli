package installation

import (
	"fmt"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/spf13/cobra"
	"github.com/strmprivacy/api-definitions-go/v2/api/installations/v1"
	"strmprivacy/strm/pkg/common"
	"strmprivacy/strm/pkg/util"
)

var printer util.Printer

func configurePrinter(command *cobra.Command) util.Printer {
	outputFormat := util.GetStringAndErr(command.Flags(), common.OutputFormatFlag)

	p := availablePrinters()[outputFormat+command.Parent().Name()]

	if p == nil {
		common.CliExit(fmt.Sprintf("Output format '%v' is not supported. Allowed values: %v", outputFormat, common.OutputFormatFlagAllowedValuesText))
	}

	return p
}

func availablePrinters() map[string]util.Printer {
	return util.MergePrinterMaps(
		util.DefaultPrinters,
		map[string]util.Printer{
			common.OutputFormatTable + common.ListCommandName: listTablePrinter{},
			common.OutputFormatTable + common.GetCommandName:  getTablePrinter{},
			common.OutputFormatPlain + common.ListCommandName: listPlainPrinter{},
			common.OutputFormatPlain + common.GetCommandName:  getPlainPrinter{},
		},
	)
}

type listPlainPrinter struct{}
type getPlainPrinter struct{}
type createPlainPrinter struct{}

type listTablePrinter struct{}
type getTablePrinter struct{}
type createTablePrinter struct{}

func (p listTablePrinter) Print(data interface{}) {
	listResponse, _ := (data).(*installations.ListInstallationsResponse)
	printTable(listResponse.Installations)
}

func (p getTablePrinter) Print(data interface{}) {
	getResponse, _ := (data).(*installations.GetInstallationResponse)
	printTable([]*installations.Installation{getResponse.Installation})
}

func (p listPlainPrinter) Print(data interface{}) {
	listResponse, _ := (data).(*installations.ListInstallationsResponse)
	printPlain(listResponse.Installations)
}

func (p getPlainPrinter) Print(data interface{}) {
	getResponse, _ := (data).(*installations.GetInstallationResponse)
	printPlain([]*installations.Installation{getResponse.Installation})
}

func printTable(installations []*installations.Installation) {
	rows := make([]table.Row, 0, len(installations))
	for _, installation := range installations {

		rows = append(rows, table.Row{
			installation.Id,
			installation.InstallationType.String(),
			installation.BillingAccountId,
		})
	}

	util.RenderTable(
		table.Row{
			"Installation id",
			"Type",
			"Billing Account id",
		},
		rows,
	)
}

func printPlain(installations []*installations.Installation) {
	var ids string
	lastIndex := len(installations) - 1

	for index, installation := range installations {
		ids = ids + installation.Id

		if index != lastIndex {
			ids = ids + "\n"
		}
	}

	util.RenderPlain(ids)
}
