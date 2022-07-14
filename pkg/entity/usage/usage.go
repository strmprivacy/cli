package usage

import (
	"context"
	"errors"
	"fmt"
	"github.com/bykof/gostradamus"
	"github.com/golang/protobuf/ptypes/duration"
	"github.com/golang/protobuf/ptypes/timestamp"
	"github.com/spf13/cobra"
	"github.com/strmprivacy/api-definitions-go/v2/api/entities/v1"
	"github.com/strmprivacy/api-definitions-go/v2/api/usage/v1"
	"google.golang.org/grpc"
	"regexp"
	"strconv"
	"strings"
	"strmprivacy/strm/pkg/common"
	"strmprivacy/strm/pkg/util"
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
		ProjectId: common.ProjectId,
		Name:      *streamName,
	},
		StartTime: &timestamp.Timestamp{Seconds: startTime.UnixTimestamp()},
		EndTime:   &timestamp.Timestamp{Seconds: endTime.UnixTimestamp()},
		Interval:  &duration.Duration{Seconds: interval},
	}

	response, err := client.GetStreamEventUsage(apiContext, req)
	common.CliExit(err)

	printer.Print(response)
}

func interpretInterval(by string) int64 {
	interval, err := strconv.ParseInt(by, 10, 64)
	if err != nil {
		// the interval is not an integer (seconds).
		// let's interpret it as an integer followed by (m)inute, (h)our, (s)econd, (d)ay
		pat := regexp.MustCompilePOSIX(`^([0-9]+)([mshd])$`)
		r := pat.FindAllStringSubmatch(strings.ToLower(by), -1)
		if r == nil {
			common.CliExit(errors.New(fmt.Sprintf("%v not understood as interval format", by)))
		}
		lookup := map[string]int64{
			"m": 60,
			"h": 3600,
			"d": 86400,
			"s": 1,
		}
		v, ok := lookup[r[0][2]]
		if !ok {
			common.CliExit(errors.New("Don't understand unit " + r[0][2]))
		}
		interval, err = strconv.ParseInt(r[0][1], 10, 64)
		common.CliExit(err)

		interval *= v
	}
	return interval
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
