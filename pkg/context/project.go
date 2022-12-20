package context

import (
	"encoding/json"
	"errors"
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/pflag"
	"os"
	"path"
	"strmprivacy/strm/pkg/common"
	"strmprivacy/strm/pkg/entity/project"
	"strmprivacy/strm/pkg/user_project"
)

const activeProjectFilename = "active_projects.json"

// ResolveProject resolves the project to use and makes its ID globally available.
// The value passed through the flag takes precedence, then the value stored in the config dir, and finally
// a fallback to default project.
func ResolveProject(f *pflag.FlagSet) {

	activeProjectFilePath := path.Join(common.ConfigPath, activeProjectFilename)
	projectFlagValue, _ := f.GetString(common.ProjectNameFlag)

	if _, err := os.Stat(activeProjectFilePath); os.IsNotExist(err) && projectFlagValue == "" {
		initActiveProject()
		fmt.Println(fmt.Sprintf("Active project was not yet set, has been set to '%v'. You can set a project "+
			"with 'strm context project <project-name>'\n", user_project.GetActiveProject()))
	}

	if user_project.GetActiveProject() == "" && projectFlagValue == "" {
		SetActiveProject(getFirstProject())
	}

	if projectFlagValue != "" {
		resolvedProject := project.GetProject(projectFlagValue)
		if len(resolvedProject.Projects) == 0 {
			message := fmt.Sprintf("Specified project '%v' does not exist, or you do not have access to it.", projectFlagValue)
			common.CliExit(errors.New(message))
		}
		common.ProjectId = resolvedProject.Projects[0].Id
	} else {
		activeProject := user_project.GetActiveProject()
		resolvedProject := project.GetProject(activeProject)
		if len(resolvedProject.Projects) == 0 {
			initActiveProject()
			common.CliExit(errors.New(fmt.Sprintf("Active project '%v' does not exist, or you do not have access "+
				"to it. The following project has been set instead: %v", activeProject, user_project.GetActiveProject())))
		}
		common.ProjectId = resolvedProject.Projects[0].Id
	}
}

func SetActiveProject(projectName string) {
	if len(project.GetProject(projectName).Projects) != 0 {
		user_project.LoadActiveProject()
		saveActiveProject(projectName)
		message := "Active project set to: " + projectName
		log.Infoln(message)
		fmt.Println(message)
	} else {
		message := fmt.Sprintf("No project '%v' found, or you do not have access to it.", projectName)
		log.Warnln(message)
		common.CliExit(errors.New(message))
	}
}

func getFirstProject() string {
	projects := project.ListProjects()
	if len(projects) == 0 {
		common.CliExit(errors.New("you do not have access to any projects; create a project first, or ask to be granted access to one"))
	}
	return projects[0].Name
}

func initActiveProject() {
	firstProjectName := getFirstProject()
	saveActiveProject(firstProjectName)
}

func saveActiveProject(projectName string) {
	activeProjectFilepath := path.Join(common.ConfigPath, activeProjectFilename)
	user_project.Projects.SetActiveProject(projectName)
	projects, err := json.Marshal(user_project.Projects)
	if err != nil {
		common.CliExit(err)
	}

	err = os.WriteFile(
		activeProjectFilepath,
		projects,
		0644,
	)
	common.CliExit(err)
}
