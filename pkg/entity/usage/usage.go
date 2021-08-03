package usage

import (
	"context"
	"encoding/csv"
	"fmt"
	"github.com/bykof/gostradamus"
	"github.com/golang/protobuf/ptypes/duration"
	"github.com/golang/protobuf/ptypes/timestamp"
	"github.com/spf13/cobra"
	"github.com/streammachineio/api-definitions-go/api/entities/v1"
	"github.com/streammachineio/api-definitions-go/api/usage/v1"
	"google.golang.org/grpc"
	"math"
	"os"
	"regexp"
	"strconv"
	"streammachine.io/strm/pkg/common"
	"streammachine.io/strm/pkg/util"
	"strings"
	"time"
)

const (
	dateTimeParseFormat = "YYYY/M/D-HH:mm"
	defaultAggInterval  = int64(300)
	defaultRangeDays    = 1
)

var client usage.UsageServiceClient
var apiContext context.Context
var tz = gostradamus.Local()

func SetupClient(clientConnection *grpc.ClientConn, ctx context.Context) {
	apiContext = ctx
	client = usage.NewUsageServiceClient(clientConnection)
}

func get(cmd *cobra.Command, streamName *string) {
	flags := cmd.Flags()
	from := util.GetStringAndErr(flags, fromFlag)
	until := util.GetStringAndErr(flags, untilFlag)
	by := util.GetStringAndErr(flags, aggregateByFlag)

	// defaults
	endTime := gostradamus.NowInTimezone(tz)
	startTime := endTime.ShiftDays(-defaultRangeDays)
	interval := defaultAggInterval

	var err error

	if len(by) != 0 {
		interval = interpretInterval(by)
	}

	if len(from) != 0 {
		startTime, err = gostradamus.ParseInTimezone(from, dateTimeParseFormat, tz)
		common.CliExit(err)
	}
	if len(until) != 0 {
		endTime, err = gostradamus.ParseInTimezone(until, dateTimeParseFormat, tz)
		common.CliExit(err)
	}

	req := &usage.GetStreamEventUsageRequest{Ref: &entities.StreamRef{
		BillingId: common.BillingId,
		Name:      *streamName,
	},
		StartTime: &timestamp.Timestamp{Seconds: startTime.UnixTimestamp()},
		EndTime:   &timestamp.Timestamp{Seconds: endTime.UnixTimestamp()},
		Interval:  &duration.Duration{Seconds: interval},
	}

	streamUsage, err := client.GetStreamEventUsage(apiContext, req)
	common.CliExit(err)
	if util.GetBoolAndErr(flags, jsonFlag) {
		util.Print(streamUsage)
	} else {
		printCsv(streamUsage)
	}
}

func printCsv(streamUsage *usage.GetStreamEventUsageResponse) {
	w := csv.NewWriter(os.Stdout)
	_ = w.Write([]string{"from", "count", "duration", "change", "rate"})

	windowCount := int64(-1)
	change := math.NaN()
	for _, window := range streamUsage.Windows {
		if windowCount != -1 {
			change = float64(window.EventCount - windowCount)
		}
		windowCount = window.EventCount

		windowDuration := window.EndTime.AsTime().Sub(window.StartTime.AsTime())
		rate := change / windowDuration.Seconds()
		record := []string{isoFormat(window.StartTime.AsTime()),
			fmt.Sprintf("%d", windowCount),
			fmt.Sprintf("%.0f", windowDuration.Seconds()),
			fmt.Sprintf("%v", change),
			fmt.Sprintf("%.2f", rate),
		}
		_ = w.Write(record)
	}
	w.Flush()
}

func interpretInterval(by string) int64 {

	interval, err := strconv.ParseInt(by, 10, 64)
	if err != nil {
		// the interval is not an integer (seconds).
		// let's interpret it as an integer followed by (m)inute, (h)our, (s)econd, (d)ay
		pat := regexp.MustCompilePOSIX(`^([0-9]+)([mshd])$`)
		r := pat.FindAllStringSubmatch(strings.ToLower(by), -1)
		if r == nil {
			common.CliExit(fmt.Sprintf("%v not understood as interval format", by))
		}
		lookup := map[string]int64{
			"m": 60,
			"h": 3600,
			"d": 86400,
			"s": 1,
		}
		v, ok := lookup[r[0][2]]
		if !ok {
			common.CliExit("Don't understand unit " + r[0][2])
		}
		interval, err = strconv.ParseInt(r[0][1], 10, 64)
		if err != nil {
			common.CliExit(err)
		}
		interval *= v
	}
	return interval
}

func isoFormat(t time.Time) string {
	n := gostradamus.DateTimeFromTime(t)
	return n.InTimezone(tz).IsoFormatTZ()
	//return t.Format(time.RFC3339)
}

func dateCompletion(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	endTime := gostradamus.NowInTimezone(tz).FloorHour()
	hourRange := 48
	startTime := endTime.ShiftHours(-hourRange)
	completions := make([]string, hourRange)
	for i := range completions {
		completions[i] = startTime.ShiftHours(i).Format(dateTimeParseFormat)
	}
	return completions, cobra.ShellCompDirectiveDefault
}
