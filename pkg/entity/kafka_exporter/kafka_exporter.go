package kafka_exporter

import (
	"context"
	"github.com/spf13/cobra"
	"github.com/strmprivacy/api-definitions-go/v3/api/entities/v1"
	"github.com/strmprivacy/api-definitions-go/v3/api/kafka_exporters/v1"
	"google.golang.org/grpc"
	"strings"
	"strmprivacy/strm/pkg/common"
	"strmprivacy/strm/pkg/entity/project"
	"strmprivacy/strm/pkg/util"
)

var client kafka_exporters.KafkaExportersServiceClient
var apiContext context.Context

func SetupClient(clientConnection *grpc.ClientConn, ctx context.Context) {
	apiContext = ctx
	client = kafka_exporters.NewKafkaExportersServiceClient(clientConnection)
}

func Get(name *string, cmd *cobra.Command) *kafka_exporters.GetKafkaExporterResponse {
	req := &kafka_exporters.GetKafkaExporterRequest{
		Ref: &entities.KafkaExporterRef{
			ProjectId: project.GetProjectId(cmd),
			Name:      *name,
		},
	}
	exporter, err := client.GetKafkaExporter(apiContext, req)
	common.CliExit(err)
	return exporter
}

func list(recursive bool, cmd *cobra.Command) {
	// TODO need api recursive addition
	req := &kafka_exporters.ListKafkaExportersRequest{
		ProjectId: project.GetProjectId(cmd),
	}
	response, err := client.ListKafkaExporters(apiContext, req)
	common.CliExit(err)
	printer.Print(response)
}

func get(name *string, cmd *cobra.Command, recursive bool) {
	response := Get(name, cmd)
	printer.Print(response)
}

func del(name *string, recursive bool, cmd *cobra.Command) {
	exporterRef := &entities.KafkaExporterRef{
		ProjectId: project.GetProjectId(cmd),
		Name:      *name,
	}
	exporter := Get(name, cmd)

	req := &kafka_exporters.DeleteKafkaExporterRequest{Ref: exporterRef, Recursive: recursive}
	response, err := client.DeleteKafkaExporter(apiContext, req)
	common.CliExit(err)

	for _, user := range exporter.KafkaExporter.Users {
		util.DeleteSaved(user, &user.Ref.Name)
	}

	printer.Print(response)
}

func create(name *string, cmd *cobra.Command) {
	flags := cmd.Flags()
	_, err := flags.GetString(clusterFlag) // TODO at the moment, the cluster flag is ignored
	common.CliExit(err)

	projectId := project.GetProjectId(cmd)
	// key streams not yet supported in data model!
	exporter := &entities.KafkaExporter{
		StreamRef: &entities.StreamRef{
			ProjectId: projectId,
			Name:      *name,
		},
		Ref: &entities.KafkaExporterRef{
			ProjectId: projectId,
		},
	}

	response, err := client.CreateKafkaExporter(
		apiContext,
		&kafka_exporters.CreateKafkaExporterRequest{KafkaExporter: exporter},
	)

	common.CliExit(err)

	save, err := flags.GetBool(saveFlag)
	if save {
		user := response.KafkaExporter.Users[0]
		util.Save(user, &user.Ref.Name)
	}

	printer.Print(response)
}

func NamesCompletion(cmd *cobra.Command, args []string, complete string) ([]string, cobra.ShellCompDirective) {
	if len(args) != 0 && strings.Fields(cmd.Short)[0] != "Delete" {
		// this one means you don't get multiple completion suggestions for one stream if it's not a delete call
		return nil, cobra.ShellCompDirectiveNoFileComp
	}

	req := &kafka_exporters.ListKafkaExportersRequest{
		ProjectId: common.ProjectId,
	}
	response, err := client.ListKafkaExporters(apiContext, req)

	if err != nil {
		return common.GrpcRequestCompletionError(err)
	}

	names := make([]string, 0, len(response.KafkaExporters))
	for _, s := range response.KafkaExporters {
		names = append(names, s.Ref.Name)
	}

	return names, cobra.ShellCompDirectiveNoFileComp
}
