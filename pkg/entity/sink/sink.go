package sink

import (
	"context"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/streammachineio/api-definitions-go/api/entities/v1"
	"github.com/streammachineio/api-definitions-go/api/sinks/v1"
	"google.golang.org/grpc"
	"io/ioutil"
	"streammachine.io/strm/pkg/common"
	"streammachine.io/strm/pkg/util"
)

var Client sinks.SinksServiceClient
var apiContext context.Context

func ref(n *string) *entities.SinkRef {
	return &entities.SinkRef{BillingId: common.BillingId, Name: *n}
}

func SetupClient(clientConnection *grpc.ClientConn, ctx context.Context) {
	apiContext = ctx
	Client = sinks.NewSinksServiceClient(clientConnection)
}

func list(recursive bool) {
	req := &sinks.ListSinksRequest{Recursive: recursive, BillingId: common.BillingId}
	sinksList, err := Client.ListSinks(apiContext, req)
	common.CliExit(err)
	printer.Print(sinksList)
}

func get(name *string, recursive bool) {
	req := &sinks.GetSinkRequest{Recursive: recursive, Ref: ref(name)}
	stream, err := Client.GetSink(apiContext, req)
	common.CliExit(err)
	printer.Print(stream)
}

func del(name *string, recursive bool) {
	req := &sinks.DeleteSinkRequest{Recursive: recursive, Ref: ref(name)}
	sink, err := Client.DeleteSink(apiContext, req)
	common.CliExit(err)
	printer.Print(sink)
}

func create(sinkName *string, bucketName *string, cmd *cobra.Command) {
	flags := cmd.Flags()
	sink := &entities.Sink{Ref: ref(sinkName),
		Config: &entities.Sink_Bucket{
			Bucket: &entities.BucketConfig{
				BucketName:  *bucketName,
				Credentials: readCredentialsFile(flags)},
		},
		SinkType: parseSyncType(flags),
	}
	response, err := Client.CreateSink(apiContext, &sinks.CreateSinkRequest{Sink: sink})
	common.CliExit(err)
	printer.Print(response.Sink)

}

func readCredentialsFile(flags *pflag.FlagSet) string {
	fn := util.GetStringAndErr(flags, credentialsFileFlag)
	buf, err := ioutil.ReadFile(fn)
	common.CliExit(err)
	return string(buf)
}

func parseSyncType(flags *pflag.FlagSet) entities.SinkType {
	typeString := util.GetStringAndErr(flags, sinkTypeFlag)
	if len(typeString) == 0 {
		return entities.SinkType_SINK_TYPE_UNSPECIFIED
	}
	sinkType, ok := entities.SinkType_value[typeString]
	if !ok {
		log.Fatalf("Can't convert %s to a known consent sink type, types are %v",
			typeString, entities.SinkType_value)
	}
	return entities.SinkType(sinkType)
}

func NamesCompletion(cmd *cobra.Command, args []string, complete string) ([]string, cobra.ShellCompDirective) {
	if common.BillingIdIsMissing() {
		return common.MissingBillingIdCompletionError(cmd.CommandPath())
	}
	if len(args) != 0 {
		// this one means you don't get two completion suggestions for one stream
		return nil, cobra.ShellCompDirectiveNoFileComp
	}

	req := &sinks.ListSinksRequest{BillingId: common.BillingId}
	response, err := Client.ListSinks(apiContext, req)

	if err != nil {
		return common.GrpcRequestCompletionError(err)
	}

	names := make([]string, 0, len(response.Sinks))
	for _, s := range response.Sinks {
		names = append(names, s.Sink.Ref.Name)
	}

	return names, cobra.ShellCompDirectiveNoFileComp
}
