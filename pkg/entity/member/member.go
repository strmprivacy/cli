package member

import (
	"context"
	"errors"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/strmprivacy/api-definitions-go/v2/api/entities/v1"
	"github.com/strmprivacy/api-definitions-go/v2/api/organizations/v1"
	"github.com/strmprivacy/api-definitions-go/v2/api/projects/v1"
	"google.golang.org/grpc"
	"strings"
	"strmprivacy/strm/pkg/common"
	"strmprivacy/strm/pkg/entity/project"
)

const (
	projectFlag      = "project"
	organizationFlag = "organization"
	usersFlag        = "users"
	rolesFlag        = "roles"
)

var (
	userRolesMap = map[string]entities.UserRole{
		"admin":         entities.UserRole_ADMIN,
		"project-admin": entities.UserRole_PROJECT_ADMIN,
		"approver":      entities.UserRole_APPROVER,
		"member":        entities.UserRole_MEMBER,
	}
)

var projectsClient projects.ProjectsServiceClient
var organizationsClient organizations.OrganizationsServiceClient
var apiContext context.Context

func SetupClient(clientConnection *grpc.ClientConn, ctx context.Context) {
	apiContext = ctx
	projectsClient = projects.NewProjectsServiceClient(clientConnection)
	organizationsClient = organizations.NewOrganizationsServiceClient(clientConnection)
}

func list(cmd *cobra.Command) []*entities.User {
	flags := cmd.Flags()
	organizationBool, err := flags.GetBool(organizationFlag)
	common.CliExit(err)
	projectBool, err := flags.GetBool(projectFlag)
	common.CliExit(err)
	if organizationBool == projectBool {
		common.CliExit(errors.New("strm list users requires organizations or projects flag, not neither nor both"))
	}
	if projectBool {
		return listProjectMembers()
	} else {
		return listOrganizationMembers()
	}
}

func listProjectMembers() []*entities.User {
	projectId := project.GetProjectIdFromName(common.GetActiveProject())
	req := &projects.ListProjectMembersRequest{
		ProjectId: projectId,
	}
	response, err := projectsClient.ListProjectMembers(apiContext, req)
	common.CliExit(err)
	return response.ProjectMembers
}

func listOrganizationMembers() []*entities.User {
	req := &organizations.ListOrganizationMembersRequest{}
	response, err := organizationsClient.ListOrganizationMembers(apiContext, req)
	common.CliExit(err)
	return response.OrganizationMembers
}

func get(email string) *entities.User {
	req := &organizations.GetUserRequest{
		Email: email,
	}
	response, err := organizationsClient.GetUser(apiContext, req)
	common.CliExit(err)
	return response.User
}

func manage(cmd *cobra.Command) *organizations.UpdateUserRolesResponse {
	flags := cmd.Flags()
	emails, err := flags.GetStringSlice(usersFlag)
	common.CliExit(err)
	roles, err := flags.GetStringSlice(rolesFlag)
	common.CliExit(err)
	var userRoles []entities.UserRole
	for i := 0; i < len(roles); i++ {
		userRole := userRolesMap[strings.ToLower(roles[i])]
		if userRole == entities.UserRole_USER_ROLE_UNSPECIFIED {
			common.CliExit(errors.New(fmt.Sprintf("\"%s\" is not a valid role.\nValid roles are: "+
				"admin, approver, project_admin, member", roles[i])))
		} else {
			userRoles = append(userRoles, userRole)
		}
	}
	if len(userRoles) == 0 {
		common.CliExit(errors.New(fmt.Sprintf("No valid user roles defined in %s.\nValid roles are: "+
			"admin, approver, project_admin, member", roles)))
	}
	req := &organizations.UpdateUserRolesRequest{
		Emails:    emails,
		UserRoles: userRoles,
	}
	response, err := organizationsClient.UpdateUserRoles(apiContext, req)
	return response
}
