package policy

import (
	"context"
	"errors"
	"fmt"
	"github.com/hashicorp/go-uuid"
	"github.com/spf13/cobra"
	"github.com/strmprivacy/api-definitions-go/v2/api/entities/v1"
	"github.com/strmprivacy/api-definitions-go/v2/api/policies/v1"
	"google.golang.org/grpc"
	"strmprivacy/strm/pkg/common"
	"strmprivacy/strm/pkg/util"
)

const (
	idFlag           = "id"
	nameFlag         = "name"
	descriptionFlag  = "description"
	legalGroundsFlag = "legal-grounds"
	stateFlag        = "state"
	retentionFlag    = "retention"
)

var client policies.PoliciesServiceClient
var apiContext context.Context

func SetupClient(clientConnection *grpc.ClientConn, ctx context.Context) {
	apiContext = ctx
	client = policies.NewPoliciesServiceClient(clientConnection)
}

func list(cmd *cobra.Command) {
	var err error
	//flags := cmd.Flags()
	req := &policies.ListPoliciesRequest{}
	response, err := client.ListPolicies(apiContext, req)
	common.CliExit(err)
	printer.Print(response)
}

func get(cmd *cobra.Command, args []string) {
	var err error
	flags := cmd.Flags()
	id := util.GetStringAndErr(flags, idFlag)
	if id != "" {
		req := &policies.GetPolicyRequest{PolicyId: id}
		response, _ := client.GetPolicy(apiContext, req)
		common.CliExit(err)
		printer.Print(response.Policy)
	} else {
		req := &policies.ListPoliciesRequest{}
		response, err := client.ListPolicies(apiContext, req)
		common.CliExit(err)
		for _, p := range response.Policies {
			if p.Name == args[0] {
				printer.Print(p)
				return
			}
		}
		common.CliExit(errors.New(fmt.Sprintf("no policy with name %s exists in your organization", args[0])))
	}
}

func create(cmd *cobra.Command) {
	flags := cmd.Flags()
	name := util.GetStringAndErr(flags, nameFlag)
	description := util.GetStringAndErr(flags, descriptionFlag)
	legalGrounds := util.GetStringAndErr(flags, legalGroundsFlag)
	retention := util.GetInt32AndErr(flags, retentionFlag)

	randomUuid, _ := uuid.GenerateUUID()
	policy := &entities.Policy{
		Id:            randomUuid,
		State:         entities.Policy_State(entities.Policy_State_value[util.GetStringAndErr(flags, stateFlag)]),
		Name:          name,
		RetentionDays: retention,
		LegalGrounds:  legalGrounds,
		Description:   description,
	}
	request := &policies.CreatePolicyRequest{Policy: policy}
	response, err := client.CreatePolicy(apiContext, request)
	common.CliExit(err)
	printer.Print(response.Policy)

}

func NamesCompletion(cmd *cobra.Command, args []string, complete string) ([]string, cobra.ShellCompDirective) {
	if len(args) != 0 {
		// this one means you don't get multiple completion suggestions for one policy
		return nil, cobra.ShellCompDirectiveNoFileComp
	}

	req := &policies.ListPoliciesRequest{}
	response, err := client.ListPolicies(apiContext, req)

	if err != nil {
		return common.GrpcRequestCompletionError(err)
	}

	// TODO the completion actually doesn't work for names with spaces in them
	// I've tried adding quotes around the items in the names array, but those don't
	// show up in the output.
	names := make([]string, 0, len(response.Policies))
	for _, s := range response.Policies {
		names = append(names, s.Name)
	}
	return names, cobra.ShellCompDirectiveNoFileComp
}

func IdsCompletion(cmd *cobra.Command, args []string, complete string) ([]string, cobra.ShellCompDirective) {
	if len(args) != 0 {
		// this one means you don't get multiple completion suggestions for one policy
		return nil, cobra.ShellCompDirectiveNoFileComp
	}

	req := &policies.ListPoliciesRequest{}
	response, err := client.ListPolicies(apiContext, req)

	if err != nil {
		return common.GrpcRequestCompletionError(err)
	}

	ids := make([]string, 0, len(response.Policies))
	for _, s := range response.Policies {
		ids = append(ids, s.Id)
	}
	return ids, cobra.ShellCompDirectiveNoFileComp
}
