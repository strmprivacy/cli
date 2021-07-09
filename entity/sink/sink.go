package sink

import (
	"context"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/streammachineio/api-definitions-go/api/entities/v1"
	"github.com/streammachineio/api-definitions-go/api/sinks/v1"
	"google.golang.org/grpc"
	"io/ioutil"
	"log"
	"streammachine.io/strm/utils"
)

// strings used in the cli
const ()

var BillingId string
var client sinks.SinksServiceClient
var apiContext context.Context

func ref(n *string) *entities.SinkRef { return &entities.SinkRef{BillingId: BillingId, Name: *n} }

func SetupClient(clientConnection *grpc.ClientConn, ctx context.Context) {
	apiContext = ctx
	client = sinks.NewSinksServiceClient(clientConnection)
}

func ExistingNamesCompletion(cmd *cobra.Command, args []string, complete string) ([]string, cobra.ShellCompDirective) {
	if len(args) > 0 {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}
	return ExistingNames(), cobra.ShellCompDirectiveNoFileComp
}

func ExistingNames() []string {

	req := &sinks.ListSinksRequest{BillingId: BillingId}
	response, err := client.ListSinks(apiContext, req)
	cobra.CheckErr(err)
	sinkNames := make([]string, 0, len(response.Sinks))
	for _, s := range response.Sinks {
		sinkNames = append(sinkNames, s.Sink.Ref.Name)
	}
	return sinkNames
}

func list(recursive bool) {
	req := &sinks.ListSinksRequest{Recursive: recursive, BillingId: BillingId}
	sinksList, err := client.ListSinks(apiContext, req)
	cobra.CheckErr(err)
	utils.Print(sinksList)
}

func get(name *string, recursive bool) {
	req := &sinks.GetSinkRequest{Recursive: recursive, Ref: ref(name)}
	stream, err := client.GetSink(apiContext, req)
	cobra.CheckErr(err)
	utils.Print(stream)
}

func del(name *string, recursive bool) {
	req := &sinks.DeleteSinkRequest{Recursive: recursive, Ref: ref(name)}
	sink, err := client.DeleteSink(apiContext, req)
	cobra.CheckErr(err)
	utils.Print(sink)
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
	response, err := client.CreateSink(apiContext, &sinks.CreateSinkRequest{Sink: sink})
	cobra.CheckErr(err)
	utils.Print(response.Sink)

}

func readCredentialsFile(flags *pflag.FlagSet) string {
	fn := utils.GetStringAndErr(flags, credentialsFileFlag)
	buf, err := ioutil.ReadFile(fn)
	cobra.CheckErr(err)
	return string(buf)
}

func parseSyncType(flags *pflag.FlagSet) entities.SinkType {
	typeString := utils.GetStringAndErr(flags, sinkTypeFlag)
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
