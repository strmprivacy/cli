package datasubject

import (
	"errors"
	"fmt"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/spf13/cobra"
	"github.com/strmprivacy/api-definitions-go/v2/api/data_subjects/v1"
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
	listResponse, _ := (data).(*data_subjects.ListDataSubjectsResponse)
	printTable(listResponse.DataSubjects)
}

func (p listPlainPrinter) Print(data interface{}) {
	listResponse, _ := (data).(*data_subjects.ListDataSubjectsResponse)
	printPlain(listResponse.DataSubjects)
}

func printTable(infos []*data_subjects.ListDataSubjectsResponse_DataSubjectInfo) {
	rows := make([]table.Row, 0, len(infos))

	for _, info := range infos {
		rows = append(rows, table.Row{
			info.DataSubjectId,
		})
	}

	util.RenderTable(
		table.Row{
			"DataSubject",
		},
		rows,
	)
}

func printPlain(infos []*data_subjects.ListDataSubjectsResponse_DataSubjectInfo) {
	var names string
	lastIndex := len(infos) - 1

	for index, info := range infos {
		names = names + info.DataSubjectId

		if index != lastIndex {
			names = names + "\n"
		}
	}
	util.RenderPlain(names)
}