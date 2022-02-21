package cmd

import (
	"github.com/spf13/cobra"
	"strmprivacy/strm/pkg/kafkaconsumer"
	"strmprivacy/strm/pkg/web_socket"
)

var ListenCmd = &cobra.Command{
	Use:               "listen",
	DisableAutoGenTag: true,
	Short:             "Listen for events on a stream",
}

func init() {
	ListenCmd.AddCommand(web_socket.WebSocketCmd)
	ListenCmd.AddCommand(kafkaconsumer.Cmd)
}
