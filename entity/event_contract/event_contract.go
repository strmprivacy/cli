package event_contract

import (
	"context"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/streammachineio/api-definitions-go/api/entities/v1"
	"github.com/streammachineio/api-definitions-go/api/event_contracts/v1"
	"google.golang.org/grpc"
	"streammachine.io/strm/utils"
	"strings"
)

// strings used in the cli
const ()

var BillingId string
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

func ExistingNamesCompletion(cmd *cobra.Command, args []string, complete string) ([]string, cobra.ShellCompDirective) {
	if len(args) > 0 {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}
	return ExistingNames(), cobra.ShellCompDirectiveNoFileComp
}

func list() {
	req := &event_contracts.ListEventContractsRequest{BillingId: BillingId}
	sinksList, err := client.ListEventContracts(apiContext, req)
	cobra.CheckErr(err)
	utils.Print(sinksList)
}

func get(name *string) {
	event_contract := GetEventContract(name)
	utils.Print(event_contract)
}
func GetEventContract(name *string) *entities.EventContract {
	req := &event_contracts.GetEventContractRequest{
		BillingId: BillingId,
		Ref:       ref(name)}
	eventContract, err := client.GetEventContract(apiContext, req)
	cobra.CheckErr(err)
	return eventContract.EventContract
}

func ExistingNames() []string {
	req := &event_contracts.ListEventContractsRequest{BillingId: BillingId}
	response, err := client.ListEventContracts(apiContext, req)
	cobra.CheckErr(err)
	names := make([]string, 0)
	for _, s := range response.EventContracts {
		names = append(names, refToString(s.Ref))
	}
	return names
}
