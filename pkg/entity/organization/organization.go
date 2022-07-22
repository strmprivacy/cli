package organization

import (
	"context"
	"errors"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/strmprivacy/api-definitions-go/v2/api/organizations/v1"
	"google.golang.org/grpc"
	"io/ioutil"
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
		common.CliExit(errors.New(fmt.Sprint("No login information found. Use: `dstrm auth login` first.")))
	}
	emails := getEmails(args, cmd)
	var invites []*organizations.UserInvite
	for _, email := range emails {
		invites = append(invites, &organizations.UserInvite{Email: email})
	}
	req := &organizations.InviteUsersRequest{
		UserInvites: invites,
	}
	_, err := client.InviteUsers(apiContext, req)
	common.CliExit(err)

	fmt.Println(fmt.Sprintf("Invited %d users to your organization", len(invites)))
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
	buf, err := ioutil.ReadFile(emailsFile)
	common.CliExit(err)
	splitFn := func(c rune) bool {
		return c == '\n'
	}
	// FieldsFunc is used instead of split to filter out any (trailing) empty lines
	return strings.FieldsFunc(string(buf), splitFn)
}