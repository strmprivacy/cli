package policy

import (
	"context"
	"errors"
	"fmt"
	"github.com/hashicorp/go-uuid"
	"github.com/samber/lo"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/strmprivacy/api-definitions-go/v2/api/entities/v1"
	policiesApi "github.com/strmprivacy/api-definitions-go/v2/api/policies/v1"
	"golang.org/x/exp/maps"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/fieldmaskpb"
	"strings"
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
	policyNameFlag   = "policy-name" // used by stream and batch-job
	policyIdFlag     = "policy-id"   // used by stream and batch-job
	updateMaskFlag   = "update-mask"
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

// del a policy by name (args[0]) or id option
func del(cmd *cobra.Command, args []string) {
	id := getPolicyIdFromArgumentsOrIdFlag(cmd.Flags(), args)
	response, err := client.DeletePolicy(apiContext, &policiesApi.DeletePolicyRequest{PolicyId: id})
	common.CliExit(err)
	printer.Print(response)
}

func create(cmd *cobra.Command) {
	flags := cmd.Flags()
	name := util.GetStringAndErr(flags, nameFlag)
	description := util.GetStringAndErr(flags, descriptionFlag)
	legalGrounds := util.GetStringAndErr(flags, legalGroundsFlag)
	retention := util.GetInt32AndErr(flags, retentionFlag)
	randomUuid, _ := uuid.GenerateUUID()
	state := stateEnum(flags)
	if state == entities.Policy_STATE_UNSPECIFIED {
		state = entities.Policy_STATE_DRAFT
	}
	policy := &entities.Policy{
		Id:            randomUuid,
		State:         state,
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
	name := util.GetStringAndErr(flags, nameFlag)
	if len(name) != 0 {
		updateMask = append(updateMask, "name")
	}
	description := util.GetStringAndErr(flags, descriptionFlag)
	if len(description) != 0 {
		updateMask = append(updateMask, "description")
	}
	legalGrounds := util.GetStringAndErr(flags, legalGroundsFlag)
	if len(legalGrounds) != 0 {
		updateMask = append(updateMask, "legal_grounds")
	}
	retention := util.GetInt32AndErr(flags, retentionFlag)
	if retention != 0 {
		updateMask = append(updateMask, "retention_days")
	}
	state := stateEnum(flags)
	if state != entities.Policy_STATE_UNSPECIFIED {
		updateMask = append(updateMask, "state")
	}
	// this is for the rare case where you want to clear a value (say description)
	// `--description "" --update-mask description
	m, err := flags.GetStringSlice(updateMaskFlag)
	updateMask = append(updateMask, m...)
	common.CliExit(err)

	request := &policiesApi.UpdatePolicyRequest{
		Policy: &entities.Policy{
			Id:            id,
			State:         state,
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

// stateEnum converts a cli option flag value to a Policy_State
// it will automatically uppercase and add a 'STATE_' prefix
func stateEnum(flags *pflag.FlagSet) entities.Policy_State {
	flagValue := strings.ToUpper(util.GetStringAndErr(flags, stateFlag))
	if len(flagValue) == 0 {
		return entities.Policy_STATE_UNSPECIFIED
	}
	if !strings.HasPrefix(flagValue, "STATE_") {
		flagValue = "STATE_" + flagValue
	}
	if val, ok := entities.Policy_State_value[flagValue]; ok {
		return entities.Policy_State(val)
	} else {
		common.Abort("State %s is not one of %s", flagValue, maps.Keys(entities.Policy_State_value))
		return 0
	}
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

func idsCompletion(_ *cobra.Command, args []string, complete string) ([]string, cobra.ShellCompDirective) {
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
