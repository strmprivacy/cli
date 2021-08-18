package schema

import (
	"fmt"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/spf13/cobra"
	v1 "github.com/streammachineio/api-definitions-go/api/entities/v1"
	"github.com/streammachineio/api-definitions-go/api/schemas/v1"
	"google.golang.org/protobuf/proto"
	"streammachine.io/strm/pkg/common"
	"streammachine.io/strm/pkg/constants"
	"streammachine.io/strm/pkg/util"
)

var printer util.Printer

func configurePrinter(command *cobra.Command) util.Printer {
	outputFormat := util.GetStringAndErr(command.Flags(), util.OutputFormatFlag)

	p := availablePrinters()[outputFormat+command.Parent().Name()]

	if p == nil {
		common.CliExit(fmt.Sprintf("Output format '%v' is not supported. Allowed values: %v", outputFormat, constants.OutputFormatFlagAllowedValuesText))
	}

	return p
}

func availablePrinters() map[string]util.Printer {
	return util.MergePrinterMaps(
		util.DefaultPrinters,
		map[string]util.Printer{
			constants.OutputFormatTable + constants.ListCommandName:   listTablePrinter{},
			constants.OutputFormatTable + constants.GetCommandName:    getTablePrinter{},
			constants.OutputFormatTable + constants.CreateCommandName: createTablePrinter{},
			constants.OutputFormatPlain + constants.ListCommandName:   listPlainPrinter{},
			constants.OutputFormatPlain + constants.GetCommandName:    getPlainPrinter{},
			constants.OutputFormatPlain + constants.CreateCommandName: createPlainPrinter{},
		},
	)
}

type listPlainPrinter struct{}
type getPlainPrinter struct{}
type createPlainPrinter struct{}

type listTablePrinter struct{}
type getTablePrinter struct{}
type createTablePrinter struct{}

func (p listTablePrinter) Print(data proto.Message) {
	listResponse, _ := (data).(*schemas.ListSchemasResponse)
	printTable(listResponse.Schemas)
}

func (p getTablePrinter) Print(data proto.Message) {
	getResponse, _ := (data).(*schemas.GetSchemaResponse)
	printTable([]*v1.Schema{getResponse.Schema})
}

func (p createTablePrinter) Print(data proto.Message) {
	createResponse, _ := (data).(*schemas.CreateSchemaResponse)
	printTable([]*v1.Schema{createResponse.Schema})
}

func (p listPlainPrinter) Print(data proto.Message) {
	listResponse, _ := (data).(*schemas.ListSchemasResponse)
	printPlain(listResponse.Schemas)
}

func (p getPlainPrinter) Print(data proto.Message) {
	getResponse, _ := (data).(*schemas.GetSchemaResponse)
	printPlain([]*v1.Schema{getResponse.Schema})
}

func (p createPlainPrinter) Print(data proto.Message) {
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
