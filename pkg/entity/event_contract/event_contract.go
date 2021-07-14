package event_contract

import (
	"context"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/streammachineio/api-definitions-go/api/entities/v1"
	"github.com/streammachineio/api-definitions-go/api/event_contracts/v1"
	"google.golang.org/grpc"
	"streammachine.io/strm/pkg/common"
	"streammachine.io/strm/pkg/util"
	"strings"
)

// strings used in the cli
const ()

var client event_contracts.EventContractsServiceClient
var apiContext context.Context

func ref(n *string) *entities.EventContractRef {
	parts := strings.Split(*n, "/")
	return &entities.EventContractRef{
		Handle:  parts[0],
		Name:    parts[1],
		Version: parts[2],
	}
}

func refToString(ref *entities.EventContractRef) string {
	return fmt.Sprintf("%v/%v/%v", ref.Handle, ref.Name, ref.Version)
}

func SetupClient(clientConnection *grpc.ClientConn, ctx context.Context) {
	apiContext = ctx
	client = event_contracts.NewEventContractsServiceClient(clientConnection)
}

func list() {
	req := &event_contracts.ListEventContractsRequest{BillingId: common.BillingId}
	sinksList, err := client.ListEventContracts(apiContext, req)
	common.CliExit(err)
	util.Print(sinksList)
}

func get(name *string) {
	eventContract := GetEventContract(name)
	util.Print(eventContract)
}

func GetEventContract(name *string) *entities.EventContract {
	req := &event_contracts.GetEventContractRequest{
		BillingId: common.BillingId,
		Ref:       ref(name)}
	eventContract, err := client.GetEventContract(apiContext, req)
	common.CliExit(err)
	return eventContract.EventContract
}
func refsCompletion(cmd *cobra.Command, args []string, complete string) ([]string, cobra.ShellCompDirective) {
	if len(args) > 0 || common.BillingIdIsMissing() {
		return common.MissingBillingIdCompletionError(cmd.CommandPath())
	}

	req := &event_contracts.ListEventContractsRequest{BillingId: common.BillingId}
	response, err := client.ListEventContracts(apiContext, req)

	if err != nil {
		return common.GrpcRequestCompletionError(err)
	}

	names := make([]string, 0, len(response.EventContracts))
	for _, s := range response.EventContracts {
		names = append(names, refToString(s.Ref))
	}

	return names, cobra.ShellCompDirectiveNoFileComp
}
