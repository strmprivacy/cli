package schema_code

import (
	"context"
	"errors"
	"github.com/spf13/cobra"
	"github.com/strmprivacy/api-definitions-go/v2/api/schemas/v1"
	"google.golang.org/grpc"
	"os"
	"strmprivacy/strm/pkg/common"
	"strmprivacy/strm/pkg/entity/schema"
	"strmprivacy/strm/pkg/util"
)

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
			common.CliExit(errors.New("Not overwriting " + outputFile))
		}
	}
	saveFile(schemaCode, outputFile)
	return outputFile
}

func saveFile(code *schemas.GetSchemaCodeResponse, file string) {
	err := os.WriteFile(file, code.Data, 0666)
	common.CliExit(err)
}
