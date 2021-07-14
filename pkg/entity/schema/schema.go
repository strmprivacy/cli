package schema

import (
	"context"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/streammachineio/api-definitions-go/api/entities/v1"
	"github.com/streammachineio/api-definitions-go/api/schemas/v1"
	"google.golang.org/grpc"
	"streammachine.io/strm/pkg/common"
	"streammachine.io/strm/pkg/util"
	"strings"
)

// strings used in the cli
const ()

var client schemas.SchemasServiceClient
var apiContext context.Context

func ref(n *string) *entities.SchemaRef {
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

func get(name *string) {
	schema := GetSchema(name)
	util.Print(schema)
}

func GetSchema(name *string) *entities.Schema {
	req := &schemas.GetSchemaRequest{
		BillingId: common.BillingId,
		Ref:       ref(name)}
	schema, err := client.GetSchema(apiContext, req)
	common.CliExit(err)
	return schema.Schema
}
func namesCompletion(cmd *cobra.Command, args []string, complete string) ([]string, cobra.ShellCompDirective) {
	if len(args) > 0 || common.BillingIdIsMissing() {
		return common.MissingBillingIdCompletionError(cmd.CommandPath())
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
