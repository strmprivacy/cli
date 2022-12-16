package monitor

import (
	"github.com/spf13/cobra"
	"github.com/strmprivacy/api-definitions-go/v2/api/monitoring/v1"
	"strings"
	"strmprivacy/strm/pkg/util"
)

var longDoc = util.LongDocsUsage(``)

func Command(entityType monitoring.EntityState_EntityType) *cobra.Command {
	typeLowercase := strings.ToLower(entityType.String())
	cmd := &cobra.Command{
		Use:               typeLowercase,
		Short:             "monitor entities of type " + typeLowercase,
		DisableAutoGenTag: true,
		Long:              longDoc,
		Run: func(cmd *cobra.Command, args []string) {
			run(cmd, entityType, args)
		},
		Args: cobra.MaximumNArgs(1), // the optional name of the entity
	}
	return cmd
}
