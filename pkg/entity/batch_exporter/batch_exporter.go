package batch_exporter

import (
	"context"
	"errors"
	"github.com/golang/protobuf/ptypes/duration"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/strmprivacy/api-definitions-go/v3/api/batch_exporters/v1"
	"github.com/strmprivacy/api-definitions-go/v3/api/data_connectors/v1"
	"github.com/strmprivacy/api-definitions-go/v3/api/entities/v1"
	"google.golang.org/grpc"
	"strings"
	"strmprivacy/strm/pkg/common"
	"strmprivacy/strm/pkg/entity/data_connector"
	"strmprivacy/strm/pkg/entity/project"
	"strmprivacy/strm/pkg/util"
)

var client batch_exporters.BatchExportersServiceClient
var apiContext context.Context

func SetupClient(clientConnection *grpc.ClientConn, ctx context.Context) {
	apiContext = ctx
	client = batch_exporters.NewBatchExportersServiceClient(clientConnection)
}

func list(cmd *cobra.Command) {
	req := &batch_exporters.ListBatchExportersRequest{
		ProjectId: project.GetProjectId(cmd),
	}
	response, err := client.ListBatchExporters(apiContext, req)
	common.CliExit(err)

	printer.Print(response)
}

func get(name *string, cmd *cobra.Command) {
	ref := &entities.BatchExporterRef{
		ProjectId: project.GetProjectId(cmd),
		Name:      *name,
	}
	req := &batch_exporters.GetBatchExporterRequest{Ref: ref}
	response, err := client.GetBatchExporter(apiContext, req)
	common.CliExit(err)
	printer.Print(response)
}

func del(name *string, cmd *cobra.Command) {
	req := &batch_exporters.DeleteBatchExporterRequest{Ref: &entities.BatchExporterRef{
		ProjectId: project.GetProjectId(cmd),
		Name:      *name,
	}}
	response, err := client.DeleteBatchExporter(apiContext, req)
	common.CliExit(err)
	printer.Print(response)
}

func create(streamName *string, cmd *cobra.Command) {
	flags := cmd.Flags()
	keyStream := util.GetBoolAndErr(flags, exportKeys)
	includeExistingEvents := util.GetBoolAndErr(flags, includeExistingEventsFlag)
	dataConnectorName := getDataConnectorName(flags)
	// this exporterName might be empty, in which case the API will set it to
	// the appropriate default
	exporterName := util.GetStringAndErr(flags, nameFlag)
	i := util.GetInt64AndErr(flags, intervalFlag)
	interval := duration.Duration{Seconds: i}

	pathPrefix := util.GetStringAndErr(flags, pathPrefix)
	projectId := project.GetProjectId(cmd)
	exporter := &entities.BatchExporter{
		Ref: &entities.BatchExporterRef{
			ProjectId: projectId,
			Name:      exporterName,
		},
		DataConnectorRef: &entities.DataConnectorRef{
			ProjectId: projectId,
			Name:      dataConnectorName,
		},
		Interval:              &interval,
		PathPrefix:            pathPrefix,
		IncludeExistingEvents: includeExistingEvents,
	}
	if keyStream {
		exporter.StreamOrKeyStreamRef = &entities.BatchExporter_KeyStreamRef{
			KeyStreamRef: &entities.KeyStreamRef{
				ProjectId: projectId,
				Name:      *streamName,
			}}
	} else {
		exporter.StreamOrKeyStreamRef = &entities.BatchExporter_StreamRef{
			StreamRef: &entities.StreamRef{
				ProjectId: projectId,
				Name:      *streamName,
			}}
	}

	response, err := client.CreateBatchExporter(apiContext,
		&batch_exporters.CreateBatchExporterRequest{BatchExporter: exporter})
	common.CliExit(err)
	printer.Print(response)
}

func getDataConnectorName(flags *pflag.FlagSet) string {
	dataConnectorName := util.GetStringAndErr(flags, dataConnectorFlag)
	availableDataConnectorNames := getDataConnectorNames()

	if len(dataConnectorName) == 0 && len(availableDataConnectorNames) == 1 {
		dataConnectorName = availableDataConnectorNames[0]
	}
	dataConnectorExists := false
	for _, name := range availableDataConnectorNames {
		if name == dataConnectorName {
			dataConnectorExists = true
			break
		}
	}
	if len(dataConnectorName) == 0 || !dataConnectorExists {
		common.CliExit(errors.New("You must provide a valid data connector name when creating a batch exporter"))
	}
	return dataConnectorName
}

func getDataConnectorNames() []string {
	req := &data_connectors.ListDataConnectorsRequest{
		ProjectId: common.ProjectId,
	}
	response, err := data_connector.Client.ListDataConnectors(apiContext, req)

	common.CliExit(err)

	names := make([]string, 0, len(response.DataConnectors))
	for _, dataConnector := range response.DataConnectors {
		names = append(names, dataConnector.Ref.Name)
	}

	return names
}

func NamesCompletion(cmd *cobra.Command, args []string, complete string) ([]string, cobra.ShellCompDirective) {
	if len(args) != 0 && strings.Fields(cmd.Short)[0] != "Delete" {
		// this one means you don't get multiple completion suggestions for one stream if it's not a delete call
		return nil, cobra.ShellCompDirectiveNoFileComp
	}
	req := &batch_exporters.ListBatchExportersRequest{
		ProjectId: common.ProjectId,
	}
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
