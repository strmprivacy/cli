package context

import (
	"errors"
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/pflag"
	"io/ioutil"
	"os"
	"path"
	"strmprivacy/strm/pkg/common"
	"strmprivacy/strm/pkg/entity/project"
)

const activeProjectFilename = "active_project"

// ResolveProject resolves the project to use and makes its ID globally available.
// The value passed through the flag takes precedence, then the value stored in the config dir, and finally
// a fallback to default project.
func ResolveProject(f *pflag.FlagSet) {

	activeProjectFilePath := path.Join(common.ConfigPath, activeProjectFilename)
	projectFlagValue, _ := f.GetString(ProjectFlag)

	if _, err := os.Stat(activeProjectFilePath); (os.IsNotExist(err) || common.GetActiveProject() == "") && projectFlagValue == "" {
		initActiveProject()
		fmt.Println(fmt.Sprintf("Active project was not yet set, has been set to '%v'. You can set a project "+
			"with 'strm context project <project-name>'\n", common.GetActiveProject()))
	}

	if projectFlagValue != "" {
		resolvedProject := project.GetProject(projectFlagValue)
		if resolvedProject == nil {
			message := fmt.Sprintf("Specified project '%v' does not exist, or you do not have access to it.", projectFlagValue)
			common.CliExit(errors.New(message))
		}
		common.ProjectId = resolvedProject.Id
	} else {
		activeProject := common.GetActiveProject()
		resolvedProject := project.GetProject(activeProject)
		if resolvedProject == nil {
			initActiveProject()
			common.CliExit(errors.New(fmt.Sprintf("Active project '%v' does not exist, or you do not have access "+
				"to it. The following project has been set instead: %v", activeProject, common.GetActiveProject())))
		}
		common.ProjectId = resolvedProject.Id
	}
}

func SetActiveProject(projectName string) {
	if project.GetProject(projectName) != nil {
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

func initActiveProject() {
	projects := project.ListProjects()
	if len(projects.Projects) == 0 {
		common.CliExit(errors.New("you do not have access to any projects; create a project first, or ask to be granted access to one"))
	}
	firstProjectName := projects.Projects[0].Name
	saveActiveProject(firstProjectName)
}

func saveActiveProject(projectName string) {
	activeProjectFilepath := path.Join(common.ConfigPath, activeProjectFilename)

	err := ioutil.WriteFile(
		activeProjectFilepath,
		[]byte(projectName),
		0644,
	)
	common.CliExit(err)
}
