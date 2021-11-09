package web_socket

import (
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/spf13/cobra"
	"github.com/streammachineio/api-definitions-go/api/entities/v1"
	"net/http"
	"streammachine.io/strm/pkg/auth"
	"streammachine.io/strm/pkg/common"
	"streammachine.io/strm/pkg/simulator"
	"streammachine.io/strm/pkg/util"
)

const (
	WebSocketUrl = "web-socket-url"
)

func Run(cmd *cobra.Command, streamName *string) {
	s := &entities.Stream{
		Ref: &entities.StreamRef{BillingId: auth.Auth.BillingId(), Name: *streamName},
	}
	flags := cmd.Flags()
	u := util.GetStringAndErr(flags, WebSocketUrl)
	// loads Stream definition from save version
	if err := util.TryLoad(s, streamName); err != nil {
		// there was no saved version, try to get credentials from the command options
		clientId := util.GetStringAndErr(flags, sim.ClientIdFlag)
		clientSecret := util.GetStringAndErr(flags, sim.ClientSecretFlag)
		if len(clientId) == 0 || len(clientSecret) == 0 {
			common.CliExit(fmt.Sprintf("There's no saved stream for %s and clientId %s clientSecret %s are missing as options",
				*streamName, clientId, clientSecret))
		}
		s.Credentials = append(s.Credentials, &entities.Credentials{
			ClientSecret: clientSecret, ClientId: clientId,
		})
	}
	for {

		token := auth.GetEventToken(s.Ref.BillingId, s.Credentials[0].ClientId, s.Credentials[0].ClientSecret)
		header := http.Header{"authorization": []string{"Bearer " + token}}
		ws, _, err := websocket.DefaultDialer.Dial(u, header)
		common.CliExit(err)

	innerLoop:
		for {
			_, message, err := ws.ReadMessage()
			if err == nil {
				fmt.Println(string(message))
			} else {
				// there was an error. Check it for normal websocket timeouts
				// and in that case just re-authenticate and reconnect
				if ce, ok := err.(*websocket.CloseError); ok {
					switch ce.Code {
					case websocket.CloseNormalClosure,
						websocket.CloseGoingAway,
						websocket.CloseNoStatusReceived:
						break innerLoop
					default:
						// not one of the above errors
						common.CliExit(err)
					}
				} else {
					// not a websocket.CloseError
					common.CliExit(err)
				}
			}
		}
	}
}
