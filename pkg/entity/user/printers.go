package user

import (
	"errors"
	"fmt"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/spf13/cobra"
	v1 "github.com/strmprivacy/api-definitions-go/v2/api/entities/v1"
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
			common.OutputFormatTable + common.GetCommandName:  getTablePrinter{},
			common.OutputFormatPlain + common.GetCommandName:  getPlainPrinter{},
		},
	)
}

type listPlainPrinter struct{}
type getPlainPrinter struct{}

type listTablePrinter struct{}
type getTablePrinter struct{}

func (p listTablePrinter) Print(data interface{}) {
	user, _ := (data).([]*v1.User)
	printTable(user)
}

func (p listPlainPrinter) Print(data interface{}) {
	user, _ := (data).([]*v1.User)
	printPlain(user)
}

func (p getTablePrinter) Print(data interface{}) {
	user, _ := (data).(*v1.User)
	printTable([]*v1.User{user})
}

func (p getPlainPrinter) Print(data interface{}) {
	user, _ := (data).(*v1.User)
	printPlain([]*v1.User{user})
}

func printTable(users []*v1.User) {
	rows := make([]table.Row, 0, len(users))

	for _, user := range users {
		rows = append(rows, table.Row{
			user.Email,
			user.FirstName,
			user.LastName,
			user.UserRoles,
		})
	}

	util.RenderTable(
		table.Row{
			"Email",
			"First Name",
			"Last Name",
			"User Roles",
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
