package context

import (
	"fmt"
	"os"
	"path"
	"strings"
	"strmprivacy/strm/pkg/auth"
	"strmprivacy/strm/pkg/common"
	"strmprivacy/strm/pkg/entity/account"
	"strmprivacy/strm/pkg/user_project"
)

type configuration struct {
	ConfigPath     string
	ConfigFilepath string
	Contents       string
	SavedEntities  []string
	ApiUrls        apiUrls
}

type apiUrls struct {
	ApiHost     string
	ApiAuthHost string
}

func showConfiguration() {
	configFilepath := findConfigFile()
	contents, err := os.ReadFile(configFilepath)
	common.CliExit(err)

	configuration := configuration{
		ConfigPath:     common.ConfigPath,
		ConfigFilepath: configFilepath,
		Contents:       string(contents),
		SavedEntities:  listSavedEntities(path.Join(common.ConfigPath, common.SavedEntitiesDirectory)),
		ApiUrls: apiUrls{
			ApiHost:     common.ApiHost,
			ApiAuthHost: common.ApiAuthHost,
		},
	}
	printer.Print(configuration)
}

type savedEntity struct {
	Path     string
	Contents []byte
}

func entityInfo(args []string) {
	filepath := path.Join(common.ConfigPath, common.SavedEntitiesDirectory, args[0]+".json")
	contents, err := os.ReadFile(filepath)
	common.CliExit(err)

	entity := savedEntity{Path: filepath, Contents: contents}
	printer.Print(entity)
}

func showAccountDetails() {
	details := account.GetAccountDetails()
	printer.Print(details)
}

func listSavedEntities(p string) []string {
	dirEntries, err := os.ReadDir(p)

	if err != nil {
		return []string{}
	}

	entityTypeTemplate := path.Base(p) + "/%v"
	var entities = make([]string, 0)

	for _, dirEntry := range dirEntries {
		if dirEntry.IsDir() {
			entities = append(entities, listSavedEntities(path.Join(p, dirEntry.Name()))...)
		} else if strings.HasSuffix(dirEntry.Name(), "json") && !strings.HasPrefix(dirEntry.Name(), auth.StrmCredsFilePrefix) {
			entities = append(entities, fmt.Sprintf(entityTypeTemplate, strings.Replace(dirEntry.Name(), ".json", "", -1)))
		}
	}

	return entities
}

func findConfigFile() string {
	dirEntries, err := os.ReadDir(common.ConfigPath)
	common.CliExit(err)

	var configFilepath string

	for _, dirEntry := range dirEntries {
		if dirEntry.Name() == common.DefaultConfigFilename+common.DefaultConfigFileSuffix || dirEntry.Name() == common.DefaultConfigFilename+common.DefaultConfigFileSuffix {
			configFilepath = path.Join(common.ConfigPath, dirEntry.Name())
		}
	}

	return configFilepath
}

func showActiveProject() {
	printer.Print("Active project: " + user_project.GetActiveProject())
}
