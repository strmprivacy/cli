package policy

import (
	"context"
	"errors"
	"fmt"
	"github.com/hashicorp/go-uuid"
	"github.com/samber/lo"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/strmprivacy/api-definitions-go/v3/api/entities/v1"
	policiesApi "github.com/strmprivacy/api-definitions-go/v3/api/policies/v1"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/fieldmaskpb"
	"strmprivacy/strm/pkg/common"
	"strmprivacy/strm/pkg/util"
)

const (
	idFlag            = "id"
	nameFlag          = "name"
	descriptionFlag   = "description"
	legalGroundsFlag  = "legal-grounds"
	retentionFlag     = "retention"
	policyNameFlag    = "policy-name" // used by stream and batch-job
	policyIdFlag      = "policy-id"   // used by stream and batch-job
	defaultPolicyFlag = "get-default-policy"
)

var client policiesApi.PoliciesServiceClient
var apiContext context.Context

func SetupClient(clientConnection *grpc.ClientConn, ctx context.Context) {
	apiContext = ctx
	client = policiesApi.NewPoliciesServiceClient(clientConnection)
}

// list all the policies owned by a certain organization
func list() {
	var err error
	response, err := client.ListPolicies(apiContext, &policiesApi.ListPoliciesRequest{})
	common.CliExit(err)
	printer.Print(response)
}

// get a policy by name (args[0]) or id option
func get(cmd *cobra.Command, args []string) {
	id := getPolicyIdFromArgumentsOrIdFlag(cmd.Flags(), args)
	response, err := client.GetPolicy(apiContext, &policiesApi.GetPolicyRequest{PolicyId: id})
	common.CliExit(err)
	printer.Print(response.Policy)
}

func activate(cmd *cobra.Command, args []string) {
	response, err := changeState(cmd, args, entities.Policy_STATE_ACTIVE)
	common.CliExit(err)
	printer.Print(response.Policy)
}

func archive(cmd *cobra.Command, args []string) {
	response, err := changeState(cmd, args, entities.Policy_STATE_ARCHIVED)
	common.CliExit(err)
	printer.Print(response.Policy)
}

func changeState(cmd *cobra.Command, args []string, state entities.Policy_State) (*policiesApi.UpdatePolicyResponse, error) {
	id := getPolicyIdFromArgumentsOrIdFlag(cmd.Flags(), args)
	if len(id) == 0 {
		// we don't allow updating the default policy state!
		if len(args) == 0 {
			// we tried to identify the policy by id
			common.Abort("No policy id %s found", util.GetStringAndErr(cmd.Flags(), idFlag))
		} else {
			common.Abort("No policy named %s found", args[0])
		}
	}
	request := &policiesApi.UpdatePolicyRequest{
		Policy: &entities.Policy{
			Id:    id,
			State: state,
		},
		UpdateMask: &fieldmaskpb.FieldMask{Paths: []string{"state"}},
	}
	return client.UpdatePolicy(apiContext, request)
}

// del delete a policy by name (args[0]) or id option
func del(cmd *cobra.Command, args []string) {
	id := getPolicyIdFromArgumentsOrIdFlag(cmd.Flags(), args)
	response, err := client.DeletePolicy(apiContext, &policiesApi.DeletePolicyRequest{PolicyId: id})
	common.CliExit(err)
	printer.Print(response)
}

// create a policy from
func create(cmd *cobra.Command) {
	flags := cmd.Flags()
	name := util.GetStringAndErr(flags, nameFlag)
	description := util.GetStringAndErr(flags, descriptionFlag)
	legalGrounds := util.GetStringAndErr(flags, legalGroundsFlag)
	retention := util.GetInt32AndErr(flags, retentionFlag)
	randomUuid, _ := uuid.GenerateUUID()
	policy := &entities.Policy{
		Id:            randomUuid,
		State:         entities.Policy_STATE_DRAFT,
		Name:          name,
		RetentionDays: retention,
		LegalGrounds:  legalGrounds,
		Description:   description,
	}
	request := &policiesApi.CreatePolicyRequest{Policy: policy}
	response, err := client.CreatePolicy(apiContext, request)
	common.CliExit(err)
	printer.Print(response.Policy)

}

func update(cmd *cobra.Command, id string) {
	flags := cmd.Flags()
	updateMask := make([]string, 0)
	// utility to check is a command flag was set and if so add it to the updateMask
	getOption := func(flagName string, updateName string) string {
		lookup := flags.Lookup(flagName)
		if lookup.Changed {
			updateMask = append(updateMask, updateName)
		}
		return lookup.Value.String()
	}
	name := getOption(nameFlag, "name")
	description := getOption(descriptionFlag, "description")
	legalGrounds := getOption(legalGroundsFlag, "legal_grounds")
	retention := util.GetInt32AndErr(flags, retentionFlag)
	if retention != 0 {
		updateMask = append(updateMask, "retention_days")
	}
	request := &policiesApi.UpdatePolicyRequest{
		Policy: &entities.Policy{
			Id:            id,
			Name:          name,
			RetentionDays: retention,
			LegalGrounds:  legalGrounds,
			Description:   description,
		},
		UpdateMask: &fieldmaskpb.FieldMask{Paths: updateMask},
	}
	response, err := client.UpdatePolicy(apiContext, request)
	common.CliExit(err)
	printer.Print(response.Policy)
}

var cachedPolicies *[]*entities.Policy

func getPolicies() (error, []*entities.Policy) {
	if cachedPolicies == nil {
		req := &policiesApi.ListPoliciesRequest{}
		response, err := client.ListPolicies(apiContext, req)
		cachedPolicies = &response.Policies
		return err, *cachedPolicies
	} else {
		return nil, *cachedPolicies
	}
}

// PoliciesNameIdMap link policy names to ids and vice versa
func PoliciesNameIdMap() (map[string]string, map[string]string) {
	var err, p = getPolicies()
	common.CliExit(err)
	return lo.Associate(p, func(f *entities.Policy) (string, string) {
			return f.Name, f.Id
		}), lo.Associate(p, func(f *entities.Policy) (string, string) {
			return f.Id, f.Name
		})
}

func namesCompletion(_ *cobra.Command, args []string, _ string) ([]string, cobra.ShellCompDirective) {
	fmt.Println(args)
	if len(args) != 0 {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}
	err, policies := getPolicies()
	if err != nil {
		return common.GrpcRequestCompletionError(err)
	}
	names := lo.Map(policies, func(policy *entities.Policy, _ int) string {
		return policy.Name
	})

	// TODO the completion actually doesn't work for names with spaces in them
	// I've tried adding quotes around the items in the names array, but those don't
	// show up in the output.
	// this will probably have to be solved in Cobra
	return names, cobra.ShellCompDirectiveNoFileComp
}

func idsCompletion(_ *cobra.Command, args []string, _ string) ([]string, cobra.ShellCompDirective) {
	if len(args) != 0 {
		// this one means you don't get multiple completion suggestions for one policy
		return nil, cobra.ShellCompDirectiveNoFileComp
	}
	err, policies := getPolicies()
	if err != nil {
		return common.GrpcRequestCompletionError(err)
	}
	ids := lo.Map(policies, func(policy *entities.Policy, _ int) string {
		return policy.Id
	})
	return ids, cobra.ShellCompDirectiveNoFileComp
}

// SetupFlags adds policy-name and if flags to a Cobra command
// used during stream and batch-job creation
func SetupFlags(stream *cobra.Command, flags *pflag.FlagSet) {
	flags.String(policyNameFlag, "", "the name of the policy to attach")
	flags.String(policyIdFlag, "", "the uuid of the policy to attach")
	err := stream.RegisterFlagCompletionFunc(policyNameFlag, namesCompletion)
	err = stream.RegisterFlagCompletionFunc(policyIdFlag, idsCompletion)
	common.CliExit(err)
}

// GetPolicyFromFlags retrieves a policy from either the policy-name or id flag
// used by stream and batch-job creation
func GetPolicyFromFlags(flags *pflag.FlagSet) string {
	policyId := util.GetStringAndErr(flags, policyIdFlag)
	if len(policyId) == 0 {
		policyName := util.GetStringAndErr(flags, policyNameFlag)
		if len(policyName) != 0 {
			m, _ := PoliciesNameIdMap()
			policyId = m[policyName]
		}
	}
	if len(policyId) > 0 {
		_, err := uuid.ParseUUID(policyId)
		if err != nil {
			common.CliExit(errors.New(fmt.Sprintf("%s is not a uuid", policyId)))
		}
	}
	return policyId

}

// getPolicyIdFromArgumentsOrIdFlag get a policy id from either an
// option or as the first string in args
func getPolicyIdFromArgumentsOrIdFlag(flags *pflag.FlagSet, args []string) string {
	defaultPolicy := flags.Lookup(defaultPolicyFlag)
	if defaultPolicy != nil && defaultPolicy.Changed {
		return ""
	}
	id := util.GetStringAndErr(flags, idFlag)
	if len(id) == 0 {
		if len(args) == 0 {
			common.Abort("No name argument and no --id option")
		}
		name := args[0]
		if len(name) == 0 {
			common.Abort("name argument is an empty string")
		}
		nameIdMap, _ := PoliciesNameIdMap()
		id = nameIdMap[name]
		if len(id) == 0 {
			common.Abort("There's no policy named %s for your organization", name)
		}
	}
	return id
}
