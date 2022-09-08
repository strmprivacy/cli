package data_contract

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/strmprivacy/api-definitions-go/v2/api/data_contracts/v1"
	"github.com/strmprivacy/api-definitions-go/v2/api/entities/v1"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/encoding/protojson"
	"io/ioutil"
	"sigs.k8s.io/yaml"
	"strings"
	"strmprivacy/strm/pkg/common"
	"strmprivacy/strm/pkg/util"
)

const (
	schemaDefinitionFlag   = "schema-definition"
	publicFlag             = "public"
	schemaTypeFlag         = "type"
	contractDefinitionFlag = "contract-definition"
)

var client data_contracts.DataContractsServiceClient
var apiContext context.Context

func SetupClient(clientConnection *grpc.ClientConn, ctx context.Context) {
	apiContext = ctx
	client = data_contracts.NewDataContractsServiceClient(clientConnection)
}

type DataContractDefinition struct {
	KeyField         string                 `json:"keyField"`
	PiiFields        map[string]int32       `json:"piiFields"`
	Validations      []*entities.Validation `json:"validations"`
	DataSubjectField string                 `json:"dataSubjectField"`
}

func readContractDefinition(filename *string) DataContractDefinition {
	file, _ := ioutil.ReadFile(*filename)

	contractDefinition := DataContractDefinition{}
	err := json.Unmarshal(file, &contractDefinition)

	common.CliExit(err)

	return contractDefinition
}

func getSchemaDefinition(filename string, ref *entities.DataContractRef, schemaType int32, isPublic bool) entities.Schema {
	// try yaml
	definition, err := ioutil.ReadFile(filename)

	convertedToJson, err := yaml.YAMLToJSON(definition)
	if err == nil {
		definition = convertedToJson
	}
	simple := &entities.Schema_SimpleSchemaDefinition{}

	// try json
	err = protojson.Unmarshal(definition, simple)
	schema := entities.Schema{
		Ref: &entities.SchemaRef{
			Name:       ref.Name,
			Handle:     ref.Handle,
			Version:    ref.Version,
			SchemaType: entities.SchemaType(schemaType),
		},
		IsPublic: isPublic,
		Metadata: &entities.SchemaMetadata{},
	}

	if err == nil {
		// it's a simple schema
		schema.SimpleSchema = simple
	} else {
		schema.Definition = string(definition)
	}
	return schema
}

func create(cmd *cobra.Command, args *string) {
	flags := cmd.Flags()
	typeString := util.GetStringAndErr(flags, schemaTypeFlag)
	schemaType, ok := entities.SchemaType_value[typeString]
	if !ok {
		common.CliExit(errors.New(fmt.Sprintf("Can't convert %s to a known consent schema type, types are %v",
			typeString, entities.SchemaType_value)))
	}
	schemaDefinitionFilename := util.GetStringAndErr(flags, schemaDefinitionFlag)
	isPublic := util.GetBoolAndErr(flags, publicFlag)
	contractDefinitionFilename := util.GetStringAndErr(flags, contractDefinitionFlag)
	contractDefinition := readContractDefinition(&contractDefinitionFilename)

	ref := ref(args)
	schema := getSchemaDefinition(schemaDefinitionFilename, ref, schemaType, isPublic)

	req := &data_contracts.CreateDataContractRequest{
		ProjectId: common.ProjectId,
		DataContract: &entities.DataContract{
			KeyField:         contractDefinition.KeyField,
			IsPublic:         isPublic,
			ProjectId:        common.ProjectId,
			DataSubjectField: contractDefinition.DataSubjectField,
			Schema:           &schema,
			Ref:              ref,
			PiiFields:        contractDefinition.PiiFields,
			Validations:      contractDefinition.Validations,
			Metadata:         &entities.DataContractMetadata{},
		},
	}
	fmt.Println(req)
	response, err := client.CreateDataContract(apiContext, req)
	common.CliExit(err)
	printer.Print(response)
}

func list() {
	req := &data_contracts.ListDataContractsRequest{
		ProjectId: common.ProjectId,
	}

	response, err := client.ListDataContracts(apiContext, req)
	common.CliExit(err)

	printer.Print(response)
}

func del(refString *string) {
	req := &data_contracts.DeleteDataContractRequest{
		DataContractRef: ref(refString),
	}
	_, err := client.DeleteDataContract(apiContext, req)
	common.CliExit(err)
}

func get(name *string) {
	req := &data_contracts.GetDataContractRequest{
		Ref: ref(name),
	}
	response, err := client.GetDataContract(apiContext, req)
	common.CliExit(err)

	printer.Print(response)
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
