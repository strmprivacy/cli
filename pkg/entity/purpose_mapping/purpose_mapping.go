package purpose_mapping

import (
	"context"
	"github.com/spf13/cobra"
	"github.com/strmprivacy/api-definitions-go/v2/api/entities/v1"
	"github.com/strmprivacy/api-definitions-go/v2/api/purpose_mapping/v1"
	"google.golang.org/grpc"
	"sort"
	"strconv"
	"strmprivacy/strm/pkg/common"
)

var client purpose_mapping.PurposeMappingServiceClient
var apiContext context.Context

func SetupClient(clientConnection *grpc.ClientConn, ctx context.Context) {
	apiContext = ctx
	client = purpose_mapping.NewPurposeMappingServiceClient(clientConnection)
}

func get(level int32) {
	req := &purpose_mapping.GetPurposeMappingRequest{
		Level: level,
	}

	response, err := client.GetPurposeMapping(apiContext, req)
	common.CliExit(err)
	printer.Print(response)
}

func list() (*purpose_mapping.ListPurposeMappingsResponse, error) {
	req := &purpose_mapping.ListPurposeMappingsRequest{}
	response, err := client.ListPurposeMappings(apiContext, req)

	sort.Slice(response.PurposeMappings, func(i, j int) bool {
		return response.PurposeMappings[i].Level < response.PurposeMappings[j].Level
	})

	return response, err
}

func listAndPrint() {
	response, err := list()
	common.CliExit(err)
	printer.Print(response)
}

func create(newPurposeMapping string) {
	req := &purpose_mapping.CreatePurposeMappingRequest{
		PurposeMapping: &entities.PurposeMapping{
			Purpose: newPurposeMapping,
		},
	}

	response, err := client.CreatePurposeMapping(apiContext, req)
	common.CliExit(err)
	printer.Print(response)
}

func NamesCompletion(cmd *cobra.Command, args []string, complete string) ([]string, cobra.ShellCompDirective) {
	response, err := list()

	if err != nil {
		return common.GrpcRequestCompletionError(err)
	}

	names := make([]string, 0, len(response.PurposeMappings))
	for _, purposeMapping := range response.PurposeMappings {
		names = append(names, strconv.Itoa(int(purposeMapping.Level)))
	}

	return names, cobra.ShellCompDirectiveNoFileComp
}
