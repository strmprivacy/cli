package purpose_mapping

import (
	"errors"
	"fmt"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/spf13/cobra"
	v1 "github.com/strmprivacy/api-definitions-go/v3/api/entities/v1"
	"github.com/strmprivacy/api-definitions-go/v3/api/purpose_mapping/v1"
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
			common.OutputFormatTable + common.GetCommandName:    getTablePrinter{},
			common.OutputFormatTable + common.CreateCommandName: createTablePrinter{},
			common.OutputFormatPlain + common.ListCommandName:   listPlainPrinter{},
			common.OutputFormatPlain + common.GetCommandName:    getPlainPrinter{},
			common.OutputFormatPlain + common.CreateCommandName: createPlainPrinter{},
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
	listResponse, _ := (data).(*purpose_mapping.ListPurposeMappingsResponse)
	printTable(listResponse.PurposeMappings)
}

func (p getTablePrinter) Print(data interface{}) {
	getResponse, _ := (data).(*purpose_mapping.GetPurposeMappingResponse)
	printTable([]*v1.PurposeMapping{getResponse.PurposeMapping})
}

func (p createTablePrinter) Print(data interface{}) {
	createResponse, _ := (data).(*purpose_mapping.CreatePurposeMappingResponse)
	printTable([]*v1.PurposeMapping{createResponse.PurposeMapping})
}

func (p listPlainPrinter) Print(data interface{}) {
	listResponse, _ := (data).(*purpose_mapping.ListPurposeMappingsResponse)
	printPlain(listResponse.PurposeMappings)
}

func (p getPlainPrinter) Print(data interface{}) {
	getResponse, _ := (data).(*purpose_mapping.GetPurposeMappingResponse)
	printPlain([]*v1.PurposeMapping{getResponse.PurposeMapping})
}

func (p createPlainPrinter) Print(data interface{}) {
	createResponse, _ := (data).(*purpose_mapping.CreatePurposeMappingResponse)
	printPlain([]*v1.PurposeMapping{createResponse.PurposeMapping})
}

func printTable(purposeMappings []*v1.PurposeMapping) {
	rows := make([]table.Row, 0, len(purposeMappings))

	for _, purposeMapping := range purposeMappings {
		row := table.Row{
			purposeMapping.Purpose,
			purposeMapping.Level,
			purposeMapping.Description,
		}
		rows = append(rows, row)
	}

	headers := table.Row{
		"Purpose",
		"Value",
		"Description",
	}
	util.RenderTable(headers, rows)
}

func printPlain(purposeMappings []*v1.PurposeMapping) {
	var names string
	lastIndex := len(purposeMappings) - 1

	for index, purposeMapping := range purposeMappings {
		names = names + fmt.Sprintf("%v (%v)", purposeMapping.Purpose, purposeMapping.Level)

		if index != lastIndex {
			names = names + "\n"
		}
	}

	util.RenderPlain(names)
}
