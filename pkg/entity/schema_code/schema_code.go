package schema_code

import (
	"context"
	"errors"
	"github.com/spf13/cobra"
	"github.com/strmprivacy/api-definitions-go/v3/api/data_contracts/v1"
	"github.com/strmprivacy/api-definitions-go/v3/api/entities/v1"
	"google.golang.org/grpc"
	"os"
	"strings"
	"strmprivacy/strm/pkg/common"
	"strmprivacy/strm/pkg/util"
)

const (
	languageFlag  = "language"
	filenameFlag  = "output-file"
	overwriteFlag = "overwrite"
)

var client data_contracts.DataContractsServiceClient
var apiContext context.Context

func SetupClient(clientConnection *grpc.ClientConn, ctx context.Context) {
	apiContext = ctx
	client = data_contracts.NewDataContractsServiceClient(clientConnection)
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
	req := &data_contracts.GetDataContractSchemaCodeRequest{
		Language:        language,
		DataContractRef: dataContractRef(name),
	}
	schemaCode, err := client.GetDataContractSchemaCode(apiContext, req)
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

func saveFile(code *data_contracts.GetDataContractSchemaCodeResponse, file string) {
	err := os.WriteFile(file, code.Data, 0666)
	common.CliExit(err)
}

func dataContractRef(refString *string) *entities.DataContractRef {
	parts := strings.Split(*refString, "/")

	if len(parts) != 3 {
		common.CliExit(errors.New("Data contract reference should consist of three parts: <handle>/<name>/<version>"))
	}

	return &entities.DataContractRef{
		Handle:  parts[0],
		Name:    parts[1],
		Version: parts[2],
	}
}
