package event_contract

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/streammachineio/api-definitions-go/api/entities/v1"
	"github.com/streammachineio/api-definitions-go/api/event_contracts/v1"
	"google.golang.org/grpc"
	"io/ioutil"
	"streammachine.io/strm/pkg/common"
	"streammachine.io/strm/pkg/entity/schema"
	"streammachine.io/strm/pkg/util"
	"strings"
)

// strings used in the cli
const ()

var client event_contracts.EventContractsServiceClient
var apiContext context.Context

type EventContractDefinition struct {
	KeyField    string                 `json:"keyField"`
	PiiFields   map[string]int32       `json:"piiFields"`
	Validations []*entities.Validation `json:"validations"`
}

func ref(refString *string) *entities.EventContractRef {
	parts := strings.Split(*refString, "/")

	if len(parts) != 3 {
		common.CliExit("Event Contract reference should consist of three parts: <handle>/<name>/<version>")
	}

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

func create(cmd *cobra.Command, contractReference *string) {
	flags := cmd.Flags()

	isPublic := util.GetBoolAndErr(flags, isPublicFlag)
	schemaRef := util.GetStringAndErr(flags, schemaRefFlag)
	definitionFilename := util.GetStringAndErr(flags, definitionFile)
	definition := readContractDefinition(&definitionFilename)

	req := &event_contracts.CreateEventContractRequest{
		BillingId: common.BillingId,
		EventContract: &entities.EventContract{
			Ref:         ref(contractReference),
			SchemaRef:   schema.Ref(&schemaRef),
			IsPublic:    isPublic,
			KeyField:    definition.KeyField,
			PiiFields:   definition.PiiFields,
			Validations: definition.Validations,
		},
	}

	response, err := client.CreateEventContract(apiContext, req)
	common.CliExit(err)
	util.Print(response)
}

func readContractDefinition(filename *string) EventContractDefinition {
	file, _ := ioutil.ReadFile(*filename)

	contractDefinition := EventContractDefinition{}
	err := json.Unmarshal(file, &contractDefinition)

	common.CliExit(err)

	return contractDefinition
}

func refsCompletion(cmd *cobra.Command, args []string, complete string) ([]string, cobra.ShellCompDirective) {
	if common.BillingIdIsMissing() {
		return common.MissingBillingIdCompletionError(cmd.CommandPath())
	}
	if len(args) != 0 {
		// this one means you don't get two completion suggestions for one stream
		return nil, cobra.ShellCompDirectiveNoFileComp
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
