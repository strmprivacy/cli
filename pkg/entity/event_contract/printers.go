package event_contract

import (
	"fmt"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/spf13/cobra"
	v1 "github.com/streammachineio/api-definitions-go/api/entities/v1"
	"github.com/streammachineio/api-definitions-go/api/event_contracts/v1"
	"google.golang.org/protobuf/proto"
	"streammachine.io/strm/pkg/constants"
	"streammachine.io/strm/pkg/entity/schema"
	"streammachine.io/strm/pkg/util"
)

var printer util.Printer

func configurePrinter(command *cobra.Command) util.Printer {
	outputFormat := util.GetStringAndErr(command.Flags(), util.OutputFormatFlag)

	switch outputFormat {
	case "json":
		return util.GenericPrettyJsonPrinter{}
	case "json-raw":
		return util.GenericRawJsonPrinter{}
	case "table":
		switch command.Parent().Name() {
		case constants.ListCommandName:
			return ListEventContractsTablePrinter{}
		case constants.GetCommandName:
			return GetEventContractTablePrinter{}
		case constants.CreateCommandName:
			return CreateEventContractTablePrinter{}
		}

		return util.GenericPrettyJsonPrinter{}
	case "plain":
		switch command.Parent().Name() {
		case constants.ListCommandName:
			return ListEventContractPlainPrinter{}
		case constants.GetCommandName:
			return GetEventContractPlainPrinter{}
		case constants.CreateCommandName:
			return CreateEventContractPlainPrinter{}
		}

		return util.GenericPrettyJsonPrinter{}
	default:
		return util.GenericPrettyJsonPrinter{}
	}
}

type ListEventContractPlainPrinter struct{}
type GetEventContractPlainPrinter struct{}
type CreateEventContractPlainPrinter struct{}

type ListEventContractsTablePrinter struct{}
type GetEventContractTablePrinter struct{}
type CreateEventContractTablePrinter struct{}

type DeleteEventContractPrinter struct{}

func (p ListEventContractsTablePrinter) Print(data proto.Message) {
	listResponse, _ := (data).(*event_contracts.ListEventContractsResponse)
	printTable(listResponse.EventContracts)
}

func (p GetEventContractTablePrinter) Print(data proto.Message) {
	getResponse, _ := (data).(*event_contracts.GetEventContractResponse)
	printTable([]*v1.EventContract{getResponse.EventContract})
}

func (p CreateEventContractTablePrinter) Print(data proto.Message) {
	createResponse, _ := (data).(*event_contracts.CreateEventContractResponse)
	printTable([]*v1.EventContract{createResponse.EventContract})
}

func (p ListEventContractPlainPrinter) Print(data proto.Message) {
	listResponse, _ := (data).(*event_contracts.ListEventContractsResponse)
	printPlain(listResponse.EventContracts)
}

func (p GetEventContractPlainPrinter) Print(data proto.Message) {
	getResponse, _ := (data).(*event_contracts.GetEventContractResponse)
	printPlain([]*v1.EventContract{getResponse.EventContract})
}

func (p CreateEventContractPlainPrinter) Print(data proto.Message) {
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

	fmt.Println(names)
}
