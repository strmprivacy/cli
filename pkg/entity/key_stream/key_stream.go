package key_stream

import (
	"context"
	"github.com/spf13/cobra"
	"github.com/streammachineio/api-definitions-go/api/entities/v1"
	"github.com/streammachineio/api-definitions-go/api/key_streams/v1"
	"google.golang.org/grpc"
	"streammachine.io/strm/pkg/common"
)

var client key_streams.KeyStreamsServiceClient
var apiContext context.Context

func SetupClient(clientConnection *grpc.ClientConn, ctx context.Context) {
	apiContext = ctx
	client = key_streams.NewKeyStreamsServiceClient(clientConnection)
}

func list() {
	req := &key_streams.ListKeyStreamsRequest{BillingId: common.BillingId}
	keyStreams, err := client.ListKeyStreams(apiContext, req)
	common.CliExit(err)
	printer.Print(keyStreams)
}

func get(name *string) {
	req := &key_streams.GetKeyStreamRequest{Ref: ref(name)}
	stream, err := client.GetKeyStream(apiContext, req)
	common.CliExit(err)
	printer.Print(stream)
}

func ref(n *string) *entities.KeyStreamRef {
	return &entities.KeyStreamRef{BillingId: common.BillingId, Name: *n}
}

func NamesCompletion(cmd *cobra.Command, args []string, complete string) ([]string, cobra.ShellCompDirective) {
	if common.BillingIdIsMissing() {
		return common.MissingBillingIdCompletionError(cmd.CommandPath())
	}
	if len(args) != 0 {
		// this one means you don't get two completion suggestions for one stream
		return nil, cobra.ShellCompDirectiveNoFileComp
	}

	req := &key_streams.ListKeyStreamsRequest{BillingId: common.BillingId}
	response, err := client.ListKeyStreams(apiContext, req)

	if err != nil {
		return common.GrpcRequestCompletionError(err)
	}

	names := make([]string, 0, len(response.KeyStreams))
	for _, s := range response.KeyStreams {
		names = append(names, s.Ref.Name)
	}

	return names, cobra.ShellCompDirectiveNoFileComp
}
