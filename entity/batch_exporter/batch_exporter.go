package batch_exporter

import (
	"context"
	"github.com/golang/protobuf/ptypes/duration"
	"github.com/spf13/cobra"
	"github.com/streammachineio/api-definitions-go/api/batch_exporters/v1"
	"github.com/streammachineio/api-definitions-go/api/entities/v1"
	"google.golang.org/grpc"
	"log"
	"streammachine.io/strm/entity/sink"
	"streammachine.io/strm/utils"
)

var BillingId string
var client batch_exporters.BatchExportersServiceClient
var apiContext context.Context

func SetupClient(clientConnection *grpc.ClientConn, ctx context.Context) {
	apiContext = ctx
	client = batch_exporters.NewBatchExportersServiceClient(clientConnection)
}

func list() {
	req := &batch_exporters.ListBatchExportersRequest{BillingId: BillingId}
	exporters, err := client.ListBatchExporters(apiContext, req)
	cobra.CheckErr(err)
	utils.Print(exporters)
}

func get(name *string, _ *cobra.Command) {
	ref := &entities.BatchExporterRef{
		Name: *name, BillingId: BillingId,
	}
	req := &batch_exporters.GetBatchExporterRequest{Ref: ref}
	exporter, err := client.GetBatchExporter(apiContext, req)
	cobra.CheckErr(err)
	utils.Print(exporter)
}

func del(name *string) {
	req := &batch_exporters.DeleteBatchExporterRequest{Ref: &entities.BatchExporterRef{
		BillingId: BillingId, Name: *name}}
	exporter, err := client.DeleteBatchExporter(apiContext, req)
	cobra.CheckErr(err)
	utils.Print(exporter)
}

func create(streamName *string, cmd *cobra.Command) {
	flags := cmd.Flags()
	keyStream := utils.GetBoolAndErr(flags, exportKeys)
	sinkName := utils.GetStringAndErr(flags, sinkFlag)
	sinkNames := sink.ExistingNames()
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
			BillingId: BillingId,
		},
		Interval:   &interval,
		PathPrefix: pathPrefix,
	}
	if keyStream {
		exporter.StreamOrKeyStreamRef = &entities.BatchExporter_KeyStreamRef{
			KeyStreamRef: &entities.KeyStreamRef{
				Name: *streamName, BillingId: BillingId}}
	} else {
		exporter.StreamOrKeyStreamRef = &entities.BatchExporter_StreamRef{
			StreamRef: &entities.StreamRef{
				Name: *streamName, BillingId: BillingId}}
	}

	response, err := client.CreateBatchExporter(apiContext,
		&batch_exporters.CreateBatchExporterRequest{BatchExporter: exporter})
	cobra.CheckErr(err)
	utils.Print(response.BatchExporter)
}

func existingNamesCompletion(cmd *cobra.Command, args []string, complete string) ([]string, cobra.ShellCompDirective) {
	if len(args) != 0 {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}
	req := &batch_exporters.ListBatchExportersRequest{BillingId: BillingId}
	response, err := client.ListBatchExporters(apiContext, req)
	cobra.CheckErr(err)
	streamNames := make([]string, 0, len(response.BatchExporters))
	for _, s := range response.BatchExporters {
		streamNames = append(streamNames, s.Ref.Name)
	}
	return streamNames, cobra.ShellCompDirectiveNoFileComp
}
