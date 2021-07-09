package egress

import (
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/spf13/cobra"
	"github.com/streammachineio/api-definitions-go/api/entities/v1"
	"log"
	"net/http"
	"streammachine.io/strm/auth"
	"streammachine.io/strm/sims"
	"streammachine.io/strm/utils"
)

const (
	UrlFlag = "egress"
)

var BillingId string

func Run(cmd *cobra.Command, streamName *string) {
	s := &entities.Stream{Ref: &entities.StreamRef{BillingId: sims.BillingId, Name: *streamName}}
	flags := cmd.Flags()
	u := utils.GetStringAndErr(flags, UrlFlag)
	// loads Stream definition from save version
	if err := utils.TryLoad(s, streamName); err != nil {
		// there was no saved version, try to get credentials from the command options
		clientId := utils.GetStringAndErr(flags, sims.ClientIdFlag)
		clientSecret := utils.GetStringAndErr(flags, sims.ClientSecretFlag)
		if len(clientId) == 0 || len(clientSecret) == 0 {
			log.Fatalf("There's no saved stream and clientId %s clientSecret %s are missing",
				clientId, clientSecret)
		}
		s.Credentials = append(s.Credentials, &entities.Credentials{
			ClientSecret: clientSecret, ClientId: clientId,
		})
	}
	sts := utils.GetStringAndErr(flags, auth.EventAuthHostFlag)
	authClient := &auth.Auth{Uri: sts}
	authClient.AuthenticateEvent(s.Ref.BillingId, s.Credentials[0].ClientId, s.Credentials[0].ClientSecret)

	token, _ := authClient.GetToken(false)
	header := http.Header{"authorization": []string{"Bearer " + token}}
	c, _, err := websocket.DefaultDialer.Dial(u, header)
	cobra.CheckErr(err)

	for {
		_, message, err := c.ReadMessage()
		cobra.CheckErr(err)
		fmt.Println(string(message))
	}
}