package kafka_cluster

import (
	"context"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/strmprivacy/api-definitions-go/v2/api/entities/v1"
	"github.com/strmprivacy/api-definitions-go/v2/api/kafka_clusters/v1"
	"google.golang.org/grpc"
	"strmprivacy/strm/pkg/auth"
	"strmprivacy/strm/pkg/common"
)

// strings used in the cli
const ()

var client kafka_clusters.KafkaClustersServiceClient
var apiContext context.Context

func ref(n *string) *entities.KafkaClusterRef {
	return &entities.KafkaClusterRef{BillingId: auth.Auth.BillingId(), Name: *n}
}

func SetupClient(clientConnection *grpc.ClientConn, ctx context.Context) {
	apiContext = ctx
	client = kafka_clusters.NewKafkaClustersServiceClient(clientConnection)
}

func list() {
	req := &kafka_clusters.ListKafkaClustersRequest{BillingId: auth.Auth.BillingId()}
	response, err := client.ListKafkaClusters(apiContext, req)
	common.CliExit(err)
	printer.Print(response)
}

func get(name *string) {
	req := &kafka_clusters.GetKafkaClusterRequest{Ref: ref(name)}
	response, err := client.GetKafkaCluster(apiContext, req)
	common.CliExit(err)

	printer.Print(response)
}

func NamesCompletion(cmd *cobra.Command, args []string, complete string) ([]string, cobra.ShellCompDirective) {
	if auth.Auth.BillingIdAbsent() {
		return common.MissingBillingIdCompletionError(cmd.CommandPath())
	}
	if len(args) != 0 {
		// this one means you don't get two completion suggestions for one stream
		return nil, cobra.ShellCompDirectiveNoFileComp
	}

	req := &kafka_clusters.ListKafkaClustersRequest{BillingId: auth.Auth.BillingId()}
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

func RefToString(clusterRef *entities.KafkaClusterRef) string {
	return fmt.Sprintf("%v/%v", clusterRef.BillingId, clusterRef.Name)
}
