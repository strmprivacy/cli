package batch_exporter

import (
	"context"
	"github.com/golang/protobuf/ptypes/duration"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/streammachineio/api-definitions-go/api/batch_exporters/v1"
	"github.com/streammachineio/api-definitions-go/api/entities/v1"
	"github.com/streammachineio/api-definitions-go/api/sinks/v1"
	"google.golang.org/grpc"
	"streammachine.io/strm/common"
	"streammachine.io/strm/entity/sink"
	"streammachine.io/strm/utils"
)

var client batch_exporters.BatchExportersServiceClient
var apiContext context.Context

func SetupClient(clientConnection *grpc.ClientConn, ctx context.Context) {
	apiContext = ctx
	client = batch_exporters.NewBatchExportersServiceClient(clientConnection)
}

func list() {
	req := &batch_exporters.ListBatchExportersRequest{BillingId: common.BillingId}
	exporters, err := client.ListBatchExporters(apiContext, req)
	common.CliExit(err)
	utils.Print(exporters)
}

func get(name *string, _ *cobra.Command) {
	ref := &entities.BatchExporterRef{
		Name: *name, BillingId: common.BillingId,
	}
	req := &batch_exporters.GetBatchExporterRequest{Ref: ref}
	exporter, err := client.GetBatchExporter(apiContext, req)
	common.CliExit(err)
	utils.Print(exporter)
}

func del(name *string) {
	req := &batch_exporters.DeleteBatchExporterRequest{Ref: &entities.BatchExporterRef{
		BillingId: common.BillingId, Name: *name}}
	exporter, err := client.DeleteBatchExporter(apiContext, req)
	common.CliExit(err)
	utils.Print(exporter)
}

func create(streamName *string, cmd *cobra.Command) {
	flags := cmd.Flags()
	keyStream := utils.GetBoolAndErr(flags, exportKeys)
	sinkName := utils.GetStringAndErr(flags, sinkFlag)
	sinkNames := getSinkNames()
	if len(sinkName) == 0 && len(sinkNames) == 1 {
		sinkName = sinkNames[0]
	}
	if len(sinkName) == 0 {
		log.Fatalf("You must provide a sink name when creating a batch exporter")
	}
	// this exporterName might be empty, in which case the API will set it to
	// the appropriate default
	exporterName := utils.GetStringAndErr(flags, nameFlag)
	i := utils.GetInt64AndErr(flags, intervalFlag)
	interval := duration.Duration{Seconds: i}

	pathPrefix := utils.GetStringAndErr(flags, pathPrefix)

	exporter := &entities.BatchExporter{
		SinkName: sinkName,
		Ref: &entities.BatchExporterRef{
			Name:      exporterName,
			BillingId: common.BillingId,
		},
		Interval:   &interval,
		PathPrefix: pathPrefix,
	}
	if keyStream {
		exporter.StreamOrKeyStreamRef = &entities.BatchExporter_KeyStreamRef{
			KeyStreamRef: &entities.KeyStreamRef{
				Name: *streamName, BillingId: common.BillingId}}
	} else {
		exporter.StreamOrKeyStreamRef = &entities.BatchExporter_StreamRef{
			StreamRef: &entities.StreamRef{
				Name: *streamName, BillingId: common.BillingId}}
	}

	response, err := client.CreateBatchExporter(apiContext,
		&batch_exporters.CreateBatchExporterRequest{BatchExporter: exporter})
	common.CliExit(err)
	utils.Print(response.BatchExporter)
}

func getSinkNames() []string {
	req := &sinks.ListSinksRequest{BillingId: common.BillingId}
	response, err := sink.Client.ListSinks(apiContext, req)

	common.CliExit(err)

	names := make([]string, 0, len(response.Sinks))
	for _, s := range response.Sinks {
		names = append(names, s.Sink.Ref.Name)
	}

	return names
}

func namesCompletion(cmd *cobra.Command, args []string, complete string) ([]string, cobra.ShellCompDirective) {
	if len(args) != 0 || common.BillingIdIsMissing() {
		return common.MissingBillingIdCompletionError(cmd.CommandPath())
	}
	req := &batch_exporters.ListBatchExportersRequest{BillingId: common.BillingId}
	response, err := client.ListBatchExporters(apiContext, req)

	if err != nil {
		return common.GrpcRequestCompletionError(err)
	}

	streamNames := make([]string, 0, len(response.BatchExporters))
	for _, s := range response.BatchExporters {
		streamNames = append(streamNames, s.Ref.Name)
	}
	return streamNames, cobra.ShellCompDirectiveNoFileComp
}
