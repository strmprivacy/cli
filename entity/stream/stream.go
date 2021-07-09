package stream

import (
	"context"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/streammachineio/api-definitions-go/api/entities/v1"
	"github.com/streammachineio/api-definitions-go/api/streams/v1"
	"google.golang.org/grpc"
	"log"
	"streammachine.io/strm/utils"
)

// strings used in the cli
const (
	linkedStreamFlag     = "derived-from"
	consentLevelTypeFlag = "consent-type"
	consentLevelsFlag    = "levels"
	tagsFlag             = "tags"
	descriptionFlag      = "description"
	saveFlag             = "save"
)

var BillingId string
var client streams.StreamsServiceClient
var apiContext context.Context

func SetupClient(clientConnection *grpc.ClientConn, ctx context.Context) {
	apiContext = ctx
	client = streams.NewStreamsServiceClient(clientConnection)
}

func Get1(streamName *string, recursive bool) *streams.GetStreamResponse {
	req := &streams.GetStreamRequest{Recursive: recursive, Ref: Ref(streamName)}
	stream, err := client.GetStream(apiContext, req)
	cobra.CheckErr(err)
	return stream
}

func list(recursive bool) {
	req := &streams.ListStreamsRequest{BillingId: BillingId, Recursive: recursive}
	streamsList, err := client.ListStreams(apiContext, req)
	cobra.CheckErr(err)
	utils.Print(streamsList)
}

func get(streamName *string, recursive bool) {
	stream := Get1(streamName, recursive)
	utils.Print(stream)
}

func del(streamName *string, recursive bool) {
	stream := Get1(streamName, recursive)
	req := &streams.DeleteStreamRequest{Recursive: recursive, Ref: Ref(streamName)}
	_, err := client.DeleteStream(apiContext, req)
	cobra.CheckErr(err)
	utils.Print(stream)
}

func create(args []string, cmd *cobra.Command) {
	var err error
	flags := cmd.Flags()
	linkedStream := utils.GetStringAndErr(flags, linkedStreamFlag)
	stream := &entities.Stream{Ref: &entities.StreamRef{BillingId: BillingId}}
	if len(args) > 0 {
		stream.Ref.Name = args[0]
	}

	if len(linkedStream) != 0 {
		stream.ConsentLevels, err = flags.GetInt32Slice(consentLevelsFlag)
		if len(stream.ConsentLevels) == 0 {
			log.Fatalf("You need consent levels when creating a derived stream")
		}
		stream.ConsentLevelType, err = parseConsentLevelType(flags)
		stream.LinkedStream = linkedStream
	} else {
		if len(stream.Ref.Name) == 0 {
			log.Fatalf("You must provide a name when creating a source stream")
		}
	}

	stream.Description = utils.GetStringAndErr(flags, descriptionFlag)
	stream.Tags, err = flags.GetStringSlice(tagsFlag)
	cobra.CheckErr(err)
	req := &streams.CreateStreamRequest{Stream: stream}
	response, err := client.CreateStream(apiContext, req)
	cobra.CheckErr(err)
	utils.Print(response.Stream)
	save, err := flags.GetBool(saveFlag)
	if save {
		utils.Save(response.Stream, &response.Stream.Ref.Name)
	}
}

func parseConsentLevelType(flags *pflag.FlagSet) (entities.ConsentLevelType, error) {
	var err error
	var consentLevelTypeString string
	consentLevelTypeString = utils.GetStringAndErr(flags, consentLevelTypeFlag)
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

func ExistingNamesCompletion(cmd *cobra.Command, args []string, complete string) ([]string, cobra.ShellCompDirective) {
	if len(args) != 0 {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}
	req := &streams.ListStreamsRequest{BillingId: BillingId}
	response, err := client.ListStreams(apiContext, req)
	cobra.CheckErr(err)
	streamResponse := response.Streams
	names := make([]string, 0, len(streamResponse))
	for _, s := range streamResponse {
		names = append(names, s.Stream.Ref.Name)
	}
	return names, cobra.ShellCompDirectiveNoFileComp
}

func ExistingSourceStreamNames(cmd *cobra.Command, args []string, complete string) ([]string, cobra.ShellCompDirective) {
	req := &streams.ListStreamsRequest{BillingId: BillingId}
	response, err := client.ListStreams(apiContext, req)
	cobra.CheckErr(err)
	streamResponse := response.Streams
	names := make([]string, 0, len(streamResponse))
	for _, s := range streamResponse {
		if len(s.Stream.LinkedStream) == 0 {
			names = append(names, s.Stream.Ref.Name)
		}
	}
	return names, cobra.ShellCompDirectiveNoFileComp
}

func Ref(n *string) *entities.StreamRef { return &entities.StreamRef{BillingId: BillingId, Name: *n} }
