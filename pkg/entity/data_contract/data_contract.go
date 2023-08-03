package data_contract

import (
	"context"
	"errors"
	"github.com/spf13/cobra"
	"github.com/strmprivacy/api-definitions-go/v3/api/data_contracts/v1"
	"github.com/strmprivacy/api-definitions-go/v3/api/entities/v1"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/encoding/protojson"
	"os"
	"sigs.k8s.io/yaml"
	"strings"
	"strmprivacy/strm/pkg/common"
	"strmprivacy/strm/pkg/entity/project"
	"strmprivacy/strm/pkg/util"
)

const (
	schemaDefinitionFlag   = "schema-definition"
	publicFlag             = "public"
	contractDefinitionFlag = "contract-definition"
)

var client data_contracts.DataContractsServiceClient
var apiContext context.Context

func SetupClient(clientConnection *grpc.ClientConn, ctx context.Context) {
	apiContext = ctx
	client = data_contracts.NewDataContractsServiceClient(clientConnection)
}

func readContractDefinition(filename *string) *entities.DataContract {
	file, _ := os.ReadFile(*filename)

	contractDefinition := &entities.DataContract{}
	err := protojson.Unmarshal(file, contractDefinition)

	common.CliExit(err)

	return contractDefinition
}

func readSchemaDefinition(filename string, ref *entities.DataContractRef, isPublic bool) *entities.Schema {
	schema := entities.Schema{
		Ref: &entities.SchemaRef{
			Name:    ref.Name,
			Handle:  ref.Handle,
			Version: ref.Version,
		},
		IsPublic: isPublic,
		Metadata: &entities.SchemaMetadata{},
	}

	definition, err := os.ReadFile(filename)

	// Try to convert YAML to JSON
	convertedToJson, err := yaml.YAMLToJSON(definition)
	if err == nil {
		definition = convertedToJson
	}

	// Try to unmarshal the JSON as Simple Schema
	simple := &entities.Schema_SimpleSchemaDefinition{}
	err = protojson.Unmarshal(definition, simple)
	if err == nil {
		// It's a Simple Schema
		schema.SimpleSchema = simple
	} else {
		// It's an Avro or JsonSchema definition
		schema.Definition = string(definition)
	}
	return &schema
}

func create(cmd *cobra.Command, args *string) {
	flags := cmd.Flags()

	schemaDefinitionFilename := util.GetStringAndErr(flags, schemaDefinitionFlag)
	isPublic := util.GetBoolAndErr(flags, publicFlag)
	contractDefinitionFilename := util.GetStringAndErr(flags, contractDefinitionFlag)
	dataContract := readContractDefinition(&contractDefinitionFilename)

	projectId := project.GetProjectId(cmd)
	ref := ref(args)
	schema := readSchemaDefinition(schemaDefinitionFilename, ref, isPublic)

	// Set the fields which are not part of the contract definition file
	dataContract.Ref = ref
	dataContract.IsPublic = isPublic
	dataContract.ProjectId = projectId
	dataContract.Schema = schema
	dataContract.Metadata = &entities.DataContractMetadata{}

	req := &data_contracts.CreateDataContractRequest{
		ProjectId:    projectId,
		DataContract: dataContract,
	}
	response, err := client.CreateDataContract(apiContext, req)
	common.CliExit(err)
	printer.Print(response)
}

func list(cmd *cobra.Command) {
	req := &data_contracts.ListDataContractsRequest{
		ProjectId: project.GetProjectId(cmd),
	}

	response, err := client.ListDataContracts(apiContext, req)
	common.CliExit(err)

	printer.Print(response)
}

func del(refString *string, cmd *cobra.Command) {
	req := &data_contracts.DeleteDataContractRequest{
		ProjectId:       project.GetProjectId(cmd),
		DataContractRef: ref(refString),
	}
	_, err := client.DeleteDataContract(apiContext, req)
	common.CliExit(err)
}

func get(refString *string) {
	req := &data_contracts.GetDataContractRequest{
		Ref: ref(refString),
	}
	response, err := client.GetDataContract(apiContext, req)
	common.CliExit(err)

	printer.Print(response)
}

func review(refString *string) {
	req := &data_contracts.ReviewDataContractRequest{
		DataContractRef: ref(refString),
	}

	_, err := client.ReviewDataContract(apiContext, req)
	common.CliExit(err)
}

func approve(refString *string) {
	req := &data_contracts.ApproveDataContractRequest{
		DataContractRef: ref(refString),
	}

	_, err := client.ApproveDataContract(apiContext, req)
	common.CliExit(err)
}

func activate(refString *string) {
	req := &data_contracts.ActivateDataContractRequest{
		DataContractRef: ref(refString),
	}
	_, err := client.ActivateDataContract(apiContext, req)
	common.CliExit(err)
}

func archive(refString *string) {
	req := &data_contracts.ArchiveDataContractRequest{
		DataContractRef: ref(refString),
	}
	_, err := client.ArchiveDataContract(apiContext, req)
	common.CliExit(err)
}

func RefsCompletion(cmd *cobra.Command, args []string, complete string) ([]string, cobra.ShellCompDirective) {
	if len(args) != 0 {
		// this one means you don't get multiple completion suggestions for one stream
		return nil, cobra.ShellCompDirectiveNoFileComp
	}

	req := &data_contracts.ListDataContractsRequest{}
	response, err := client.ListDataContracts(apiContext, req)

	if err != nil {
		return common.GrpcRequestCompletionError(err)
	}

	names := make([]string, 0, len(response.DataContracts))
	for _, s := range response.DataContracts {
		names = append(names, refToString(s.Ref))
	}

	return names, cobra.ShellCompDirectiveNoFileComp
}

func ref(refString *string) *entities.DataContractRef {
	parts := strings.Split(*refString, "/")

	if len(parts) != 3 {
		common.CliExit(errors.New("Event Contract reference should consist of three parts: <handle>/<name>/<version>"))
	}

	return &entities.DataContractRef{
		Handle:  parts[0],
		Name:    parts[1],
		Version: parts[2],
	}
}
