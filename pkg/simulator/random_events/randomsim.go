package random_events

import (
	"fmt"
	"math/rand"
	"net/http"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/streammachineio/api-definitions-go/api/entities/v1"
	"streammachine.io/strm/pkg/auth"
	"streammachine.io/strm/pkg/common"
	"streammachine.io/strm/pkg/entity/stream"
	"streammachine.io/strm/pkg/simulator"
	"streammachine.io/strm/pkg/util"
)

// start a random simulator
func run(cmd *cobra.Command, streamName *string) {
	s := &entities.Stream{Ref: &entities.StreamRef{BillingId: auth.Auth.BillingId(), Name: *streamName}}
	flags := cmd.Flags()
	// loads Stream definition from save version
	if err := util.TryLoad(s, streamName); err != nil {
		// there was no saved version, try to get credentials from the command options
		clientId := util.GetStringAndErr(flags, sim.ClientIdFlag)
		clientSecret := util.GetStringAndErr(flags, sim.ClientSecretFlag)
		if len(clientId) == 0 || len(clientSecret) == 0 {
			log.Fatalf("There are no credentials stored for stream '%s'", *streamName)
		}
		s.Credentials = append(s.Credentials, &entities.Credentials{
			ClientSecret: clientSecret, ClientId: clientId,
		})
	}
	streamInfo := stream.Get(streamName, false)
	if len(streamInfo.StreamTree.Stream.LinkedStream) != 0 {
		log.Fatalf("You can't run a simulator on a derived stream")
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
		common.CliExit(fmt.Sprintf("Can't simulate for schema %s", schema))
	}

	if len(consentLevels) == 0 {
		log.Fatalf("%v is not a valid set of consent levels", consentLevels)
	}

	if !quiet {
		fmt.Printf("Starting to simulate random %s events to stream %s. ",
			schema, *streamName)
		fmt.Printf("Sending one event every %d ms.\n", interval)
	}

	client := http.Client{}
	var sender sim.Sender
	if schema == "clickstream" {
		sender = sim.LegacySender{Client: client, Gateway: gateway, Schema: schema}
	} else {
		sender = sim.ModernSender{Client: client, Gateway: gateway, Schema: schema}
	}

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