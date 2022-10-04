package project

import (
	"errors"
	"fmt"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/spf13/cobra"
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
			common.OutputFormatTable + common.GetCommandName:    getTablePrinter{},
			common.OutputFormatTable + common.DeleteCommandName: deleteTablePrinter{},
			common.OutputFormatPlain + common.ListCommandName:   listPlainPrinter{},
			common.OutputFormatPlain + common.CreateCommandName: createPlainPrinter{},
			common.OutputFormatPlain + common.ManageCommandName: managePlainPrinter{},
			common.OutputFormatPlain + common.GetCommandName:    getPlainPrinter{},
			common.OutputFormatPlain + common.DeleteCommandName: deletePlainPrinter{},
		},
	)
}

type listPlainPrinter struct{}
type createPlainPrinter struct{}
type managePlainPrinter struct{}
type getPlainPrinter struct{}
type deletePlainPrinter struct{}

type listTablePrinter struct{}
type createTablePrinter struct{}
type manageTablePrinter struct{}
type getTablePrinter struct{}
type deleteTablePrinter struct{}

func (p listTablePrinter) Print(data interface{}) {
	listResponse, _ := (data).(ProjectsAndActiveProject)
	printTable(listResponse)
}

func (p createTablePrinter) Print(data interface{}) {
	createResponse, _ := (data).(ProjectsAndActiveProject)
	printTable(createResponse)
}

func (p manageTablePrinter) Print(_ interface{}) {
	printTable(ProjectsAndActiveProject{})
}

func (p getTablePrinter) Print(data interface{}) {
	getResponse, _ := (data).(ProjectsAndActiveProject)
	printTable(getResponse)
}

func (p deleteTablePrinter) Print(data interface{}) {
	printTable(ProjectsAndActiveProject{})
}

func (p listPlainPrinter) Print(data interface{}) {
	listResponse, _ := (data).(ProjectsAndActiveProject)
	printPlain(listResponse)
}

func (p createPlainPrinter) Print(data interface{}) {
	createResponse, _ := (data).(ProjectsAndActiveProject)
	printPlain(createResponse)
}

func (p managePlainPrinter) Print(_ interface{}) {
	printPlain(ProjectsAndActiveProject{})
}

func (p getPlainPrinter) Print(data interface{}) {
	getResponse, _ := (data).(ProjectsAndActiveProject)
	printPlain(getResponse)
}

func (p deletePlainPrinter) Print(_ interface{}) {
	printPlain(ProjectsAndActiveProject{})
}

func printTable(projectsWithActive ProjectsAndActiveProject) {
	rows := make([]table.Row, 0, len(projectsWithActive.Projects))

	for _, project := range projectsWithActive.Projects {
		projectName := ""
		if project.Name == projectsWithActive.activeProject {
			projectName = project.Name + " (active)"
		} else {
			projectName = project.Name
		}
		rows = append(rows, table.Row{
			projectName,
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

func printPlain(projectsWithActive ProjectsAndActiveProject) {
	var names string
	lastIndex := len(projectsWithActive.Projects) - 1

	for index, project := range projectsWithActive.Projects {
		if project.Name == projectsWithActive.activeProject {
			names = names + project.Name + " (active)"
		} else {
			names = names + project.Name
		}
		if index != lastIndex {
			names = names + "\n"
		}
	}

	util.RenderPlain(names)
}
