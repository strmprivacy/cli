package installation

import (
	"context"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/strmprivacy/api-definitions-go/v2/api/installations/v1"
	"google.golang.org/grpc"
	"strings"
	"strmprivacy/strm/pkg/auth"
	"strmprivacy/strm/pkg/common"
)

var client installations.InstallationsServiceClient
var apiContext context.Context

func SetupClient(clientConnection *grpc.ClientConn, ctx context.Context) {
	apiContext = ctx
	client = installations.NewInstallationsServiceClient(clientConnection)
}

func get(_ *cobra.Command) {
	req := &installations.GetInstallationRequest{}
	response, err := client.GetInstallation(apiContext, req)
	common.CliExit(err)
	printer.Print(response)
}

func namesCompletion(cmd *cobra.Command, args []string, complete string) ([]string, cobra.ShellCompDirective) {
	if len(args) != 0 && strings.Fields(cmd.Short)[0] != "Delete" {
		// this one means you don't get multiple completion suggestions for one stream if it's not a delete call
		return nil, cobra.ShellCompDirectiveNoFileComp
	}
	if auth.Auth.BillingIdAbsent() {
		return common.MissingBillingIdCompletionError(cmd.CommandPath())
	}
	req := &installations.ListInstallationsRequest{}
	response, err := client.ListInstallations(apiContext, req)
	fmt.Println(response)
	if err != nil {
		return common.GrpcRequestCompletionError(err)
	}

	installationIds := make([]string, 0, len(response.Installations))
	for _, s := range response.Installations {
		installationIds = append(installationIds, s.Id)
	}
	return installationIds, cobra.ShellCompDirectiveNoFileComp
}
