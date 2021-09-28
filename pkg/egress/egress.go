package egress

import (
	"fmt"
	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/streammachineio/api-definitions-go/api/entities/v1"
	"net/http"
	"streammachine.io/strm/pkg/auth"
	"streammachine.io/strm/pkg/common"
	"streammachine.io/strm/pkg/sim"
	"streammachine.io/strm/pkg/util"
)

const (
	UrlFlag = "egress"
)

func Run(cmd *cobra.Command, streamName *string) {
	s := &entities.Stream{
		Ref: &entities.StreamRef{BillingId: auth.Auth.BillingId(), Name: *streamName},
	}
	flags := cmd.Flags()
	u := util.GetStringAndErr(flags, UrlFlag)
	// loads Stream definition from save version
	if err := util.TryLoad(s, streamName); err != nil {
		// there was no saved version, try to get credentials from the command options
		clientId := util.GetStringAndErr(flags, sim.ClientIdFlag)
		clientSecret := util.GetStringAndErr(flags, sim.ClientSecretFlag)
		if len(clientId) == 0 || len(clientSecret) == 0 {
			log.Fatalf("There's no saved stream and clientId %s clientSecret %s are missing",
				clientId, clientSecret)
		}
		s.Credentials = append(s.Credentials, &entities.Credentials{
			ClientSecret: clientSecret, ClientId: clientId,
		})
	}
	sts := util.GetStringAndErr(flags, auth.EventAuthHostFlag)
	authClient := &auth.Authenticator{Uri: sts}
	authClient.AuthenticateEvent(s.Ref.BillingId, s.Credentials[0].ClientId, s.Credentials[0].ClientSecret)

	token := authClient.GetToken()
	header := http.Header{"authorization": []string{"Bearer " + *token}}
	c, _, err := websocket.DefaultDialer.Dial(u, header)
	common.CliExit(err)

	for {
		_, message, err := c.ReadMessage()
		common.CliExit(err)
		fmt.Println(string(message))
	}
}