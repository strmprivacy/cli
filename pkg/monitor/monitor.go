package monitor

import (
	"context"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/strmprivacy/api-definitions-go/v3/api/monitoring/v1"
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
	entityName := ""
	if len(args) > 0 {
		entityName = args[0]
	}
	projectId := project.GetProjectId(cmd)
	follow := util.GetBoolAndErr(flags, followFlag)
	entityRef := &monitoring.EntityState_Ref{
		ProjectId: projectId,
		Type:      entityType,
		Name:      entityName,
	}
	mask := &field_mask.FieldMask{
		Paths: []string{"state_time", "ref", "status", "message", "resource_type"},
	}
	if follow {
		monitorFollow(entityRef, mask)
	} else {
		monitorGetLatest(entityRef, mask)
	}
}

func monitorGetLatest(ref *monitoring.EntityState_Ref, mask *field_mask.FieldMask) {
	request := &monitoring.GetLatestEntityStatesRequest{
		Ref:            ref,
		ProjectionMask: mask,
	}
	resp, err := client.GetLatestEntityStates(apiContext, request)
	common.CliExit(err)
	printer.Print(resp)
}

func monitorFollow(ref *monitoring.EntityState_Ref, mask *field_mask.FieldMask) {
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
