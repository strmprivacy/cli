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

	//entityName := ""
	//if len(args)>0 {
	//	entityName = args[0]
	// figure out entity Id from name!
	//}

	in := &monitoring.GetEntityStateRequest{
		Ref: &monitoring.EntityState_Ref{
			// TODO Robin setting these values crash the CRS with a NPE.
			// but this is my old CRS running.
			// I don't want to touch the CRS because you're in the middle of a
			// major change
			// ProjectId: project.GetProjectId(cmd),
			// Type:      entityType,
			// Id: entityId,
		},
		ProjectionMask: &field_mask.FieldMask{
			Paths: []string{"state_time", "ref", "status", "message", "resource_type"},
		},
	}
	stream, err := client.GetEntityState(apiContext, in)
	common.CliExit(err)

	done := make(chan bool)

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
				resp.State.Ref.Id,
				resp.State.Status,
				resp.State.ResourceType,
				resp.State.Message,
			)
		}
	}()
	<-done

}
