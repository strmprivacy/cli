package usage

import (
	"context"
	"encoding/csv"
	"fmt"
	"github.com/golang/protobuf/ptypes/duration"
	"github.com/golang/protobuf/ptypes/timestamp"
	"github.com/spf13/cobra"
	"github.com/streammachineio/api-definitions-go/api/entities/v1"
	"github.com/streammachineio/api-definitions-go/api/usage/v1"
	"google.golang.org/grpc"
	"os"
	"strconv"
	"streammachine.io/strm/pkg/common"
	"streammachine.io/strm/pkg/util"
	"time"
)

var client usage.UsageServiceClient
var apiContext context.Context


func SetupClient(clientConnection *grpc.ClientConn, ctx context.Context) {
	apiContext = ctx
	client = usage.NewUsageServiceClient(clientConnection)
}

func get(cmd *cobra.Command, streamName *string) {
	flags := cmd.Flags()
	from := util.GetStringAndErr(flags, fromFlag)
	until := util.GetStringAndErr(flags, untilFlag)
	by := util.GetStringAndErr(flags, aggregateBy)

	startTime := time.Now().Add(time.Second * -86400)
	endTime := time.Now()
	interval := int64(300)
	var err error

	if len(by)!=0 {
		interval, err = strconv.ParseInt(by, 10, 64) ; if err != nil {
			common.CliExit(err)
		}
	}


	if len(from)!=0 {
		startTime, _ = time.Parse(time.RFC3339, from )
	}
	if len(until)!=0 {
		endTime, _ = time.Parse(time.RFC3339, from )
	}

	req := &usage.GetStreamEventUsageRequest{Ref: &entities.StreamRef{
		BillingId: common.BillingId,
		Name:      *streamName,
	},
		StartTime : &timestamp.Timestamp{Seconds: startTime.Unix()},
		EndTime : &timestamp.Timestamp{Seconds: endTime.Unix()},
		Interval: &duration.Duration{Seconds: interval},
	}

	streamUsage, err := client.GetStreamEventUsage(apiContext, req)
	common.CliExit(err)
	if util.GetBoolAndErr(flags, jsonFlag) {
		util.Print(streamUsage)

	} else {
		w := csv.NewWriter(os.Stdout)
		_ = w.Write([]string{"From", "Duration", "Count", "Change"})

		windowCount := int64(-1)
		change := ""
		for _, window := range streamUsage.Windows {
			if windowCount != -1 {
				change = fmt.Sprintf("%d", window.EventCount -windowCount)
			}
			windowCount = window.EventCount

			windowDuration := window.EndTime.AsTime().Sub(window.StartTime.AsTime())
			record := []string{isoFormat(window.StartTime.AsTime()),
				fmt.Sprintf("%.0f", windowDuration.Seconds()),
				strconv.FormatInt(windowCount, 10),
				change,
			}
			_ = w.Write(record)
		}
		w.Flush()
	}
}

func isoFormat(t time.Time) string {
	return t.Format(time.RFC3339)
}