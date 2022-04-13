package batch_exporter

import (
	"context"
	"errors"
	"strings"

	"github.com/golang/protobuf/ptypes/duration"
	"github.com/spf13/cobra"
	"github.com/strmprivacy/api-definitions-go/v2/api/batch_exporters/v1"
	"github.com/strmprivacy/api-definitions-go/v2/api/entities/v1"
	"github.com/strmprivacy/api-definitions-go/v2/api/sinks/v1"
	"google.golang.org/grpc"
	"strmprivacy/strm/pkg/auth"
	"strmprivacy/strm/pkg/common"
	"strmprivacy/strm/pkg/entity/sink"
	"strmprivacy/strm/pkg/util"
)

var client batch_exporters.BatchExportersServiceClient
var apiContext context.Context

func SetupClient(clientConnection *grpc.ClientConn, ctx context.Context) {
	apiContext = ctx
	client = batch_exporters.NewBatchExportersServiceClient(clientConnection)
}

func list() {
	req := &batch_exporters.ListBatchExportersRequest{BillingId: auth.Auth.BillingId()}
	response, err := client.ListBatchExporters(apiContext, req)
	common.CliExit(err)

	printer.Print(response)
}

func get(name *string, _ *cobra.Command) {
	ref := &entities.BatchExporterRef{
		Name: *name, BillingId: auth.Auth.BillingId(),
	}
	req := &batch_exporters.GetBatchExporterRequest{Ref: ref}
	response, err := client.GetBatchExporter(apiContext, req)
	common.CliExit(err)
	printer.Print(response)
}

func del(name *string) {
	req := &batch_exporters.DeleteBatchExporterRequest{Ref: &entities.BatchExporterRef{
		BillingId: auth.Auth.BillingId(), Name: *name}}
	response, err := client.DeleteBatchExporter(apiContext, req)
	common.CliExit(err)
	printer.Print(response)
}

func create(streamName *string, cmd *cobra.Command) {
	flags := cmd.Flags()
	keyStream := util.GetBoolAndErr(flags, exportKeys)
	includeExistingEvents := util.GetBoolAndErr(flags, includeExistingEventsFlag)
	sinkName := util.GetStringAndErr(flags, sinkFlag)
	sinkNames := getSinkNames()
	if len(sinkName) == 0 && len(sinkNames) == 1 {
		sinkName = sinkNames[0]
	}
	if len(sinkName) == 0 {
		common.CliExit(errors.New("You must provide a sink name when creating a batch exporter"))
	}
	// this exporterName might be empty, in which case the API will set it to
	// the appropriate default
	exporterName := util.GetStringAndErr(flags, nameFlag)
	i := util.GetInt64AndErr(flags, intervalFlag)
	interval := duration.Duration{Seconds: i}

	pathPrefix := util.GetStringAndErr(flags, pathPrefix)

	exporter := &entities.BatchExporter{
		SinkName: sinkName,
		Ref: &entities.BatchExporterRef{
			Name:      exporterName,
			BillingId: auth.Auth.BillingId(),
		},
		Interval:              &interval,
		PathPrefix:            pathPrefix,
		IncludeExistingEvents: includeExistingEvents,
	}
	if keyStream {
		exporter.StreamOrKeyStreamRef = &entities.BatchExporter_KeyStreamRef{
			KeyStreamRef: &entities.KeyStreamRef{
				Name: *streamName, BillingId: auth.Auth.BillingId()}}
	} else {
		exporter.StreamOrKeyStreamRef = &entities.BatchExporter_StreamRef{
			StreamRef: &entities.StreamRef{
				Name: *streamName, BillingId: auth.Auth.BillingId()}}
	}

	response, err := client.CreateBatchExporter(apiContext,
		&batch_exporters.CreateBatchExporterRequest{BatchExporter: exporter})
	common.CliExit(err)
	printer.Print(response)
}

func getSinkNames() []string {
	req := &sinks.ListSinksRequest{BillingId: auth.Auth.BillingId()}
	response, err := sink.Client.ListSinks(apiContext, req)

	common.CliExit(err)

	names := make([]string, 0, len(response.Sinks))
	for _, s := range response.Sinks {
		names = append(names, s.Sink.Ref.Name)
	}

	return names
}

func namesCompletion(cmd *cobra.Command, args []string, complete string) ([]string, cobra.ShellCompDirective) {
	if len(args) != 0 && strings.Fields(cmd.Short)[0] != "Delete" {
		// this one means you don't get multiple completion suggestions for one stream if it's not a delete call
		return nil, cobra.ShellCompDirectiveNoFileComp
	}
	if auth.Auth.BillingIdAbsent() {
		return common.MissingBillingIdCompletionError(cmd.CommandPath())
	}
	req := &batch_exporters.ListBatchExportersRequest{BillingId: auth.Auth.BillingId()}
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
