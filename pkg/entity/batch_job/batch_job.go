package batch_job

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"github.com/golang/protobuf/jsonpb"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/strmprivacy/api-definitions-go/v2/api/batch_jobs/v1"
	"github.com/strmprivacy/api-definitions-go/v2/api/entities/v1"
	"google.golang.org/grpc"
	"os"
	"strings"
	"strmprivacy/strm/pkg/common"
	"strmprivacy/strm/pkg/entity/policy"
	"strmprivacy/strm/pkg/entity/project"
	"strmprivacy/strm/pkg/util"
)

var client batch_jobs.BatchJobsServiceClient
var apiContext context.Context

type refWithStates struct {
	ref    *entities.BatchJobRef
	states []*entities.BatchJobState
}

func SetupClient(clientConnection *grpc.ClientConn, ctx context.Context) {
	apiContext = ctx
	client = batch_jobs.NewBatchJobsServiceClient(clientConnection)
}

func list(cmd *cobra.Command) {
	req := &batch_jobs.ListBatchJobsRequest{
		ProjectId: project.GetProjectId(cmd),
	}
	response, err := client.ListBatchJobs(apiContext, req)
	common.CliExit(err)

	printer.Print(response)
}

func get(id *string, cmd *cobra.Command) {
	ref := &entities.BatchJobRef{
		ProjectId: project.GetProjectId(cmd),
		Id:        *id,
	}
	req := &batch_jobs.GetBatchJobRequest{Ref: ref}
	response, err := client.GetBatchJob(apiContext, req)
	common.CliExit(err)
	printer.Print(response)
}

func del(id *string, cmd *cobra.Command) {
	req := &batch_jobs.DeleteBatchJobRequest{
		Ref: &entities.BatchJobRef{
			ProjectId: project.GetProjectId(cmd),
			Id:        *id,
		},
	}
	response, err := client.DeleteBatchJob(apiContext, req)
	common.CliExit(err)
	printer.Print(response)
}

func create(cmd *cobra.Command) {
	flags := cmd.Flags()
	batchJobFile := util.GetStringAndErr(flags, batchJobFileFlagName)
	batchJobType := getBatchJobType(flags)

	batchJobData, err := os.ReadFile(batchJobFile)
	if err != nil {
		common.CliExit(err)
	}

	var batchJobWrapper *entities.BatchJobWrapper
	if batchJobType == encryptionType {
		batchJobWrapper = createEncryptionBatchJob(cmd, batchJobData, flags)
	} else if batchJobType == microAggregationType {
		batchJobWrapper = createMicroAggregationBatchJob(cmd, batchJobData)
	}
	createBatchJobRequest := &batch_jobs.CreateBatchJobRequest{Job: batchJobWrapper}
	response, err := client.CreateBatchJob(apiContext, createBatchJobRequest)
	common.CliExit(err)

	printer.Print(response)
}

func getBatchJobType(flags *pflag.FlagSet) string {
	batchJobType := util.GetStringAndErr(flags, batchJobTypeFlagName)
	if !(batchJobType == encryptionType || batchJobType == microAggregationType) {
		common.CliExit(errors.New(fmt.Sprintf("Batch job type should be one of: %s, %s",
			encryptionType, microAggregationType)))
	}
	return batchJobType
}

func createEncryptionBatchJob(cmd *cobra.Command, batchJobData []byte, flags *pflag.FlagSet) *entities.BatchJobWrapper {
	batchJob := &entities.BatchJob{}
	err := jsonpb.Unmarshal(bytes.NewReader(batchJobData), batchJob)
	if err != nil {
		common.CliExit(err)
	}
	projectId := project.GetProjectId(cmd)
	policyId := policy.GetPolicyFromFlags(flags)
	if policyId != "" {
		batchJob.PolicyId = policyId
	}
	setEncryptionBatchJobProjectIds(batchJob, projectId)
	return &entities.BatchJobWrapper{
		Job: &entities.BatchJobWrapper_EncryptionBatchJob{
			EncryptionBatchJob: batchJob,
		},
	}
}

func setEncryptionBatchJobProjectIds(batchJob *entities.BatchJob, projectId string) {
	if batchJob.Ref == nil {
		// normal situation where the whole ref attribute in the json is absent.
		batchJob.Ref = &entities.BatchJobRef{}
	}
	batchJob.Ref.ProjectId = projectId
	batchJob.SourceData.DataConnectorRef.ProjectId = projectId
	batchJob.EncryptedData.Target.DataConnectorRef.ProjectId = projectId
	batchJob.EncryptionKeysData.Target.DataConnectorRef.ProjectId = projectId
	for _, d := range batchJob.DerivedData {
		d.Target.DataConnectorRef.ProjectId = projectId
	}
}

func createMicroAggregationBatchJob(cmd *cobra.Command, data []byte) *entities.BatchJobWrapper {
	batchJob := &entities.MicroAggregationBatchJob{}
	err := jsonpb.Unmarshal(bytes.NewReader(data), batchJob)
	if err != nil {
		common.CliExit(err)
	}
	projectId := project.GetProjectId(cmd)
	setMicroAggregationBatchJobProjectIds(batchJob, projectId)
	return &entities.BatchJobWrapper{
		Job: &entities.BatchJobWrapper_MicroAggregationBatchJob{
			MicroAggregationBatchJob: batchJob,
		},
	}
}

func setMicroAggregationBatchJobProjectIds(batchJob *entities.MicroAggregationBatchJob, projectId string) {
	if batchJob.Ref == nil {
		// normal situation where the whole ref attribute in the json is absent.
		batchJob.Ref = &entities.BatchJobRef{}
	}
	batchJob.Ref.ProjectId = projectId
	batchJob.SourceData.DataConnectorRef.ProjectId = projectId
	batchJob.TargetData.DataConnectorRef.ProjectId = projectId
}

func NamesCompletion(cmd *cobra.Command, args []string, complete string) ([]string, cobra.ShellCompDirective) {
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

func toRefWithStates(batchJob *entities.BatchJobWrapper) refWithStates {
	var ref *entities.BatchJobRef
	var states []*entities.BatchJobState
	if encryptionBatchJob := batchJob.GetEncryptionBatchJob(); encryptionBatchJob != nil {
		ref = encryptionBatchJob.Ref
		states = encryptionBatchJob.States
	} else if microAggregationBatchJob := batchJob.GetMicroAggregationBatchJob(); microAggregationBatchJob != nil {
		ref = microAggregationBatchJob.Ref
		states = microAggregationBatchJob.States
	}
	return refWithStates{
		states: states,
		ref:    ref,
	}
}
