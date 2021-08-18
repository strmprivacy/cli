package event_contract

import (
	"fmt"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/spf13/cobra"
	v1 "github.com/streammachineio/api-definitions-go/api/entities/v1"
	"github.com/streammachineio/api-definitions-go/api/event_contracts/v1"
	"streammachine.io/strm/pkg/common"
	"streammachine.io/strm/pkg/constants"
	"streammachine.io/strm/pkg/entity/schema"
	"streammachine.io/strm/pkg/util"
)

var printer util.Printer

func configurePrinter(command *cobra.Command) util.Printer {
	outputFormat := util.GetStringAndErr(command.Flags(), constants.OutputFormatFlag)

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

func (p listTablePrinter) Print(data interface{}) {
	listResponse, _ := (data).(*event_contracts.ListEventContractsResponse)
	printTable(listResponse.EventContracts)
}

func (p getTablePrinter) Print(data interface{}) {
	getResponse, _ := (data).(*event_contracts.GetEventContractResponse)
	printTable([]*v1.EventContract{getResponse.EventContract})
}

func (p createTablePrinter) Print(data interface{}) {
	createResponse, _ := (data).(*event_contracts.CreateEventContractResponse)
	printTable([]*v1.EventContract{createResponse.EventContract})
}

func (p listPlainPrinter) Print(data interface{}) {
	listResponse, _ := (data).(*event_contracts.ListEventContractsResponse)
	printPlain(listResponse.EventContracts)
}

func (p getPlainPrinter) Print(data interface{}) {
	getResponse, _ := (data).(*event_contracts.GetEventContractResponse)
	printPlain([]*v1.EventContract{getResponse.EventContract})
}

func (p createPlainPrinter) Print(data interface{}) {
	createResponse, _ := (data).(*event_contracts.CreateEventContractResponse)
	printPlain([]*v1.EventContract{createResponse.EventContract})
}

func printTable(contracts []*v1.EventContract) {
	rows := make([]table.Row, 0, len(contracts))

	for _, contract := range contracts {
		rows = append(rows, table.Row{
			refToString(contract.Ref),
			schema.RefToString(contract.SchemaRef),
			contract.IsPublic,
			contract.KeyField,
			len(contract.PiiFields),
			len(contract.Validations),
		})
	}

	util.RenderTable(
		table.Row{
			"Event Contract",
			"Schema",
			"Public",
			"Key Field",
			"# PII Fields",
			"# Validations",
		},
		rows,
	)
}

func printPlain(contracts []*v1.EventContract) {
	var names string
	lastIndex := len(contracts) - 1

	for index, contract := range contracts {
		names = names + refToString(contract.Ref)

		if index != lastIndex {
			names = names + "\n"
		}
	}

	util.RenderPlain(names)
}
