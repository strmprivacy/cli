package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"strmprivacy/strm/pkg/common"
)

var VersionCmd = &cobra.Command{
	Use:               "version",
	DisableAutoGenTag: true,
	Short:             "Print CLI version",
	Run: func(cmd *cobra.Command, args []string) {
		println(fmt.Sprintf("STRM Privacy CLI version: %s", common.Version))
		println(fmt.Sprintf("Git commit: %s", common.GitSha))
		println(fmt.Sprintf("Built on: %s", common.BuiltOn))
	},
}
