package cmd

import (
	"github.com/spf13/cobra"
	"streammachine.io/strm/pkg/web_socket"
)

var ListenCmd = &cobra.Command{
	Use:   "listen",
	Short: "Listen for events on a stream",
}

func init() {
	ListenCmd.AddCommand(web_socket.WebSocketCmd)
}
