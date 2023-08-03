package kafka_user

import (
	"context"
	"github.com/spf13/cobra"
	"github.com/strmprivacy/api-definitions-go/v3/api/entities/v1"
	"github.com/strmprivacy/api-definitions-go/v3/api/kafka_users/v1"
	"google.golang.org/grpc"
	"strings"
	"strmprivacy/strm/pkg/common"
	"strmprivacy/strm/pkg/entity/kafka_exporter"
	"strmprivacy/strm/pkg/entity/project"
	"strmprivacy/strm/pkg/util"
)

var client kafka_users.KafkaUsersServiceClient
var apiContext context.Context

func ref(n *string, cmd *cobra.Command) *entities.KafkaUserRef {
	return &entities.KafkaUserRef{
		ProjectId: project.GetProjectId(cmd),
		Name:      *n,
	}
}

func SetupClient(clientConnection *grpc.ClientConn, ctx context.Context) {
	apiContext = ctx
	client = kafka_users.NewKafkaUsersServiceClient(clientConnection)
}

func list(exporterName *string, cmd *cobra.Command) {
	req := &kafka_users.ListKafkaUsersRequest{
		Ref: &entities.KafkaExporterRef{
			ProjectId: project.GetProjectId(cmd),
			Name:      *exporterName,
		},
	}
	response, err := client.ListKafkaUsers(apiContext, req)
	common.CliExit(err)
	printer.Print(response)
}

func get(name *string, cmd *cobra.Command) {
	// TODO need api recursive addition
	req := &kafka_users.GetKafkaUserRequest{Ref: ref(name, cmd)}
	response, err := client.GetKafkaUser(apiContext, req)
	common.CliExit(err)
	printer.Print(response)
}

func del(name *string, cmd *cobra.Command) {
	response := ref(name, cmd)
	req := &kafka_users.DeleteKafkaUserRequest{Ref: response}
	_, err := client.DeleteKafkaUser(apiContext, req)
	common.CliExit(err)
	util.DeleteSaved(response, &response.Name)
	printer.Print(response)
}

func create(kafkaExporterName *string, cmd *cobra.Command) {
	flags := cmd.Flags()
	exporter := kafka_exporter.Get(kafkaExporterName, cmd).KafkaExporter

	kafkaUser := &entities.KafkaUser{
		Ref: &entities.KafkaUserRef{
			ProjectId: project.GetProjectId(cmd),
		},
		KafkaExporterName: exporter.Ref.Name,
	}

	response, err := client.CreateKafkaUser(apiContext,
		&kafka_users.CreateKafkaUserRequest{KafkaUser: kafkaUser})
	common.CliExit(err)

	save, err := flags.GetBool(saveFlag)
	if save {
		util.Save(response.KafkaUser, &response.KafkaUser.Ref.Name)
	}

	printer.Print(response)
}

func namesCompletion(cmd *cobra.Command, args []string, complete string) ([]string, cobra.ShellCompDirective) {
	if len(args) != 0 && strings.Fields(cmd.Short)[0] != "Delete" {
		// this one means you don't get multiple completion suggestions for one stream if it's not a delete call
		return nil, cobra.ShellCompDirectiveNoFileComp
	}

	req := &kafka_users.ListKafkaUsersRequest{
		Ref: &entities.KafkaExporterRef{
			ProjectId: common.ProjectId,
		},
	}
	response, err := client.ListKafkaUsers(apiContext, req)

	if err != nil {
		return common.GrpcRequestCompletionError(err)
	}

	names := make([]string, 0, len(response.KafkaUsers))
	for _, s := range response.KafkaUsers {
		names = append(names, s.Ref.Name)
	}
	return names, cobra.ShellCompDirectiveNoFileComp
}
