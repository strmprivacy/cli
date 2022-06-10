package random_events

import (
	"errors"
	"fmt"
	"math/rand"
	"net/http"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"github.com/strmprivacy/api-definitions-go/v2/api/entities/v1"
	"strmprivacy/strm/pkg/auth"
	"strmprivacy/strm/pkg/common"
	"strmprivacy/strm/pkg/entity/stream"
	"strmprivacy/strm/pkg/simulator"
	"strmprivacy/strm/pkg/util"
)

// start a random simulator
func run(cmd *cobra.Command, streamName *string) {
	s := &entities.Stream{Ref: &entities.StreamRef{
		BillingId: auth.Auth.BillingId(),
		ProjectId: common.ProjectId,
		Name: *streamName,
	}}
	flags := cmd.Flags()
	// loads Stream definition from save version
	if err := util.TryLoad(s, streamName); err != nil {
		// there was no saved version, try to get credentials from the command options
		clientId := util.GetStringAndErr(flags, common.ClientIdFlag)
		clientSecret := util.GetStringAndErr(flags, common.ClientSecretFlag)
		if len(clientId) == 0 || len(clientSecret) == 0 {
			common.CliExit(errors.New(fmt.Sprintf("There are no credentials stored for stream '%s'", *streamName)))
		}
		s.Credentials = append(s.Credentials, &entities.Credentials{
			ClientSecret: clientSecret, ClientId: clientId,
		})
	}
	streamInfo := stream.Get(streamName, false)
	if len(streamInfo.StreamTree.Stream.LinkedStream) != 0 {
		common.CliExit(errors.New("You can't run a simulator on a derived stream"))
	}
	interval := time.Duration(util.GetIntAndErr(flags, sim.IntervalFlag))
	sessionRange := util.GetIntAndErr(flags, sim.SessionRangeFlag)
	sessionPrefix := util.GetStringAndErr(flags, sim.SessionPrefixFlag)
	gateway := util.GetStringAndErr(flags, sim.EventsApiUrlFlag)
	quiet := util.GetBoolAndErr(flags, sim.QuietFlag)
	consentLevels, err := flags.GetStringSlice(sim.ConsentLevelsFlag)
	common.CliExit(err)

	schema := util.GetStringAndErr(flags, sim.SchemaFlag)
	f := EventGenerators[schema]
	if f == nil {
		common.CliExit(errors.New(fmt.Sprintf("Can't simulate for schema %s", schema)))
	}

	if len(consentLevels) == 0 {
		common.CliExit(errors.New(fmt.Sprintf("%v is not a valid set of consent levels", consentLevels)))
	}

	if !quiet {
		fmt.Printf("Starting to simulate random %s events to stream %s. ",
			schema, *streamName)
		fmt.Printf("Sending one event every %d ms.\n", interval)
	}

	client := http.Client{}
	var sender sim.Sender
	sender = sim.ModernSender{Client: client, Gateway: gateway, Schema: schema}

	var ct = 0
	now := time.Now()
	for {
		sessionId := fmt.Sprintf("%s-%d", sessionPrefix, rand.Intn(sessionRange))
		event := f(randomConsentLevels(consentLevels), sessionId)
		token := auth.GetEventToken(s.Ref.BillingId, s.Credentials[0].ClientId, s.Credentials[0].ClientSecret)

		go sender.Send(event, token)
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
