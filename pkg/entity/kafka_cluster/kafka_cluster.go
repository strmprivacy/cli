package kafka_cluster

import (
	"context"
	"github.com/spf13/cobra"
	"github.com/streammachineio/api-definitions-go/api/entities/v1"
	"github.com/streammachineio/api-definitions-go/api/kafka_clusters/v1"
	"google.golang.org/grpc"
	"streammachine.io/strm/pkg/common"
)

// strings used in the cli
const ()

var client kafka_clusters.KafkaClustersServiceClient
var apiContext context.Context

func ref(n *string) *entities.KafkaClusterRef {
	return &entities.KafkaClusterRef{BillingId: common.BillingId, Name: *n}
}

func SetupClient(clientConnection *grpc.ClientConn, ctx context.Context) {
	apiContext = ctx
	client = kafka_clusters.NewKafkaClustersServiceClient(clientConnection)
}

func list() {
	req := &kafka_clusters.ListKafkaClustersRequest{BillingId: common.BillingId}
	sinksList, err := client.ListKafkaClusters(apiContext, req)
	common.CliExit(err)
	printer.Print(sinksList)
}

func get(name *string) {
	cluster := GetCluster(name)
	printer.Print(cluster)
}

func GetCluster(name *string) *entities.KafkaCluster {
	req := &kafka_clusters.GetKafkaClusterRequest{Ref: ref(name)}
	cluster, err := client.GetKafkaCluster(apiContext, req)
	common.CliExit(err)
	return cluster.KafkaCluster
}

func NamesCompletion(cmd *cobra.Command, args []string, complete string) ([]string, cobra.ShellCompDirective) {
	if common.BillingIdIsMissing() {
		return common.MissingBillingIdCompletionError(cmd.CommandPath())
	}
	if len(args) != 0 {
		// this one means you don't get two completion suggestions for one stream
		return nil, cobra.ShellCompDirectiveNoFileComp
	}

	req := &kafka_clusters.ListKafkaClustersRequest{BillingId: common.BillingId}
	response, err := client.ListKafkaClusters(apiContext, req)

	if err != nil {
		return common.GrpcRequestCompletionError(err)
	}

	names := make([]string, 0, len(response.KafkaClusters))
	for _, s := range response.KafkaClusters {
		names = append(names, s.Ref.Name)
	}

	return names, cobra.ShellCompDirectiveNoFileComp
}
