package installation

import (
	"github.com/spf13/cobra"
	"strmprivacy/strm/pkg/util"
)

var longDoc = util.LongDocsUsage(`
Get your Installation information.

For customers that run STRM Privacy as a service the Installation is the STRM cloud and contains
no further interesting information.

For customers that run their processing in a [Customer Data Plane](https://docs.strmprivacy.io/docs/latest/concepts/deployment-modes/)
the Installation defines:

	id:                       a uuid for the installation
	installation_type:        one of SELF_HOSTED, AWS_MARKETPLACE, AWS_MARKETPLACE_PAYG
	installation_credentials: OAuth2 client credentials,
	                          and an image pull secret for Kubernetes images
`)

var example = util.DedentTrim(`
strm get installation 39b0fcf6-b4f2-4e63-8a17-6d122500e546 -o json

{
	"id": "39b0fcf6-b4f2-4e63-8a17-6d122500e546",
	"installationType": "SELF_HOSTED",
	"installationCredentials": {
		"clientId": "aadb51f7...",
		"clientSecret": "CeXXRu..."
		"imagePullSecret": "..."
	}
}
`)

func GetCmd() *cobra.Command {
	return &cobra.Command{
		Use:               "installation (id)",
		Short:             "Get your installation by id",
		Long:              longDoc,
		Example:           example,
		DisableAutoGenTag: true,
		PreRun: func(cmd *cobra.Command, args []string) {
			printer = configurePrinter(cmd)
		},
		Run: func(cmd *cobra.Command, args []string) {
			get(&args[0], cmd)
		},
		Args:              cobra.ExactArgs(1),
		ValidArgsFunction: namesCompletion,
	}
}

func ListCmd() *cobra.Command {
	return &cobra.Command{
		Use:               "installations",
		Short:             "List your installations",
		DisableAutoGenTag: true,
		PreRun: func(cmd *cobra.Command, args []string) {
			printer = configurePrinter(cmd)
		},
		Run: func(cmd *cobra.Command, args []string) {
			list()
		},
	}
}
