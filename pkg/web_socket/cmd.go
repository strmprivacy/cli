package web_socket

import (
	"github.com/spf13/cobra"
	"io/ioutil"
	"strmprivacy/strm/pkg/common"
	"strmprivacy/strm/pkg/entity/stream"
)

var content, _ = ioutil.ReadFile("pkg/web_socket/docstring.md")

var WebSocketCmd = &cobra.Command{
	Use:               "web-socket (stream-name)",
	Short:             "Read events via the web-socket (not for production purposes)",
	DisableAutoGenTag: true,
	Long:              string(content),
	Run: func(cmd *cobra.Command, args []string) {
		Run(cmd, &args[0])
	},
	Args:              cobra.ExactArgs(1), // the stream name
	ValidArgsFunction: stream.NamesCompletion,
}

func init() {
	flags := WebSocketCmd.Flags()
	flags.String(common.ClientIdFlag, "", "Client id to be used for receiving data")
	flags.String(common.ClientSecretFlag, "", "Client secret to be used for receiving data")
}
