package member

import (
	"errors"
	"fmt"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/spf13/cobra"
	v1 "github.com/strmprivacy/api-definitions-go/v2/api/entities/v1"
	"github.com/strmprivacy/api-definitions-go/v2/api/projects/v1"
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
			common.OutputFormatTable + common.ListCommandName: listTablePrinter{},
			common.OutputFormatPlain + common.ListCommandName: listPlainPrinter{},
		},
	)
}

type listPlainPrinter struct{}

type listTablePrinter struct{}

func (p listTablePrinter) Print(data interface{}) {
	listMembersResponse, _ := (data).(*projects.ListProjectMembersResponse)
	printTable(listMembersResponse.ProjectMembers)
}

func (p listPlainPrinter) Print(data interface{}) {
	listMembersResponse, _ := (data).(*projects.ListProjectMembersResponse)
	printPlain(listMembersResponse.ProjectMembers)
}

func printTable(users []*v1.User) {
	rows := make([]table.Row, 0, len(users))

	for _, user := range users {
		rows = append(rows, table.Row{
			user.Email,
			user.FirstName,
			user.LastName,
		})
	}

	util.RenderTable(
		table.Row{
			"Email",
			"First Name",
			"Last Name",
		},
		rows,
	)
}

func printPlain(users []*v1.User) {
	var names string
	lastIndex := len(users) - 1

	for index, user := range users {
		names = names + user.Email

		if index != lastIndex {
			names = names + "\n"
		}
	}

	util.RenderPlain(names)
}
