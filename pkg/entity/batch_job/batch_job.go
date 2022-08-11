package batch_job

import (
	"bytes"
	"context"
	"github.com/golang/protobuf/jsonpb"
	"github.com/spf13/cobra"
	"github.com/strmprivacy/api-definitions-go/v2/api/batch_jobs/v1"
	"github.com/strmprivacy/api-definitions-go/v2/api/entities/v1"
	"google.golang.org/grpc"
	"io/ioutil"
	"strings"
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
	req := &batch_jobs.ListBatchJobsRequest{
		ProjectId: common.ProjectId,
	}
	response, err := client.ListBatchJobs(apiContext, req)
	common.CliExit(err)

	printer.Print(response)
}

func get(id *string, _ *cobra.Command) {
	ref := &entities.BatchJobRef{
		ProjectId: common.ProjectId,
		Id:        *id,
	}
	req := &batch_jobs.GetBatchJobRequest{Ref: ref}
	response, err := client.GetBatchJob(apiContext, req)
	common.CliExit(err)
	printer.Print(response)
}

func del(id *string) {
	req := &batch_jobs.DeleteBatchJobRequest{
		Ref: &entities.BatchJobRef{
			ProjectId: common.ProjectId,
			Id:        *id,
		},
	}
	response, err := client.DeleteBatchJob(apiContext, req)
	common.CliExit(err)
	printer.Print(response)
}

func create(cmd *cobra.Command) {
	flags := cmd.Flags()
	batchJobFile := util.GetStringAndErr(flags, batchJobsFileFlagName)

	batchJobData, err := ioutil.ReadFile(batchJobFile)
	if err != nil {
		common.CliExit(err)
	}

	batchJob := &entities.BatchJob{}
	err = jsonpb.Unmarshal(bytes.NewReader(batchJobData), batchJob)
	if err != nil {
		common.CliExit(err)
	}

	createBatchJobRequest := &batch_jobs.CreateBatchJobRequest{BatchJob: batchJob}
	batchJob.Ref.ProjectId = common.ProjectId

	response, err := client.CreateBatchJob(apiContext, createBatchJobRequest)
	common.CliExit(err)

	printer.Print(response)
}

func namesCompletion(cmd *cobra.Command, args []string, complete string) ([]string, cobra.ShellCompDirective) {
	if len(args) != 0 && strings.Fields(cmd.Short)[0] != "Delete" {
		// this one means you don't get multiple completion suggestions for one stream if it's not a delete call
		return nil, cobra.ShellCompDirectiveNoFileComp
	}
	req := &batch_jobs.ListBatchJobsRequest{
		ProjectId: common.ProjectId,
	}
	response, err := client.ListBatchJobs(apiContext, req)

	if err != nil {
		return common.GrpcRequestCompletionError(err)
	}

	batchJobIds := make([]string, 0, len(response.BatchJobs))
	for _, s := range response.BatchJobs {
		batchJobIds = append(batchJobIds, s.Ref.Id)
	}
	return batchJobIds, cobra.ShellCompDirectiveNoFileComp
}
