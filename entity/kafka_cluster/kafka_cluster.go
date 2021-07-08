package kafka_cluster

import (
	"context"
	"github.com/spf13/cobra"
	"github.com/streammachineio/api-definitions-go/api/entities/v1"
	"github.com/streammachineio/api-definitions-go/api/kafka_clusters/v1"
	"google.golang.org/grpc"
	"streammachine.io/strm/utils"
)

// strings used in the cli
const ()

var BillingId string
var client kafka_clusters.KafkaClustersServiceClient
var apiContext context.Context

func ref(n *string) *entities.KafkaClusterRef {
	return &entities.KafkaClusterRef{BillingId: BillingId, Name: *n}
}

func SetupClient(clientConnection *grpc.ClientConn, ctx context.Context) {
	apiContext = ctx
	client = kafka_clusters.NewKafkaClustersServiceClient(clientConnection)
}

func ExistingNamesCompletion(cmd *cobra.Command, args []string, complete string) ([]string, cobra.ShellCompDirective) {
	if len(args) > 0 {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}
	return ExistingNames(), cobra.ShellCompDirectiveNoFileComp
}

func list() {
	req := &kafka_clusters.ListKafkaClustersRequest{BillingId: BillingId}
	sinksList, err := client.ListKafkaClusters(apiContext, req)
	cobra.CheckErr(err)
	utils.Print(sinksList)
}

func get(name *string) {
	cluster := GetCluster(name)
	utils.Print(cluster)
}
func GetCluster(name *string) *entities.KafkaCluster {
	req := &kafka_clusters.GetKafkaClusterRequest{Ref: ref(name)}
	cluster, err := client.GetKafkaCluster(apiContext, req)
	cobra.CheckErr(err)
	return cluster.KafkaCluster
}

func ExistingNames() []string {
	req := &kafka_clusters.ListKafkaClustersRequest{BillingId: BillingId}
	response, err := client.ListKafkaClusters(apiContext, req)
	cobra.CheckErr(err)
	clusters := response.KafkaClusters
	names := make([]string, 0, len(clusters))
	for _, s := range clusters {
		names = append(names, s.Ref.Name)
	}
	return names
}
