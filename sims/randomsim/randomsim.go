package randomsim

import (
	"bytes"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/streammachineio/api-definitions-go/api/entities/v1"
	"io"
	"log"
	"math/rand"
	"net/http"
	"streammachine.io/strm/auth"
	"streammachine.io/strm/clickstream"
	"streammachine.io/strm/entity/stream"
	"streammachine.io/strm/sims"
	"streammachine.io/strm/utils"
	"strings"
	"time"
)

// start a random simulator
func run(cmd *cobra.Command, streamName *string) {
	s := &entities.Stream{Ref: &entities.StreamRef{BillingId: sims.BillingId, Name: *streamName}}
	flags := cmd.Flags()
	// loads Stream definition from save version
	if err := utils.TryLoad(s, streamName); err != nil {
		// there was no saved version, try to get credentials from the command options
		clientId := utils.GetStringAndErr(flags, sims.ClientIdFlag)
		clientSecret := utils.GetStringAndErr(flags, sims.ClientSecretFlag)
		if len(clientId) == 0 || len(clientSecret) == 0 {
			log.Fatalf("There are no credentials stored for stream '%s'", *streamName)
		}
		s.Credentials = append(s.Credentials, &entities.Credentials{
			ClientSecret: clientSecret, ClientId: clientId,
		})
	}
	streamInfo := stream.Get1(streamName, false)
	if len(streamInfo.StreamTree.Stream.LinkedStream) != 0 {
		log.Fatalf("You can't run a simulator on a derived stream")
	}
	interval := time.Duration(utils.GetIntAndErr(flags, sims.IntervalFlag))
	sts := utils.GetStringAndErr(flags, auth.EventAuthHostFlag)
	sessionRange := utils.GetIntAndErr(flags, sims.SessionRangeFlag)
	sessionPrefix := utils.GetStringAndErr(flags, sims.SessionPrefixFlag)
	gateway := utils.GetStringAndErr(flags, sims.EventGatewayFlag)
	quiet := utils.GetBoolAndErr(flags, sims.QuietFlag)
	consentLevels, err := flags.GetStringSlice(sims.ConsentLevelsFlag)
	cobra.CheckErr(err)
	if len(consentLevels) == 0 {
		log.Fatalf("%v is not a valid set of consent levels", consentLevels)
	}
	authClient := &auth.Auth{Uri: sts}
	authClient.AuthenticateEvent(s.Ref.BillingId, s.Credentials[0].ClientId, s.Credentials[0].ClientSecret)
	if !quiet {

		println("Starting sim to stream '"+*streamName+"'. Sending 1 event every", interval, "ms")
	}

	client := http.Client{}
	var ct = 0
	now := time.Now()
	for {
		event := clickstream.NewClickstreamEvent()
		event.StrmMeta = &clickstream.StrmMeta{
			ConsentLevels: randomConsentLevels(consentLevels),
		}
		event.ProducerSessionId = fmt.Sprintf("%s-%d", sessionPrefix, rand.Intn(sessionRange))
		event.Customer = &clickstream.Customer{Id: "customer-" + event.ProducerSessionId}
		event.Url = "https://www.streammachine.io/rules"
		token, _ := authClient.GetToken(quiet)
		go sendEvent(client, event, gateway, token)
		ct += 1
		time.Sleep(interval * time.Millisecond)
		if !quiet && time.Now().Sub(now) > 5*time.Second {
			println("Sent", ct, "events")
			now = time.Now()
		}
	}
}

// randomConsentLevels returns a slice of integers for the simulated event.
// It starts with [ "0", "0/1", "3/8", "3/7/10", ...] so a slice of strings that define
// what we want to send. This method picks a random one.
func randomConsentLevels(levels []string) []int32 {
	l := strings.Split(levels[rand.Intn(len(levels))], "/")
	return utils.StringsArrayToInt32(l)
}

func sendEvent(client http.Client, event *clickstream.ClickstreamEvent,
	gateway string, token string) {
	b := &bytes.Buffer{}
	err := event.Serialize(b)
	cobra.CheckErr(err)
	req, err := http.NewRequest("POST", gateway, b)
	cobra.CheckErr(err)
	req.Header.Set("Strm-Schema-Id", "clickstream")
	req.Header.Set("Authorization", "Bearer "+token)
	resp, err := client.Do(req)
	if err != nil || resp.StatusCode != 204 {
		if resp != nil {
			defer resp.Body.Close()
			body, err := io.ReadAll(resp.Body)
			fmt.Printf("%v %s\n", err, string(body))
		}
	}
}
