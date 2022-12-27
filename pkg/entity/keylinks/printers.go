package keylinks

import (
	"fmt"
	"github.com/bykof/gostradamus"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/spf13/cobra"
	"github.com/strmprivacy/api-definitions-go/v2/api/data_subjects/v1"
	"strmprivacy/strm/pkg/common"
	"strmprivacy/strm/pkg/util"
)

import (
	"errors"
)

var printer util.Printer

var tz = gostradamus.Local()

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
	listResponse, _ := (data).(*data_subjects.ListDataSubjectKeylinksResponse)
	printTable(listResponse.DataSubjectKeylinks)
}

func (p listPlainPrinter) Print(data interface{}) {
	listResponse, _ := (data).(*data_subjects.ListDataSubjectKeylinksResponse)
	printPlain(listResponse.DataSubjectKeylinks)
}

func printTable(infos []*data_subjects.DatasubjectKeylinks) {
	rowCt := 0
	for _, dataSubject := range infos {
		rowCt += len(dataSubject.KeylinksAndTimestamps)
	}
	rows := make([]table.Row, 0, rowCt)

	for _, dataSubject := range infos {
		for _, keyLinkAndTimestamp := range dataSubject.KeylinksAndTimestamps {
			rows = append(rows, table.Row{
				dataSubject.DataSubjectId,
				keyLinkAndTimestamp.KeyLink,
				util.IsoFormat(tz, keyLinkAndTimestamp.ExpiryTime),
			})

		}

	}
	util.RenderTable(
		table.Row{
			"DataSubject",
			"KeyLink",
			"Expiry",
		},
		rows,
	)
}

func printPlain(infos []*data_subjects.DatasubjectKeylinks) {
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
