package stream

import (
	"context"
	"errors"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/strmprivacy/api-definitions-go/v2/api/entities/v1"
	"github.com/strmprivacy/api-definitions-go/v2/api/streams/v1"
	"google.golang.org/grpc"
	"strings"
	"strmprivacy/strm/pkg/common"
	"strmprivacy/strm/pkg/entity/policy"
	"strmprivacy/strm/pkg/entity/project"
	"strmprivacy/strm/pkg/util"
)

// strings used in the cli
const (
	linkedStreamFlag     = "derived-from"
	consentLevelTypeFlag = "consent-type"
	consentLevelsFlag    = "levels"
	tagsFlag             = "tags"
	descriptionFlag      = "description"
	saveFlag             = "save"
	maskedFieldsFlag     = "masked-fields"
	maskedFieldsSeed     = "mask-seed"
	maskedFieldHelp      = `-M strmprivacy/example/1.3.0:sensitiveValue,consistentValue \
-M strmprivacy/clickstream/1.0.0:sessionId

Masks fields values in the output stream via hashing.
	`
)

var client streams.StreamsServiceClient
var apiContext context.Context

func SetupClient(clientConnection *grpc.ClientConn, ctx context.Context) {
	apiContext = ctx
	client = streams.NewStreamsServiceClient(clientConnection)
}

func Get(streamName *string, recursive bool, cmd *cobra.Command) *streams.GetStreamResponse {
	req := &streams.GetStreamRequest{
		Recursive: recursive,
		Ref: &entities.StreamRef{
			ProjectId: project.GetProjectId(cmd),
			Name:      *streamName,
		},
	}
	stream, err := client.GetStream(apiContext, req)
	common.CliExit(err)
	return stream
}

func list(recursive bool, cmd *cobra.Command) {
	req := &streams.ListStreamsRequest{
		ProjectId: project.GetProjectId(cmd),
		Recursive: recursive,
	}
	response, err := client.ListStreams(apiContext, req)
	common.CliExit(err)
	printer.Print(response)
}

func get(streamName *string, recursive bool, cmd *cobra.Command) {
	response := Get(streamName, recursive, cmd)
	printer.Print(response)
}

func del(streamName *string, recursive bool, cmd *cobra.Command) {
	response := Get(streamName, recursive, cmd)
	req := &streams.DeleteStreamRequest{
		Recursive: recursive, Ref: &entities.StreamRef{
			ProjectId: project.GetProjectId(cmd),
			Name:      *streamName,
		},
	}
	_, err := client.DeleteStream(apiContext, req)
	common.CliExit(err)
	util.DeleteSaved(response.StreamTree.Stream, streamName)
	printer.Print(response)
}

func create(args []string, cmd *cobra.Command) {
	var err error
	flags := cmd.Flags()
	linkedStream := util.GetStringAndErr(flags, linkedStreamFlag)
	stream := &entities.Stream{Ref: &entities.StreamRef{
		ProjectId: project.GetProjectId(cmd),
	}}
	if len(args) > 0 {
		stream.Ref.Name = args[0]
	}

	if len(linkedStream) != 0 {
		stream.ConsentLevels, err = flags.GetInt32Slice(consentLevelsFlag)
		if len(stream.ConsentLevels) == 0 {
			common.CliExit(errors.New("You need consent levels when creating a derived stream"))
		}
		stream.ConsentLevelType, err = parseConsentLevelType(flags)
		stream.LinkedStream = linkedStream
	} else {
		if len(stream.Ref.Name) == 0 {
			common.CliExit(errors.New("You must provide a name when creating a source stream"))
		}
	}

	stream.Description = util.GetStringAndErr(flags, descriptionFlag)
	stream.Tags, err = flags.GetStringSlice(tagsFlag)
	stream.MaskedFields = parseMaskedFields(flags)
	stream.PolicyId = policy.GetPolicyFromFlags(flags)
	common.CliExit(err)

	req := &streams.CreateStreamRequest{Stream: stream}
	response, err := client.CreateStream(apiContext, req)
	common.CliExit(err)
	save, err := flags.GetBool(saveFlag)
	if save {
		util.Save(response.Stream, &response.Stream.Ref.Name)
	}

	printer.Print(response)
}

/*
			-M strmprivacy/example/1.3.0:sensitiveValue,anotherOne \
	   		-M dpg/nps_unified/v3:kiosk_v1,customer_id --masked_fields_file
*/
func parseMaskedFields(flags *pflag.FlagSet) *entities.MaskedFields {
	masked, err := flags.GetStringArray(maskedFieldsFlag)
	seed, err := flags.GetString(maskedFieldsSeed)
	common.CliExit(err)
	maskedField := &entities.MaskedFields{
		HashType:      "",
		Seed:          seed,
		FieldPatterns: map[string]*entities.MaskedFields_PatternList{},
	}
	for _, s := range masked {
		parts := strings.Split(s, ":")
		ecRef := parts[0]
		p := &entities.MaskedFields_PatternList{FieldPatterns: strings.Split(parts[1], ",")}
		maskedField.FieldPatterns[ecRef] = p
	}
	return maskedField
}

func parseConsentLevelType(flags *pflag.FlagSet) (entities.ConsentLevelType, error) {
	var err error
	var consentLevelTypeString string
	consentLevelTypeString = util.GetStringAndErr(flags, consentLevelTypeFlag)
	consentLevelType, ok := entities.ConsentLevelType_value[consentLevelTypeString]
	if !ok {
		log.Fatalf("Can't convert %s to a known consent level type, types are %v",
			consentLevelTypeString, entities.ConsentLevelType_value)
	}
	if consentLevelType == int32(entities.ConsentLevelType_CONSENT_LEVEL_TYPE_UNSPECIFIED) {
		consentLevelType = int32(entities.ConsentLevelType_CUMULATIVE)
	}
	return entities.ConsentLevelType(consentLevelType), err
}

func NamesCompletion(cmd *cobra.Command, args []string, complete string) ([]string, cobra.ShellCompDirective) {
	if len(args) != 0 && strings.Fields(cmd.Short)[0] != "Delete" {
		// this one means you don't get multiple completion suggestions for one stream if it's not a delete call
		return nil, cobra.ShellCompDirectiveNoFileComp
	}

	req := &streams.ListStreamsRequest{
		ProjectId: common.ProjectId,
	}
	response, err := client.ListStreams(apiContext, req)

	if err != nil {
		return common.GrpcRequestCompletionError(err)
	}

	names := make([]string, 0, len(response.Streams))
	for _, s := range response.Streams {
		names = append(names, s.Stream.Ref.Name)
	}

	return names, cobra.ShellCompDirectiveNoFileComp
}

func SourceNamesCompletion(cmd *cobra.Command, args []string, complete string) ([]string, cobra.ShellCompDirective) {
	req := &streams.ListStreamsRequest{
		ProjectId: common.ProjectId,
	}
	response, err := client.ListStreams(apiContext, req)

	if err != nil {
		return common.GrpcRequestCompletionError(err)
	}

	names := make([]string, 0, len(response.Streams))
	for _, s := range response.Streams {
		if len(s.Stream.LinkedStream) == 0 {
			names = append(names, s.Stream.Ref.Name)
		}
	}

	return names, cobra.ShellCompDirectiveNoFileComp
}
