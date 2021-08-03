package schema

import (
	"context"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/streammachineio/api-definitions-go/api/entities/v1"
	"github.com/streammachineio/api-definitions-go/api/schemas/v1"
	"google.golang.org/grpc"
	"io/ioutil"
	"streammachine.io/strm/pkg/common"
	"streammachine.io/strm/pkg/util"
	"strings"
)

// strings used in the cli
const (
	kafkaClusterFlag = "kafka-cluster"
	definitionFlag   = "definition"
	publicFlag       = "public"
)

var client schemas.SchemasServiceClient
var apiContext context.Context

func Ref(n *string) *entities.SchemaRef {
	parts := strings.Split(*n, "/")
	return &entities.SchemaRef{
		Handle:  parts[0],
		Name:    parts[1],
		Version: parts[2],
	}
}

func refToString(ref *entities.SchemaRef) string {
	return fmt.Sprintf("%v/%v/%v", ref.Handle, ref.Name, ref.Version)
}

func SetupClient(clientConnection *grpc.ClientConn, ctx context.Context) {
	apiContext = ctx
	client = schemas.NewSchemasServiceClient(clientConnection)
}

func list() {
	req := &schemas.ListSchemasRequest{BillingId: common.BillingId}
	sinksList, err := client.ListSchemas(apiContext, req)
	common.CliExit(err)
	util.Print(sinksList)
}

func get(name *string, cmd *cobra.Command) {
	flags := cmd.Flags()
	clusterRef, err := getClusterRef(flags)
	common.CliExit(err)

	schema := GetSchema(name, clusterRef)
	util.Print(schema)
}

func getClusterRef(flags *pflag.FlagSet) (*entities.KafkaClusterRef, error) {
	flag := util.GetStringAndErr(flags, kafkaClusterFlag)
	if len(flag) > 0 {
		parts := strings.Split(flag, "/")
		if len(parts) == 2 {
			return &entities.KafkaClusterRef{
				BillingId: parts[0],
				Name:      parts[1],
			}, nil
		} else {
			return nil, fmt.Errorf("invalid %v. Should be formatted as 'billing_id/cluster_name'", kafkaClusterFlag)
		}
	} else {
		return &entities.KafkaClusterRef{}, nil
	}
}

func GetSchema(name *string, clusterRef *entities.KafkaClusterRef) *schemas.GetSchemaResponse {
	req := &schemas.GetSchemaRequest{
		BillingId:  common.BillingId,
		Ref:        Ref(name),
		ClusterRef: clusterRef,
	}
	response, err := client.GetSchema(apiContext, req)
	common.CliExit(err)
	return response
}

func create(cmd *cobra.Command, args *string) {
	flags := cmd.Flags()

	definitionFilename := util.GetStringAndErr(flags, definitionFlag)
	definition, err := ioutil.ReadFile(definitionFilename)

	isPublic := util.GetBoolAndErr(flags, publicFlag)

	ref := Ref(args)
	req := &schemas.CreateSchemaRequest{
		BillingId: common.BillingId,
		Schema: &entities.Schema{
			Ref:        ref,
			Definition: string(definition),
			IsPublic:   isPublic,
		},
	}
	response, err := client.CreateSchema(apiContext, req)
	common.CliExit(err)
	util.Print(response)
}

func NamesCompletion(cmd *cobra.Command, args []string, complete string) ([]string, cobra.ShellCompDirective) {
	if common.BillingIdIsMissing() {
		return common.MissingBillingIdCompletionError(cmd.CommandPath())
	}
	if len(args) != 0 {
		// this one means you don't get two completion suggestions for one stream
		return nil, cobra.ShellCompDirectiveNoFileComp
	}

	req := &schemas.ListSchemasRequest{BillingId: common.BillingId}
	response, err := client.ListSchemas(apiContext, req)

	if err != nil {
		return common.GrpcRequestCompletionError(err)
	}

	names := make([]string, 0)
	for _, s := range response.Schemas {
		names = append(names, refToString(s.Ref))
	}

	return names, cobra.ShellCompDirectiveNoFileComp
}
