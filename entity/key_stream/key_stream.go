package key_stream

import (
	"context"
	"github.com/spf13/cobra"
	"github.com/streammachineio/api-definitions-go/api/entities/v1"
	"github.com/streammachineio/api-definitions-go/api/key_streams/v1"
	"google.golang.org/grpc"
	"streammachine.io/strm/utils"
)

var BillingId string
var client key_streams.KeyStreamsServiceClient
var apiContext context.Context

func SetupClient(clientConnection *grpc.ClientConn, ctx context.Context) {
	apiContext = ctx
	client = key_streams.NewKeyStreamsServiceClient(clientConnection)
}

func ExistingNamesCompletion(cmd *cobra.Command, args []string, complete string) ([]string, cobra.ShellCompDirective) {
	return ExistingNames(), cobra.ShellCompDirectiveNoFileComp
}

func ExistingNames() []string {

	req := &key_streams.ListKeyStreamsRequest{BillingId: BillingId}
	response, err := client.ListKeyStreams(apiContext, req)
	cobra.CheckErr(err)
	names := make([]string, 0, len(response.KeyStreams))
	for _, s := range response.KeyStreams {
		names = append(names, s.Ref.Name)
	}
	return names
}

func list() {
	req := &key_streams.ListKeyStreamsRequest{BillingId: BillingId}
	sinksList, err := client.ListKeyStreams(apiContext, req)
	cobra.CheckErr(err)
	utils.Print(sinksList)
}

func get(name *string) {
	req := &key_streams.GetKeyStreamRequest{Ref: ref(name)}
	stream, err := client.GetKeyStream(apiContext, req)
	cobra.CheckErr(err)
	utils.Print(stream)
}

func ref(n *string) *entities.KeyStreamRef {
	return &entities.KeyStreamRef{BillingId: BillingId, Name: *n}
}
