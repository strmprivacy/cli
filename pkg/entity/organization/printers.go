package organization

import (
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/strmprivacy/api-definitions-go/v3/api/organizations/v1"
	"strmprivacy/strm/pkg/util"
)

type inviteIssuesPrinter struct{}

func (p inviteIssuesPrinter) Print(data interface{}) {
	issues, _ := (data).([]*organizations.InviteUsersResponse_UserInviteIssue)
	rows := make([]table.Row, 0, len(issues))

	for _, issue := range issues {
		row := table.Row{
			issue.Invite.Email,
			issue.Message,
		}
		rows = append(rows, row)
	}

	header := table.Row{
		"Email",
		"Issue",
	}

	util.RenderTable(header, rows)
}
