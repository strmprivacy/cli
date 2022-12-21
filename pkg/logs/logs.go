package logs

import (
	"context"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/strmprivacy/api-definitions-go/v2/api/monitoring/v1"
	"google.golang.org/genproto/protobuf/field_mask"
	"google.golang.org/grpc"
	"io"
	"os"
	"strmprivacy/strm/pkg/common"
	"strmprivacy/strm/pkg/entity/project"
	"strmprivacy/strm/pkg/util"
)

var client monitoring.MonitoringServiceClient
var apiContext context.Context

func SetupClient(clientConnection *grpc.ClientConn, ctx context.Context) {
	apiContext = ctx
	client = monitoring.NewMonitoringServiceClient(clientConnection)
}

func run(cmd *cobra.Command, entityType monitoring.EntityState_EntityType, args []string) {
	flags := cmd.Flags()
	follow := util.GetBoolAndErr(flags, followFlag)
	entityRef := &monitoring.EntityState_Ref{
		ProjectId: project.GetProjectId(cmd),
		Type:      entityType,
		Name:      args[0],
	}
	mask := &field_mask.FieldMask{
		Paths: []string{"logs"},
	}
	if follow {
		logsFollow(entityRef, mask)
	} else {
		logsGetLatest(entityRef, mask)
	}
}

func logsGetLatest(ref *monitoring.EntityState_Ref, mask *field_mask.FieldMask) {
	request := &monitoring.GetLatestEntityStatesRequest{
		Ref:            ref,
		ProjectionMask: mask,
	}
	resp, err := client.GetLatestEntityStates(apiContext, request)
	common.CliExit(err)
	printer.Print(resp)
}

func logsFollow(ref *monitoring.EntityState_Ref, mask *field_mask.FieldMask) {
	in := &monitoring.GetEntityStateRequest{
		Ref:            ref,
		ProjectionMask: mask,
	}
	stream, err := client.GetEntityState(apiContext, in)
	common.CliExit(err)

	done := make(chan bool)

	go func() {
		for {
			resp, err := stream.Recv()
			if err == io.EOF {
				fmt.Fprintln(os.Stderr, "No more data to be received, closing connection.")
				done <- true //close(done)
				return
			}
			if err != nil {
				done <- true //close(done)
				common.CliExit(err)
				return
			}
			printer.Print(resp)
		}
	}()
	<-done
}
