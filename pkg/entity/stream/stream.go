package stream

import (
	"context"
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/streammachineio/api-definitions-go/api/entities/v1"
	"github.com/streammachineio/api-definitions-go/api/streams/v1"
	"google.golang.org/grpc"
	"streammachine.io/strm/pkg/auth"
	"streammachine.io/strm/pkg/common"
	"streammachine.io/strm/pkg/util"
	"strings"
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

var client streams.StreamsServiceClient
var apiContext context.Context

func SetupClient(clientConnection *grpc.ClientConn, ctx context.Context) {
	apiContext = ctx
	client = streams.NewStreamsServiceClient(clientConnection)
}

func Get(streamName *string, recursive bool) *streams.GetStreamResponse {
	if len(strings.TrimSpace(auth.Auth.BillingId())) == 0 {
		common.CliExit(fmt.Sprintf("No login information found. Use: `%v auth login` first.", common.RootCommandName))
	}

	req := &streams.GetStreamRequest{
		Recursive: recursive,
		Ref:       &entities.StreamRef{BillingId: auth.Auth.BillingId(), Name: *streamName},
	}
	stream, err := client.GetStream(apiContext, req)
	common.CliExit(err)
	return stream
}

func list(recursive bool) {
	req := &streams.ListStreamsRequest{BillingId: auth.Auth.BillingId(), Recursive: recursive}
	response, err := client.ListStreams(apiContext, req)
	common.CliExit(err)
	printer.Print(response)
}

func get(streamName *string, recursive bool) {
	response := Get(streamName, recursive)
	printer.Print(response)
}

func del(streamName *string, recursive bool) {
	response := Get(streamName, recursive)
	req := &streams.DeleteStreamRequest{
		Recursive: recursive, Ref: &entities.StreamRef{BillingId: auth.Auth.BillingId(), Name: *streamName},
	}
	_, err := client.DeleteStream(apiContext, req)
	common.CliExit(err)
	printer.Print(response)
}

func create(args []string, cmd *cobra.Command) {
	var err error
	flags := cmd.Flags()
	linkedStream := util.GetStringAndErr(flags, linkedStreamFlag)
	stream := &entities.Stream{Ref: &entities.StreamRef{BillingId: auth.Auth.BillingId()}}
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

	stream.Description = util.GetStringAndErr(flags, descriptionFlag)
	stream.Tags, err = flags.GetStringSlice(tagsFlag)
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
	if auth.Auth.BillingIdAbsent() {
		return common.MissingBillingIdCompletionError(cmd.CommandPath())
	}
	if len(args) != 0 {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}

	req := &streams.ListStreamsRequest{BillingId: auth.Auth.BillingId()}
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
	if auth.Auth.BillingIdAbsent() {
		return common.MissingBillingIdCompletionError(cmd.CommandPath())
	}
	if len(args) != 0 {
		// this one means you don't get two completion suggestions for one stream
		return nil, cobra.ShellCompDirectiveNoFileComp
	}

	req := &streams.ListStreamsRequest{BillingId: auth.Auth.BillingId()}
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
