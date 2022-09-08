package data_contract

import (
	"errors"
	"fmt"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/spf13/cobra"
	"github.com/strmprivacy/api-definitions-go/v2/api/data_contracts/v1"
	"github.com/strmprivacy/api-definitions-go/v2/api/entities/v1"
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
			common.OutputFormatTable + common.CreateCommandName: createTablePrinter{},
			common.OutputFormatPlain + common.CreateCommandName: createPlainPrinter{},
			common.OutputFormatTable + common.ListCommandName:   listTablePrinter{},
			common.OutputFormatPlain + common.ListCommandName:   listPlainPrinter{},
			common.OutputFormatTable + common.GetCommandName:    getTablePrinter{},
			common.OutputFormatPlain + common.GetCommandName:    getPlainPrinter{},
		},
	)
}

type createTablePrinter struct{}
type createPlainPrinter struct{}
type listTablePrinter struct{}
type listPlainPrinter struct{}
type getTablePrinter struct{}
type getPlainPrinter struct{}

func (p createTablePrinter) Print(data interface{}) {
	createResponse, _ := (data).(*data_contracts.CreateDataContractResponse)
	printTable([]*v1.DataContract{createResponse.DataContract})
}

func (p createPlainPrinter) Print(data interface{}) {
	createResponse, _ := (data).(*data_contracts.CreateDataContractResponse)
	printPlain([]*v1.DataContract{createResponse.DataContract})
}

func (p listTablePrinter) Print(data interface{}) {
	listResponse, _ := (data).(*data_contracts.ListDataContractsResponse)
	printTable(listResponse.DataContracts)
}

func (p listPlainPrinter) Print(data interface{}) {
	listResponse, _ := (data).(*data_contracts.ListDataContractsResponse)
	printPlain(listResponse.DataContracts)
}

func (p getTablePrinter) Print(data interface{}) {
	getResponse, _ := (data).(*data_contracts.GetDataContractResponse)
	printTable([]*v1.DataContract{getResponse.DataContract})
}

func (p getPlainPrinter) Print(data interface{}) {
	getResponse, _ := (data).(*data_contracts.GetDataContractResponse)
	printPlain([]*v1.DataContract{getResponse.DataContract})
}

func refToString(ref *entities.DataContractRef) string {
	return fmt.Sprintf("%v/%v/%v", ref.Handle, ref.Name, ref.Version)
}

func printTable(contracts []*v1.DataContract) {
	rows := make([]table.Row, 0, len(contracts))

	for _, contract := range contracts {
		rows = append(rows, table.Row{
			refToString(contract.Ref),
			contract.State,
			contract.IsPublic,
			contract.KeyField,
			len(contract.PiiFields),
			len(contract.Validations),
		})
	}

	util.RenderTable(
		table.Row{
			"Event Contract",
			"State",
			"Public",
			"Key Field",
			"# PII Fields",
			"# Validations",
		},
		rows,
	)
}
func printPlain(contracts []*v1.DataContract) {
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
