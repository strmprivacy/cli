package batch_job

import (
	"context"
	"encoding/json"
	"github.com/spf13/cobra"
	"github.com/strmprivacy/api-definitions-go/v2/api/batch_jobs/v1"
	"google.golang.org/grpc"
	"io/ioutil"
	"strmprivacy/strm/pkg/auth"
	"strmprivacy/strm/pkg/common"
	"strmprivacy/strm/pkg/util"
)

var client batch_jobs.BatchJobsServiceClient
var apiContext context.Context

func SetupClient(clientConnection *grpc.ClientConn, ctx context.Context) {
	apiContext = ctx
	client = batch_jobs.NewBatchJobsServiceClient(clientConnection)
}

func list() {
	req := &batch_jobs.ListBatchJobsRequest{BillingId: auth.Auth.BillingId()}
	response, err := client.ListBatchJobs(apiContext, req)
	common.CliExit(err)

	printer.Print(response)
}

func get(id *string, _ *cobra.Command) {
	ref := &batch_jobs.BatchJobRef{
		BillingId: auth.Auth.BillingId(), Id: *id,
	}
	req := &batch_jobs.GetBatchJobRequest{Ref: ref}
	response, err := client.GetBatchJob(apiContext, req)
	common.CliExit(err)
	printer.Print(response)
}

func del(id *string) {
	req := &batch_jobs.DeleteBatchJobRequest{Ref: &batch_jobs.BatchJobRef{
		BillingId: auth.Auth.BillingId(), Id: *id}}
	response, err := client.DeleteBatchJob(apiContext, req)
	common.CliExit(err)
	printer.Print(response)
}

func create(cmd *cobra.Command) {
	flags := cmd.Flags()

	derivedDataFile := util.GetStringAndErr(flags, file)
	if len(derivedDataFile) == 0 {
		common.CliExit("Please provide a json to create a batch job")
	}

	derivedData, err := ioutil.ReadFile(derivedDataFile)
	if err != nil {
		common.CliExit(err)
	}

	job := &batch_jobs.BatchJob{}

	err = json.Unmarshal(derivedData, job)
	if err != nil {
		common.CliExit(err)
	}

	createBatchJobRequest := &batch_jobs.CreateBatchJobRequest{BatchJob: job}
	job.Ref.BillingId = auth.Auth.BillingId()

	if job.Consent.DefaultConsentLevels == nil && job.Consent.ConsentLevelExtractor == nil {
		common.CliExit("Please provide consent levels in json file")
	}

	response, err := client.CreateBatchJob(apiContext,
		createBatchJobRequest)
	common.CliExit(err)

	printer.Print(response)
}

func namesCompletion(cmd *cobra.Command, args []string, complete string) ([]string, cobra.ShellCompDirective) {
	if len(args) != 0 {
		// this one means you don't get two completion suggestions for one stream
		return nil, cobra.ShellCompDirectiveNoFileComp
	}
	if auth.Auth.BillingIdAbsent() {
		return common.MissingBillingIdCompletionError(cmd.CommandPath())
	}
	req := &batch_jobs.ListBatchJobsRequest{BillingId: auth.Auth.BillingId()}
	response, err := client.ListBatchJobs(apiContext, req)

	if err != nil {
		return common.GrpcRequestCompletionError(err)
	}

	streamNames := make([]string, 0, len(response.BatchJobs))
	for _, s := range response.BatchJobs {
		streamNames = append(streamNames, s.Ref.Id)
	}
	return streamNames, cobra.ShellCompDirectiveNoFileComp
}
