package project

import (
	"context"
	"errors"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/strmprivacy/api-definitions-go/v2/api/entities/v1"
	"github.com/strmprivacy/api-definitions-go/v2/api/projects/v1"
	"google.golang.org/grpc"
	"io/ioutil"
	"path"
	"strmprivacy/strm/pkg/common"
	"strmprivacy/strm/pkg/util"
)

var client projects.ProjectsServiceClient
var apiContext context.Context

const (
	descriptionFlag       = "description"
	addMemberFlag         = "add-member"
	removeMemberFlag      = "remove-member"
	activeProjectFilename = "active_project"
)

func SetupClient(clientConnection *grpc.ClientConn, ctx context.Context) {
	apiContext = ctx
	client = projects.NewProjectsServiceClient(clientConnection)
}

func ListProjects() *projects.ListProjectsResponse {
	if apiContext == nil {
		common.CliExit(errors.New(fmt.Sprint("No login information found. Use: `dstrm auth login` first.")))
	}
	req := &projects.ListProjectsRequest{}
	response, err := client.ListProjects(apiContext, req)
	common.CliExit(err)
	return response
}

func GetProject(projectName string) *entities.Project {
	for _, project := range ListProjects().Projects {
		if project.Name == projectName {
			return project
		}
	}
	return nil
}

func create(projectName *string, cmd *cobra.Command) *projects.CreateProjectResponse {
	flags := cmd.Flags()
	description := util.GetStringAndErr(flags, descriptionFlag)
	req := &projects.CreateProjectRequest{
		Project: &entities.Project{
			Name:        *projectName,
			Description: description,
		},
	}
	response, err := client.CreateProject(apiContext, req)
	common.CliExit(err)
	return response
}

func manage(cmd *cobra.Command) {
	flags := cmd.Flags()
	membersToAdd, err := flags.GetStringArray(addMemberFlag)
	membersToRemove, err := flags.GetStringArray(removeMemberFlag)
	activeProject := GetActiveProject()

	addReq := &projects.AddProjectMembersRequest{
		Emails:    membersToAdd,
		ProjectId: activeProject,
	}
	removeReq := &projects.RemoveProjectMembersRequest{
		Emails:    membersToRemove,
		ProjectId: activeProject,
	}

	_, err = client.AddProjectMembers(apiContext, addReq)
	common.CliExit(err)
	_, err = client.RemoveProjectMembers(apiContext, removeReq)
	common.CliExit(err)
	return
}

func GetActiveProject() string {
	activeProjectFilePath := path.Join(common.ConfigPath, activeProjectFilename)
	bytes, err := ioutil.ReadFile(activeProjectFilePath)
	common.CliExit(err)
	return string(bytes)
}
