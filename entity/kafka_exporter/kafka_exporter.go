package kafka_exporter

import (
	"context"
	"github.com/spf13/cobra"
	"github.com/streammachineio/api-definitions-go/api/entities/v1"
	"github.com/streammachineio/api-definitions-go/api/kafka_exporters/v1"
	"google.golang.org/grpc"
	"streammachine.io/strm/common"
	"streammachine.io/strm/utils"
)

var client kafka_exporters.KafkaExportersServiceClient
var apiContext context.Context

func SetupClient(clientConnection *grpc.ClientConn, ctx context.Context) {
	apiContext = ctx
	client = kafka_exporters.NewKafkaExportersServiceClient(clientConnection)
}

func Get(name *string) *kafka_exporters.GetKafkaExporterResponse {
	req := &kafka_exporters.GetKafkaExporterRequest{
		Ref: &entities.KafkaExporterRef{BillingId: common.BillingId, Name: *name},
	}
	exporter, err := client.GetKafkaExporter(apiContext, req)
	cobra.CheckErr(err)
	return exporter
}

func list(recursive bool) {
	// TODO need api recursive addition
	req := &kafka_exporters.ListKafkaExportersRequest{BillingId: common.BillingId}
	exporters, err := client.ListKafkaExporters(apiContext, req)
	cobra.CheckErr(err)
	utils.Print(exporters)
}

func get(name *string, recursive bool) {
	exporter := Get(name)
	utils.Print(exporter)
}

func del(name *string, recursive bool) {
	exporterRef := &entities.KafkaExporterRef{BillingId: common.BillingId, Name: *name}
	exporter := Get(name)

	req := &kafka_exporters.DeleteKafkaExporterRequest{Ref: exporterRef, Recursive: recursive}
	_, err := client.DeleteKafkaExporter(apiContext, req)
	cobra.CheckErr(err)
	for _, user := range exporter.KafkaExporter.Users {
		utils.DeleteSaved(user, &user.Ref.Name)
	}
}

func create(name *string, cmd *cobra.Command) {
	flags := cmd.Flags()
	_, err := flags.GetString(clusterFlag) // TODO at the moment, the cluster flag is ignored
	cobra.CheckErr(err)

	// key streams not yet supported in data model!
	exporter := &entities.KafkaExporter{
		StreamRef: &entities.StreamRef{BillingId: common.BillingId, Name: *name},
	}

	response, err := client.CreateKafkaExporter(
		apiContext,
		&kafka_exporters.CreateKafkaExporterRequest{KafkaExporter: exporter},
	)

	cobra.CheckErr(err)
	utils.Print(response.KafkaExporter)

	save, err := flags.GetBool(saveFlag)
	if save {
		user := response.KafkaExporter.Users[0]
		utils.Save(user, &user.Ref.Name)
	}

}

func KafkaExporterNamesCompletion(cmd *cobra.Command, args []string, complete string) ([]string, cobra.ShellCompDirective) {
	if len(args) > 0 || common.BillingIdIsMissing() {
		return common.MissingBillingIdCompletionError(cmd.CommandPath())
	}

	req := &kafka_exporters.ListKafkaExportersRequest{BillingId: common.BillingId}
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
