package batch_job

import (
	"fmt"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/spf13/cobra"
	"github.com/strmprivacy/api-definitions-go/v2/api/batch_jobs/v1"
	"sort"
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
	printTable(listResponse.BatchJobs)
}

func (p getTablePrinter) Print(data interface{}) {
	getResponse, _ := (data).(*batch_jobs.GetBatchJobResponse)
	printTable([]*batch_jobs.BatchJob{getResponse.BatchJob})
}

func (p createTablePrinter) Print(data interface{}) {
	createResponse, _ := (data).(*batch_jobs.CreateBatchJobResponse)
	printTable([]*batch_jobs.BatchJob{createResponse.BatchJob})
}

func (p listPlainPrinter) Print(data interface{}) {
	listResponse, _ := (data).(*batch_jobs.ListBatchJobsResponse)
	printPlain(listResponse.BatchJobs)
}

func (p getPlainPrinter) Print(data interface{}) {
	getResponse, _ := (data).(*batch_jobs.GetBatchJobResponse)
	printPlain([]*batch_jobs.BatchJob{getResponse.BatchJob})
}

func (p createPlainPrinter) Print(data interface{}) {
	createResponse, _ := (data).(*batch_jobs.CreateBatchJobResponse)
	printPlain([]*batch_jobs.BatchJob{createResponse.BatchJob})
}

func (p deletePrinter) Print(data interface{}) {
	fmt.Println("Batch Job has been deleted")
}

func printTable(batchJobs []*batch_jobs.BatchJob) {
	rows := make([]table.Row, 0, len(batchJobs))
	for _, batchJob := range batchJobs {

		states := batchJob.States[:]
		sort.Slice(states, func(i, j int) bool {
			// Reverse sort, j > i
			return states[j].StateTime.AsTime().Before(states[i].StateTime.AsTime())
		})

		rows = append(rows, table.Row{
			batchJob.Ref.Id,
			states[0].State.String(),
			states[0].StateTime.AsTime(),
		})
	}

	util.RenderTable(
		table.Row{
			"Batch Job id",
			"State",
			"Timestamp",
		},
		rows,
	)
}

func printPlain(batchJobs []*batch_jobs.BatchJob) {
	var ids string
	lastIndex := len(batchJobs) - 1

	for index, batchJob := range batchJobs {
		ids = ids + batchJob.Ref.Id

		if index != lastIndex {
			ids = ids + "\n"
		}
	}

	util.RenderPlain(ids)
}
