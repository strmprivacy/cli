package schema

import (
	"fmt"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/spf13/cobra"
	v1 "github.com/strmprivacy/api-definitions-go/api/entities/v1"
	"github.com/strmprivacy/api-definitions-go/api/schemas/v1"
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
	listResponse, _ := (data).(*schemas.ListSchemasResponse)
	printTable(listResponse.Schemas)
}

func (p getTablePrinter) Print(data interface{}) {
	getResponse, _ := (data).(*schemas.GetSchemaResponse)
	printTable([]*v1.Schema{getResponse.Schema})
}

func (p createTablePrinter) Print(data interface{}) {
	createResponse, _ := (data).(*schemas.CreateSchemaResponse)
	printTable([]*v1.Schema{createResponse.Schema})
}

func (p listPlainPrinter) Print(data interface{}) {
	listResponse, _ := (data).(*schemas.ListSchemasResponse)
	printPlain(listResponse.Schemas)
}

func (p getPlainPrinter) Print(data interface{}) {
	getResponse, _ := (data).(*schemas.GetSchemaResponse)
	printPlain([]*v1.Schema{getResponse.Schema})
}

func (p createPlainPrinter) Print(data interface{}) {
	createResponse, _ := (data).(*schemas.CreateSchemaResponse)
	printPlain([]*v1.Schema{createResponse.Schema})
}

func printTable(schemas []*v1.Schema) {
	rows := make([]table.Row, 0, len(schemas))

	for _, schema := range schemas {
		rows = append(rows, table.Row{
			RefToString(schema.Ref),
			schema.Ref.SchemaType.String(),
			schema.IsPublic,
			schema.Fingerprint,
		})
	}

	util.RenderTable(
		table.Row{
			"Schema",
			"Type",
			"Public",
			"Fingerprint",
		},
		rows,
	)
}

func printPlain(schemas []*v1.Schema) {
	var names string
	lastIndex := len(schemas) - 1

	for index, schema := range schemas {
		names = names + RefToString(schema.Ref)

		if index != lastIndex {
			names = names + "\n"
		}
	}

	util.RenderPlain(names)
}
