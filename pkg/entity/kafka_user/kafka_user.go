package kafka_user

import (
	"context"
	"github.com/spf13/cobra"
	"github.com/streammachineio/api-definitions-go/api/entities/v1"
	"github.com/streammachineio/api-definitions-go/api/kafka_users/v1"
	"google.golang.org/grpc"
	"streammachine.io/strm/pkg/common"
	"streammachine.io/strm/pkg/entity/kafka_exporter"
	"streammachine.io/strm/pkg/utils"
)

var client kafka_users.KafkaUsersServiceClient
var apiContext context.Context

func ref(n *string) *entities.KafkaUserRef {
	return &entities.KafkaUserRef{BillingId: common.BillingId, Name: *n}
}

func SetupClient(clientConnection *grpc.ClientConn, ctx context.Context) {
	apiContext = ctx
	client = kafka_users.NewKafkaUsersServiceClient(clientConnection)
}

func list(exporterName *string) {
	req := &kafka_users.ListKafkaUsersRequest{
		Ref: &entities.KafkaExporterRef{
			BillingId: common.BillingId,
			Name:      *exporterName,
		},
	}
	users, err := client.ListKafkaUsers(apiContext, req)
	common.CliExit(err)
	utils.Print(users)
}

func get(name *string) {
	// TODO need api recursive addition
	req := &kafka_users.GetKafkaUserRequest{Ref: ref(name)}
	user, err := client.GetKafkaUser(apiContext, req)
	common.CliExit(err)
	utils.Print(user)
}

func del(name *string) {
	user := ref(name)
	req := &kafka_users.DeleteKafkaUserRequest{Ref: user}
	_, err := client.DeleteKafkaUser(apiContext, req)
	common.CliExit(err)
	utils.Print(user)
	utils.DeleteSaved(user, &user.Name)
}

func create(kafkaExporterName *string, cmd *cobra.Command) {

	flags := cmd.Flags()
	exporter := kafka_exporter.Get(kafkaExporterName).KafkaExporter
	kafkaUser := &entities.KafkaUser{
		Ref: &entities.KafkaUserRef{
			BillingId: common.BillingId,
		},
		KafkaExporterName: exporter.Ref.Name,
	}

	response, err := client.CreateKafkaUser(apiContext,
		&kafka_users.CreateKafkaUserRequest{KafkaUser: kafkaUser})
	common.CliExit(err)
	utils.Print(response.KafkaUser)
	save, err := flags.GetBool(saveFlag)
	if save {
		utils.Save(response.KafkaUser, &response.KafkaUser.Ref.Name)
	}

}

func namesCompletion(cmd *cobra.Command, args []string, complete string) ([]string, cobra.ShellCompDirective) {
	if len(args) > 0 || common.BillingIdIsMissing() {
		return common.MissingBillingIdCompletionError(cmd.CommandPath())
	}

	req := &kafka_users.ListKafkaUsersRequest{
		Ref: &entities.KafkaExporterRef{BillingId: common.BillingId},
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