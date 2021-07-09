package kafka_exporter

import (
	"context"
	"github.com/spf13/cobra"
	"github.com/streammachineio/api-definitions-go/api/entities/v1"
	"github.com/streammachineio/api-definitions-go/api/kafka_exporters/v1"
	"google.golang.org/grpc"
	"streammachine.io/strm/entity/kafka_cluster"
	"streammachine.io/strm/entity/stream"
	"streammachine.io/strm/utils"
)

var BillingId string
var client kafka_exporters.KafkaExportersServiceClient
var apiContext context.Context

func SetupClient(clientConnection *grpc.ClientConn, ctx context.Context) {
	apiContext = ctx
	client = kafka_exporters.NewKafkaExportersServiceClient(clientConnection)
}

func Get(name *string) *kafka_exporters.GetKafkaExporterResponse {
	req := &kafka_exporters.GetKafkaExporterRequest{Ref: ref(name)}
	exporter, err := client.GetKafkaExporter(apiContext, req)
	cobra.CheckErr(err)
	return exporter
}

func list(recursive bool) {
	// TODO need api recursive addition
	req := &kafka_exporters.ListKafkaExportersRequest{BillingId: BillingId}
	exporters, err := client.ListKafkaExporters(apiContext, req)
	cobra.CheckErr(err)
	utils.Print(exporters)
}

func get(name *string, recursive bool) {
	exporter := Get(name)
	utils.Print(exporter)
}

func del(name *string, recursive bool) {
	exporterRef := ref(name)
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
	clusterName, err := flags.GetString(clusterFlag)
	cobra.CheckErr(err)
	clusters := kafka_cluster.ExistingNames()
	if len(clusterName) == 0 && len(clusters) == 1 {
		clusterName = clusters[0]
	}
	cluster := kafka_cluster.GetCluster(&clusterName)

	// key streams not yet supported in data model!
	exporter := &entities.KafkaExporter{
		Ref:       ref(name),
		StreamRef: stream.Ref(name),
		Target:    &entities.KafkaExporterTarget{ClusterRef: cluster.Ref},
	}

	response, err := client.CreateKafkaExporter(apiContext,
		&kafka_exporters.CreateKafkaExporterRequest{KafkaExporter: exporter})

	cobra.CheckErr(err)
	utils.Print(response.KafkaExporter)

	save, err := flags.GetBool(saveFlag)
	if save {
		user := response.KafkaExporter.Users[0]
		utils.Save(user, &user.Ref.Name)
	}

}

func ExistingNamesCompletion(cmd *cobra.Command, args []string, complete string) ([]string, cobra.ShellCompDirective) {
	if len(args) > 0 {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}
	return ExistingNames(), cobra.ShellCompDirectiveNoFileComp
}

func ExistingNames() []string {

	req := &kafka_exporters.ListKafkaExportersRequest{BillingId: BillingId}
	response, err := client.ListKafkaExporters(apiContext, req)
	cobra.CheckErr(err)
	streamNames := make([]string, 0, len(response.KafkaExporters))
	for _, s := range response.KafkaExporters {
		streamNames = append(streamNames, s.Ref.Name)
	}
	return streamNames
}

func ref(n *string) *entities.KafkaExporterRef {
	return &entities.KafkaExporterRef{BillingId: BillingId, Name: *n}
}
