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

const defaultProjectFilename = "default_project"

// ResolveProject resolves the project to use and makes its ID globally available.
// The value passed through the flag takes precedence, then the value stored in the config dir, and finally
// a fallback to default project.
func ResolveProject(f *pflag.FlagSet) {

	activeProjectFilepath := path.Join(common.ConfigPath, defaultProjectFilename)
	projectFlagValue, _ := f.GetString(ProjectFlag)

	if _, err := os.Stat(activeProjectFilepath); (os.IsNotExist(err) || GetDefaultProject() == "") && projectFlagValue == "" {
		initDefaultProject()
	}

	if projectFlagValue != "" {
		activeProject := project.GetProject(projectFlagValue)
		if activeProject == nil {
			message := fmt.Sprintf("Specified project '%v' does not exist, or you do not have access to it.", projectFlagValue)
			common.CliExit(errors.New(message))
		}
		common.ProjectId = activeProject.Id
	} else {
		activeProject := project.GetProject(GetDefaultProject())
		if activeProject == nil {
			message := fmt.Sprintf("Default project '%v' does not exist, or you do not have access to it. " +
				"A new default we be set.", GetDefaultProject())
			fmt.Println(message)
			initDefaultProject()
			os.Exit(1)
		}
		common.ProjectId = activeProject.Id
	}
}

func SetDefaultProject(projectName string) {
	if project.GetProject(projectName) != nil {
		saveDefaultProject(projectName)
		message := "Default project set to: " + projectName
		log.Infoln(message)
		fmt.Println(message)
	} else {
		message := fmt.Sprintf("No project '%v' found, or you do not have access to it.", projectName)
		log.Warnln(message)
		common.CliExit(errors.New(message))
	}
}

func GetDefaultProject() string {
	activeProjectFilepath := path.Join(common.ConfigPath, defaultProjectFilename)

	bytes, err := ioutil.ReadFile(activeProjectFilepath)
	common.CliExit(err)
	activeProject := string(bytes)
	log.Infoln("Current default project is: " + activeProject)
	return activeProject
}

func initDefaultProject() {
	projects := project.ListProjects()
	if len(projects.Projects) == 0 {
		common.CliExit(errors.New("you do not have access to any projects; create a project first, or ask to be granted access to one"))
	}
	firstProjectName := projects.Projects[0].Name
	saveDefaultProject(firstProjectName)
	message := fmt.Sprintf("Default project was not yet set, has been set to '%v'. You can set a project with 'strm context project <project-name>'\n", firstProjectName)
	log.Infoln(message)
	fmt.Println(message)
}

func saveDefaultProject(projectName string) {
	activeProjectFilepath := path.Join(common.ConfigPath, defaultProjectFilename)

	err := ioutil.WriteFile(
		activeProjectFilepath,
		[]byte(projectName),
		0644,
	)
	common.CliExit(err)
}
