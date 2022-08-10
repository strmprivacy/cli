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
			common.OutputFormatPlain + common.ListCommandName: listPlainPrinter{},
		},
	)
}

type listPlainPrinter struct{}

func (p listPlainPrinter) Print(data interface{}) {
	listResponse, _ := (data).(*data_subjects.ListDataSubjectsResponse)
	printPlain(listResponse)
}

func printPlain(response *data_subjects.ListDataSubjectsResponse) {
	fmt.Println(response.NextPageToken)
	for _, dataSubject := range response.DataSubjects {
		fmt.Println(dataSubject.DataSubjectId)
	}
}