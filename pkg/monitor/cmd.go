package monitor

import (
	"github.com/spf13/cobra"
	"github.com/strmprivacy/api-definitions-go/v2/api/monitoring/v1"
	"strings"
	"strmprivacy/strm/pkg/util"
)

var longDoc = util.LongDocsUsage(``)

const followFlag = "follow"

func Command(entityType monitoring.EntityState_EntityType) *cobra.Command {
	typeLowercase := "all"
	maxArgs := 0
	short := "monitor all entity types"
	if entityType != 0 {
		typeLowercase = strings.ReplaceAll(strings.ToLower(entityType.String()), "_", "-")
		maxArgs = 1
		short = "monitor entities of type " + typeLowercase
	}

	cmd := &cobra.Command{
		Use:               typeLowercase,
		Short:             short,
		DisableAutoGenTag: true,
		Long:              longDoc,
		PreRun: func(cmd *cobra.Command, args []string) {
			printer = configurePrinter(cmd)
		},
		Run: func(cmd *cobra.Command, args []string) {
			run(cmd, entityType, args)
		},
		Args: cobra.MaximumNArgs(maxArgs), // the optional followFlag of the entity
	}

	flags := cmd.Flags()
	flags.Bool(followFlag, false, "continuously monitor these events")

	return cmd
}
