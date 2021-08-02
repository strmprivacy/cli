package schema_code

import (
	"context"
	"github.com/streammachineio/api-definitions-go/api/schemas/v1"
	"google.golang.org/grpc"
	"streammachine.io/strm/pkg/common"
	"streammachine.io/strm/pkg/entity/schema"
	"streammachine.io/strm/pkg/util"
)

// strings used in the cli
const ()

var client schemas.SchemasServiceClient
var apiContext context.Context


func SetupClient(clientConnection *grpc.ClientConn, ctx context.Context) {
	apiContext = ctx
	client = schemas.NewSchemasServiceClient(clientConnection)
}

func get(name *string) {
	schemaCode := GetSchemaCode(name)
	util.Print(schemaCode)
}

func GetSchemaCode(name *string) *schemas.GetSchemaCodeResponse {
	req := &schemas.GetSchemaCodeRequest{
		BillingId: common.BillingId,
		Language: "typescript",
		Ref: schema.Ref(name)}
	schemaCode, err := client.GetSchemaCode(apiContext, req)
	common.CliExit(err)
	return schemaCode
}