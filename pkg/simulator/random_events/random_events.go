package random_events

import (
	"context"
	"errors"
	"fmt"
	"golang.org/x/oauth2/clientcredentials"
	"math/rand"
	"net/http"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"github.com/strmprivacy/api-definitions-go/v2/api/entities/v1"
	"strmprivacy/strm/pkg/common"
	"strmprivacy/strm/pkg/entity/stream"
	"strmprivacy/strm/pkg/util"
)

const (
	IntervalFlag      = "interval"
	EventsApiUrlFlag  = "events-api-url"
	SessionRangeFlag  = "session-range"
	SessionPrefixFlag = "session-prefix"
	purposesFlag      = "purposes"
	QuietFlag         = "quiet"
	SchemaFlag        = "schema"
)

// start a random simulator
func randomEvents(cmd *cobra.Command, streamName *string) {
	s := getStream(streamName, cmd)

	flags := cmd.Flags()
	interval := time.Duration(util.GetIntAndErr(flags, IntervalFlag))
	sessionRange := util.GetIntAndErr(flags, SessionRangeFlag)
	sessionPrefix := util.GetStringAndErr(flags, SessionPrefixFlag)
	gateway := util.GetStringAndErr(flags, EventsApiUrlFlag)
	quiet := util.GetBoolAndErr(flags, QuietFlag)
	purposeConsentLevels, err := flags.GetStringSlice(purposesFlag)
	common.CliExit(err)

	schema := util.GetStringAndErr(flags, SchemaFlag)

	if EventGenerators[schema] == nil {
		common.CliExit(errors.New(fmt.Sprintf("The schema %s is not supported in the CLI", schema)))
	}

	if len(purposeConsentLevels) == 0 {
		common.CliExit(errors.New(fmt.Sprintf("%v is not a valid set of purpose levels", purposeConsentLevels)))
	}

	if !quiet {
		fmt.Printf("Starting to simulate random %s events to stream %s. ",
			schema, *streamName)
		fmt.Printf("Sending one event every %d ms.\n", interval)
	}

	simulate(s, gateway, schema, sessionPrefix, sessionRange, purposeConsentLevels, interval, quiet)
}

func simulate(s *entities.Stream, gateway string, schema string, sessionPrefix string, sessionRange int, consentLevels []string, interval time.Duration, quiet bool) {
	client := authenticatedHttpClient(s.Credentials[0])
	f := EventGenerators[schema]

	simulator := Simulator{Client: client, Gateway: gateway, Schema: schema}

	var ct = 0
	now := time.Now()
	for {
		sessionId := fmt.Sprintf("%s-%d", sessionPrefix, rand.Intn(sessionRange))
		event := f(randomConsentLevels(consentLevels), sessionId)

		go simulator.Send(event)
		ct += 1
		time.Sleep(interval * time.Millisecond)
		if !quiet && time.Now().Sub(now) > 5*time.Second {
			println("Sent", ct, "events")
			now = time.Now()
		}
	}
}

// randomConsentLevels returns a slice of integers for the simulated event.
// It starts with [ "", "0", "0/1", "3/8", "3/7/10", ...] so a slice of strings that define
// what we want to send. This method picks a random one.
func randomConsentLevels(levels []string) []int32 {
	l := strings.Split(levels[rand.Intn(len(levels))], "/")
	if len(l) == 1 && l[0] == "" {
		// No consent at all
		return []int32{}
	}
	return util.StringsArrayToInt32(l)
}

func getStream(streamName *string, cmd *cobra.Command) *entities.Stream {
	s := stream.Get(streamName, false, cmd).StreamTree.Stream

	if len(s.LinkedStream) != 0 {
		common.CliExit(errors.New("simulators cannot be randomEvents on a derived stream"))
	}

	return s
}

func authenticatedHttpClient(credentials *entities.Credentials) http.Client {
	config := clientcredentials.Config{
		ClientID:     credentials.ClientId,
		ClientSecret: credentials.ClientSecret,
		TokenURL:     fmt.Sprintf("%v/auth/realms/streams/protocol/openid-connect/token", common.ApiAuthHost),
		Scopes:       []string{"offline_access"},
	}

	return *config.Client(context.Background())
}
