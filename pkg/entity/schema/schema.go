package schema

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"strings"
	"strmprivacy/strm/pkg/entity/project"

	"sigs.k8s.io/yaml"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	entities "github.com/strmprivacy/api-definitions-go/v2/api/entities/v1"
	schemas "github.com/strmprivacy/api-definitions-go/v2/api/schemas/v1"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/encoding/protojson"
	"strmprivacy/strm/pkg/common"
	"strmprivacy/strm/pkg/util"
)

const (
	kafkaClusterFlag = "kafka-cluster"
	definitionFlag   = "definition"
	publicFlag       = "public"
	schemaTypeFlag   = "type"
	projectName      = "project"
)

var client schemas.SchemasServiceClient
var apiContext context.Context

func Ref(refString *string) *entities.SchemaRef {
	parts := strings.Split(*refString, "/")

	if len(parts) != 3 {
		common.CliExit(errors.New("Schema reference should consist of three parts: <handle>/<name>/<version>"))
	}

	return &entities.SchemaRef{
		Handle:  parts[0],
		Name:    parts[1],
		Version: parts[2],
	}
}

func RefToString(ref *entities.SchemaRef) string {
	return fmt.Sprintf("%v/%v/%v", ref.Handle, ref.Name, ref.Version)
}

func SetupClient(clientConnection *grpc.ClientConn, ctx context.Context) {
	apiContext = ctx
	client = schemas.NewSchemasServiceClient(clientConnection)
}

func list() {
	req := &schemas.ListSchemasRequest{}
	response, err := client.ListSchemas(apiContext, req)
	common.CliExit(err)

	printer.Print(response)
}

func del(name *string) {
	req := &schemas.DeleteSchemaRequest{
		ProjectId: common.ProjectId,
		SchemaRef: Ref(name)}
	response, err := client.DeleteSchema(apiContext, req)
	common.CliExit(err)

	printer.Print(response)
}

func activate(name *string) {
	req := &schemas.ActivateSchemaRequest{
		ProjectId: common.ProjectId,
		SchemaRef: Ref(name)}
	response, err := client.ActivateSchema(apiContext, req)
	common.CliExit(err)

	printer.Print(response)
}

func archive(name *string) {
	req := &schemas.ArchiveSchemaRequest{
		ProjectId: common.ProjectId,
		SchemaRef: Ref(name)}
	response, err := client.ArchiveSchema(apiContext, req)
	common.CliExit(err)

	printer.Print(response)
}

func get(name *string, cmd *cobra.Command) {
	flags := cmd.Flags()
	clusterRef, err := getClusterRef(flags)
	common.CliExit(err)

	response := GetSchema(name, clusterRef)
	printer.Print(response)
}

func getClusterRef(flags *pflag.FlagSet) (*entities.KafkaClusterRef, error) {
	flag := util.GetStringAndErr(flags, kafkaClusterFlag)
	if len(flag) > 0 {
		return &entities.KafkaClusterRef{
			// Todo: actual transition to different ref
			ProjectId: common.ProjectId,
			Name:      flag,
		}, nil
	} else {
		return &entities.KafkaClusterRef{}, nil
	}
}

func GetSchema(name *string, clusterRef *entities.KafkaClusterRef) *schemas.GetSchemaResponse {
	req := &schemas.GetSchemaRequest{
		Ref:        Ref(name),
		ClusterRef: clusterRef,
	}
	response, err := client.GetSchema(apiContext, req)
	common.CliExit(err)
	return response
}

func create(cmd *cobra.Command, args *string) {
	flags := cmd.Flags()
	projectName := util.GetStringAndErr(flags, projectName)
	var projectId string
	if len(projectName) > 0 {
		projectId = project.GetProjectId(projectName)
	} else {
		projectId = common.ProjectId
	}
	typeString := util.GetStringAndErr(flags, schemaTypeFlag)
	schemaType, ok := entities.SchemaType_value[typeString]
	if !ok {
		common.CliExit(errors.New(fmt.Sprintf("Can't convert %s to a known consent schema type, types are %v",
			typeString, entities.SchemaType_value)))
	}
	definitionFilename := util.GetStringAndErr(flags, definitionFlag)
	definition, err := ioutil.ReadFile(definitionFilename)
	simple := &entities.Schema_SimpleSchemaDefinition{}
	isPublic := util.GetBoolAndErr(flags, publicFlag)
	ref := Ref(args)
	ref.SchemaType = entities.SchemaType(schemaType)
	req := &schemas.CreateSchemaRequest{
		ProjectId: projectId,
		Schema: &entities.Schema{
			Ref:      ref,
			IsPublic: isPublic,
		},
	}
	// try yaml
	convertedToJson, err := yaml.YAMLToJSON(definition)
	if err == nil {
		definition = convertedToJson
	}
	// try json
	err = protojson.Unmarshal(definition, simple)
	if err == nil {
		// it's a simple schema
		req.Schema.SimpleSchema = simple
	} else {
		req.Schema.Definition = string(definition)
	}

	response, err := client.CreateSchema(apiContext, req)
	common.CliExit(err)
	printer.Print(response)
}

func RefsCompletion(cmd *cobra.Command, args []string, complete string) ([]string, cobra.ShellCompDirective) {

	req := &schemas.ListSchemasRequest{}
	response, err := client.ListSchemas(apiContext, req)

	if err != nil {
		return common.GrpcRequestCompletionError(err)
	}

	names := make([]string, 0)
	for _, s := range response.Schemas {
		names = append(names, RefToString(s.Ref))
	}

	return names, cobra.ShellCompDirectiveNoFileComp
}
