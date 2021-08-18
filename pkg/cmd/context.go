package cmd

import (
	"github.com/spf13/cobra"
	"streammachine.io/strm/pkg/context"
)

var ContextCommand = &cobra.Command{
	Use:   "context",
	Short: "Print the CLI context",
}

func init() {
	ContextCommand.AddCommand(context.Configuration())
	ContextCommand.AddCommand(context.EntityInfo())
}
