package key_stream

import (
	"github.com/spf13/cobra"
	"io/ioutil"
)

var content, _ = ioutil.ReadFile("pkg/entity/key_stream/docstring.md")

func GetCmd() *cobra.Command {
	return &cobra.Command{
		Use:               "key-stream [name]",
		Short:             "Get key stream by name",
		Long:              string(content),
		DisableAutoGenTag: true,
		PreRun: func(cmd *cobra.Command, args []string) {
			printer = configurePrinter(cmd)
		},
		Run: func(cmd *cobra.Command, args []string) {
			get(&args[0])
		},
		Args:              cobra.ExactArgs(1), // the stream name
		ValidArgsFunction: NamesCompletion,
	}
}
func ListCmd() *cobra.Command {
	return &cobra.Command{
		Use:               "key-streams",
		Short:             "List key streams",
		Long:              string(content),
		DisableAutoGenTag: true,
		PreRun: func(cmd *cobra.Command, args []string) {
			printer = configurePrinter(cmd)
		},
		Run: func(cmd *cobra.Command, args []string) {
			list()
		},
	}
}
