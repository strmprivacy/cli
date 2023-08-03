package organization

import (
	"context"
	"errors"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/strmprivacy/api-definitions-go/v3/api/organizations/v1"
	"google.golang.org/grpc"
	"os"
	"strings"
	"strmprivacy/strm/pkg/common"
	"strmprivacy/strm/pkg/util"
)

var client organizations.OrganizationsServiceClient
var apiContext context.Context

func SetupClient(clientConnection *grpc.ClientConn, ctx context.Context) {
	apiContext = ctx
	client = organizations.NewOrganizationsServiceClient(clientConnection)
}

func inviteUsers(args []string, cmd *cobra.Command) {
	if apiContext == nil {
		common.CliExit(errors.New(fmt.Sprintf("No login information found. Use: %s auth login first.", common.RootCommandName)))
	}
	emails := getEmails(args, cmd)
	var invites []*organizations.UserInvite
	for _, email := range emails {
		invites = append(invites, &organizations.UserInvite{Email: email})
	}
	req := &organizations.InviteUsersRequest{
		UserInvites: invites,
	}
	response, err := client.InviteUsers(apiContext, req)
	common.CliExit(err)
	handleInviteResponse(response)
}

func handleInviteResponse(response *organizations.InviteUsersResponse) {
	fmt.Println(fmt.Sprintf("Invited %d users to your organization.", response.InviteCount))
	if len(response.Issues) > 0 {
		fmt.Println(fmt.Sprintf("There were %d invites with issues:\n", len(response.Issues)))
		inviteIssuesPrinter{}.Print(response.Issues)
	}
}

func getEmails(args []string, cmd *cobra.Command) []string {
	emailsFile := util.GetStringAndErr(cmd.Flags(), userEmailsFileFlag)
	var emails []string

	if len(args) == 0 && emailsFile == "" {
		common.CliExit(errors.New(fmt.Sprint("Either provide comma-separated emails on the command line, or a file containing emails.")))
	} else if len(args) > 0 {
		emails = strings.Split(args[0], ",")
	} else {
		emails = read(emailsFile)
	}
	return emails
}

func read(emailsFile string) []string {
	buf, err := os.ReadFile(emailsFile)
	common.CliExit(err)
	splitFn := func(c rune) bool {
		return c == '\n'
	}
	// FieldsFunc is used instead of split to filter out any (trailing) empty lines
	return strings.FieldsFunc(string(buf), splitFn)
}
