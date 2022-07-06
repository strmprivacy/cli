package web_socket

import (
	"context"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/spf13/cobra"
	"github.com/strmprivacy/api-definitions-go/v2/api/entities/v1"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/clientcredentials"
	"net/http"
	"strmprivacy/strm/pkg/common"
	"strmprivacy/strm/pkg/entity/stream"
	"strmprivacy/strm/pkg/util"
)

const (
	WebSocketUrl = "web-socket-url"
)

var tokenSource oauth2.TokenSource

func run(cmd *cobra.Command, streamName *string) {
	s := stream.Get(streamName, false).StreamTree.Stream

	flags := cmd.Flags()
	url := util.GetStringAndErr(flags, WebSocketUrl)

	initializeTokenSource(s.Credentials[0])

	for {
		header := http.Header{"authorization": []string{"Bearer " + getToken()}}
		ws, _, err := websocket.DefaultDialer.Dial(url, header)
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
						websocket.CloseAbnormalClosure, // this happens when the server is not sending anything!
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

func initializeTokenSource(credentials *entities.Credentials) oauth2.TokenSource {
	config := clientcredentials.Config{
		ClientID:     credentials.ClientId,
		ClientSecret: credentials.ClientSecret,
		TokenURL:     fmt.Sprintf("%v/auth/realms/streams/protocol/openid-connect/token", common.ApiAuthHost),
		Scopes:       []string{"offline_access"},
	}

	return config.TokenSource(context.Background())
}

func getToken() string {
	tokens, err := tokenSource.Token()

	common.CliExit(err)

	return tokens.AccessToken
}
