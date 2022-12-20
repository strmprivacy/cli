package monitor

import (
	"context"
	"fmt"
	"github.com/bykof/gostradamus"
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
var tz = gostradamus.Local()

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

/*
Call CRS to get the latest entity states of all entities defined by ref.
*/
func monitorGetLatest(ref *monitoring.EntityState_Ref, mask *field_mask.FieldMask) {
	request := &monitoring.GetLatestEntityStateRequest{
		Ref:            ref,
		ProjectionMask: mask,
	}
	_, err := client.GetLatestEntityState(apiContext, request)
	common.CliExit(err)
	// TODO create printer and show response
}

/*
Call CRS to get a stream of entity states of all entities defined by ref.
*/
func monitorFollow(ref *monitoring.EntityState_Ref, mask *field_mask.FieldMask) {
	in := &monitoring.GetEntityStateRequest{
		Ref:            ref,
		ProjectionMask: mask,
	}
	stream, err := client.GetEntityState(apiContext, in)
	common.CliExit(err)

	done := make(chan bool)

	// TODO create printer
	go func() {
		for {
			resp, err := stream.Recv()
			if err == io.EOF {
				done <- true //close(done)
				return
			}
			if err != nil {
				fmt.Fprintf(os.Stderr, "can not receive %v", err)
				done <- true //close(done)
				return
			}
			fmt.Printf("%s %v %s %s %s %s\n",
				util.IsoFormat(tz, resp.State.StateTime),
				resp.State.Ref.Type,
				resp.State.Ref.Name,
				resp.State.Status,
				resp.State.ResourceType,
				resp.State.Message,
			)
		}
	}()
	<-done
}
