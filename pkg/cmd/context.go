package cmd

import (
	"github.com/spf13/cobra"
	"strmprivacy/strm/pkg/context"
)

var ContextCommand = &cobra.Command{
	Use:               "context",
	DisableAutoGenTag: true,
	Short:             "Print the CLI context",
}

func init() {
	ContextCommand.AddCommand(context.Configuration())
	ContextCommand.AddCommand(context.EntityInfo())
	ContextCommand.AddCommand(context.BillingIdInfo())
	ContextCommand.AddCommand(context.Account())
}