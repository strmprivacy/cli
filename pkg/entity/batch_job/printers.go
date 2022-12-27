package batch_job

import (
	"errors"
	"fmt"
	"github.com/bykof/gostradamus"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/spf13/cobra"
	"github.com/strmprivacy/api-definitions-go/v2/api/batch_jobs/v1"
	"github.com/strmprivacy/api-definitions-go/v2/api/entities/v1"
	"sort"
	"strmprivacy/strm/pkg/common"
	"strmprivacy/strm/pkg/util"
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
			common.OutputFormatTable + common.ListCommandName:   listTablePrinter{},
			common.OutputFormatTable + common.GetCommandName:    getTablePrinter{},
			common.OutputFormatTable + common.DeleteCommandName: deletePrinter{},
			common.OutputFormatTable + common.CreateCommandName: createTablePrinter{},
			common.OutputFormatPlain + common.ListCommandName:   listPlainPrinter{},
			common.OutputFormatPlain + common.GetCommandName:    getPlainPrinter{},
			common.OutputFormatPlain + common.DeleteCommandName: deletePrinter{},
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

type deletePrinter struct{}

func (p listTablePrinter) Print(data interface{}) {
	listResponse, _ := (data).(*batch_jobs.ListBatchJobsResponse)
	printTable(listResponse.Jobs)
}

func (p getTablePrinter) Print(data interface{}) {
	getResponse, _ := (data).(*batch_jobs.GetBatchJobResponse)
	printTable([]*entities.BatchJobWrapper{getResponse.Job})
}

func (p createTablePrinter) Print(data interface{}) {
	createResponse, _ := (data).(*batch_jobs.CreateBatchJobResponse)
	printTable([]*entities.BatchJobWrapper{createResponse.Job})
}

func (p listPlainPrinter) Print(data interface{}) {
	listResponse, _ := (data).(*batch_jobs.ListBatchJobsResponse)
	printPlain(listResponse.Jobs)
}

func (p getPlainPrinter) Print(data interface{}) {
	getResponse, _ := (data).(*batch_jobs.GetBatchJobResponse)
	printPlain([]*entities.BatchJobWrapper{getResponse.Job})
}

func (p createPlainPrinter) Print(data interface{}) {
	createResponse, _ := (data).(*batch_jobs.CreateBatchJobResponse)
	printPlain([]*entities.BatchJobWrapper{createResponse.Job})
}

func (p deletePrinter) Print(data interface{}) {
	fmt.Println("Batch Job has been deleted")
}

func printTable(batchJobs []*entities.BatchJobWrapper) {
	rows := make([]table.Row, 0, len(batchJobs))
	for _, batchJob := range batchJobs {
		batchJobRefWithStates := toRefWithStates(batchJob)
		states := batchJobRefWithStates.states[:]
		sort.Slice(states, func(i, j int) bool {
			// Reverse sort, j > i
			return states[j].StateTime.AsTime().Before(states[i].StateTime.AsTime())
		})

		batchJobState := states[0]

		var message = ""
		if batchJobState.State == entities.BatchJobStateType_ERROR {
			message = batchJobState.Message
		}

		rows = append(rows, table.Row{
			batchJobRefWithStates.ref.Id,
			util.IsoFormat(tz, batchJobState.StateTime),
			batchJobState.State.String(),
			message,
		})
	}

	util.RenderTable(
		table.Row{
			"Batch Job id",
			"Timestamp",
			"State",
			"Details",
		},
		rows,
	)
}

func printPlain(batchJobs []*entities.BatchJobWrapper) {
	var ids string
	lastIndex := len(batchJobs) - 1

	for index, batchJob := range batchJobs {
		ids = ids + toRefWithStates(batchJob).ref.Id

		if index != lastIndex {
			ids = ids + "\n"
		}
	}

	util.RenderPlain(ids)
}
