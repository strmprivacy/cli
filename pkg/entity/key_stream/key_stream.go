package key_stream

import (
	"context"
	"github.com/spf13/cobra"
	"github.com/strmprivacy/api-definitions-go/v2/api/entities/v1"
	"github.com/strmprivacy/api-definitions-go/v2/api/key_streams/v1"
	"google.golang.org/grpc"
	"strmprivacy/strm/pkg/common"
)

var client key_streams.KeyStreamsServiceClient
var apiContext context.Context

func SetupClient(clientConnection *grpc.ClientConn, ctx context.Context) {
	apiContext = ctx
	client = key_streams.NewKeyStreamsServiceClient(clientConnection)
}

func list() {
	req := &key_streams.ListKeyStreamsRequest{
		ProjectId: common.ProjectId,
	}
	response, err := client.ListKeyStreams(apiContext, req)
	common.CliExit(err)
	printer.Print(response)
}

func get(name *string) {
	req := &key_streams.GetKeyStreamRequest{Ref: ref(name)}
	response, err := client.GetKeyStream(apiContext, req)
	common.CliExit(err)
	printer.Print(response)
}

func ref(n *string) *entities.KeyStreamRef {
	return &entities.KeyStreamRef{
		ProjectId: common.ProjectId,
		Name: *n,
	}
}

func NamesCompletion(cmd *cobra.Command, args []string, complete string) ([]string, cobra.ShellCompDirective) {
	if len(args) != 0 {
		// this one means you don't get multiple completion suggestions for one stream
		return nil, cobra.ShellCompDirectiveNoFileComp
	}

	req := &key_streams.ListKeyStreamsRequest{
		ProjectId: common.ProjectId,
	}
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
