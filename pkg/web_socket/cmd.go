package web_socket

import (
	"github.com/spf13/cobra"
	"strmprivacy/strm/pkg/common"
	"strmprivacy/strm/pkg/entity/stream"
	"strmprivacy/strm/pkg/util"
)

var longDoc = util.LongDocsUsage(`Listens to the web-socket debug endpoint to receive format events on streams transformed to json.

Note that this is not meant for production purposes. There's no access to the underlying Kafka consumer
and no performance guarantees.`)

var example = util.DedentTrim(`

# Simulate some events
strm create stream test
strm simulate random-events test
Starting to simulate random strmprivacy/example/1.5.0 events to stream test. Sending one event every 1000 ms.
Sent 5 events
Sent 10 events

# And in another terminal

strm listen web-socket test
{"strmMeta": {"eventContractRef": "strmprivacy/example/1.5.0", "nonce": 1782462093, "timestamp": 1669990806395, "keyLink": "e6f...
{"strmMeta": {"eventContractRef": "strmprivacy/example/1.5.0", "nonce": 1159687711, "timestamp": 1669990807404, "keyLink": "b58...
{"strmMeta": {"eventContractRef": "strmprivacy/example/1.5.0", "nonce": -192240390, "timestamp": 1669990808413, "keyLink": "ba0...
`)

var WebSocketCmd = &cobra.Command{
	Use:               "web-socket (stream-name)",
	Short:             "Read events via the web-socket (not for production purposes)",
	DisableAutoGenTag: true,
	Long:              longDoc,
	Example:           example,
	Run: func(cmd *cobra.Command, args []string) {
		run(cmd, &args[0])
	},
	Args:              cobra.ExactArgs(1), // the stream name
	ValidArgsFunction: stream.NamesCompletion,
}

func init() {
	flags := WebSocketCmd.Flags()
	flags.String(common.ClientIdFlag, "", "client id to be used for receiving data")
	flags.String(common.ClientSecretFlag, "", "client secret to be used for receiving data")
	flags.String(WebSocketUrl, "wss://websocket.strmprivacy.io/ws", "websocket to receive events from")
}
