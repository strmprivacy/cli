package datasubject

import (
	"errors"
	"fmt"
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
			common.OutputFormatPlain + common.ListCommandName:       listPlainPrinter{},
			common.OutputFormatPlain + "0" + common.ListCommandName: listPlain0Printer{},
		},
	)
}

type listPlainPrinter struct{}
type listPlain0Printer struct{}

func (p listPlainPrinter) Print(data interface{}) {
	listResponse, _ := (data).(*data_subjects.ListDataSubjectsResponse)
	printPlain(listResponse, false)
}
func (p listPlain0Printer) Print(data interface{}) {
	listResponse, _ := (data).(*data_subjects.ListDataSubjectsResponse)
	printPlain(listResponse, true)
}

func printPlain(response *data_subjects.ListDataSubjectsResponse, print0 bool) {
	var sep = '\n'
	if print0 {
		sep = '\000'
	}
	fmt.Printf("%s%c", response.NextPageToken, sep)
	for _, dataSubject := range response.DataSubjects {
		fmt.Printf("%s%c", dataSubject.DataSubjectId, sep)
	}
}