package schema_code

import (
	"context"
	"github.com/spf13/cobra"
	"github.com/streammachineio/api-definitions-go/api/schemas/v1"
	"google.golang.org/grpc"
	"os"
	"streammachine.io/strm/pkg/common"
	"streammachine.io/strm/pkg/entity/schema"
	"streammachine.io/strm/pkg/util"
)

// strings used in the cli
const (
	languageFlag  = "language"
	filenameFlag  = "output-file"
	overwriteFlag = "overwrite"
)

var client schemas.SchemasServiceClient
var apiContext context.Context

func SetupClient(clientConnection *grpc.ClientConn, ctx context.Context) {
	apiContext = ctx
	client = schemas.NewSchemasServiceClient(clientConnection)
}

func get(cmd *cobra.Command, schemaRef *string) {
	outputFile := GetSchemaCode(cmd, schemaRef)
	println("Saved to", outputFile)
}

func GetSchemaCode(cmd *cobra.Command, name *string) string {
	flags := cmd.Flags()
	language := util.GetStringAndErr(flags, languageFlag)
	outputFile := util.GetStringAndErr(flags, filenameFlag)
	overwrite := util.GetBoolAndErr(flags, overwriteFlag)
	req := &schemas.GetSchemaCodeRequest{
		BillingId: common.BillingId,
		Language:  language,
		Ref:       schema.Ref(name),
	}
	schemaCode, err := client.GetSchemaCode(apiContext, req)
	common.CliExit(err)
	if len(outputFile) == 0 {
		outputFile = schemaCode.Filename
	}
	if !overwrite {
		_, err = os.Stat(outputFile)
		if !os.IsNotExist(err) {
			common.CliExit("Not overwriting " + outputFile)
		}
	}
	saveFile(schemaCode, outputFile)
	return outputFile
}

func saveFile(code *schemas.GetSchemaCodeResponse, file string) {
	err := os.WriteFile(file, code.Data, 0666)
	common.CliExit(err)
}
