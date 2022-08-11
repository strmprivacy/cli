package project

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
			common.OutputFormatTable + common.ListCommandName:   listTablePrinter{},
			common.OutputFormatTable + common.CreateCommandName: createTablePrinter{},
			common.OutputFormatTable + common.ManageCommandName: manageTablePrinter{},
			common.OutputFormatPlain + common.ListCommandName:   listPlainPrinter{},
			common.OutputFormatPlain + common.CreateCommandName: createPlainPrinter{},
			common.OutputFormatPlain + common.ManageCommandName: managePlainPrinter{},
		},
	)
}

type listPlainPrinter struct{}
type createPlainPrinter struct{}
type managePlainPrinter struct{}

type listTablePrinter struct{}
type createTablePrinter struct{}
type manageTablePrinter struct{}

func (p listTablePrinter) Print(data interface{}) {
	listResponse, _ := (data).(*projects.ListProjectsResponse)
	printTable(listResponse.Projects)
}

func (p createTablePrinter) Print(data interface{}) {
	createResponse, _ := (data).(*projects.CreateProjectResponse)
	printTable([]*v1.Project{createResponse.Project})
}

func (p manageTablePrinter) Print(_ interface{}) {
	printTable([]*v1.Project{})
}

func (p listPlainPrinter) Print(data interface{}) {
	listResponse, _ := (data).(*projects.ListProjectsResponse)
	printPlain(listResponse.Projects)
}

func (p createPlainPrinter) Print(data interface{}) {
	createResponse, _ := (data).(*projects.CreateProjectResponse)
	printPlain([]*v1.Project{createResponse.Project})
}

func (p managePlainPrinter) Print(_ interface{}) {
	printPlain([]*v1.Project{})
}

func printTable(projects []*v1.Project) {
	rows := make([]table.Row, 0, len(projects))

	for _, project := range projects {
		rows = append(rows, table.Row{
			project.Name,
			project.Description,
		})
	}

	util.RenderTable(
		table.Row{
			"Name",
			"Description",
		},
		rows,
	)
}

func printPlain(projects []*v1.Project) {
	var names string
	lastIndex := len(projects) - 1

	for index, project := range projects {
		names = names + project.Name

		if index != lastIndex {
			names = names + "\n"
		}
	}

	util.RenderPlain(names)
}
