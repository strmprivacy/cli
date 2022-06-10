package context

import (
	"errors"
	"fmt"
	"io/fs"
	"io/ioutil"
	"path"
	"strings"
	"strmprivacy/strm/pkg/auth"
	"strmprivacy/strm/pkg/common"
	"strmprivacy/strm/pkg/entity/account"
)

const ProjectFlag = "project"

type configuration struct {
	ConfigPath     string
	ConfigFilepath string
	Contents       string
	SavedEntities  []string
	ApiUrls        apiUrls
}

type apiUrls struct {
	ApiHost       string
	ApiAuthHost   string
	EventAuthHost string
}

func showConfiguration() {
	configFilepath := findConfigFile()
	contents, err := ioutil.ReadFile(configFilepath)
	common.CliExit(err)

	configuration := configuration{
		ConfigPath:     common.ConfigPath,
		ConfigFilepath: configFilepath,
		Contents:       string(contents),
		SavedEntities:  listSavedEntities(path.Join(common.ConfigPath, common.SavedEntitiesDirectory)),
		ApiUrls: apiUrls{
			ApiHost:       common.ApiHost,
			ApiAuthHost:   common.ApiAuthHost,
			EventAuthHost: common.EventAuthHost,
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
	contents, err := ioutil.ReadFile(filepath)
	common.CliExit(err)

	entity := savedEntity{Path: filepath, Contents: contents}
	printer.Print(entity)
}

func showAccountDetails() {
	details := account.GetAccountDetails()
	printer.Print(details)
}

func billingIdInfo() {
	b, err := auth.GetBillingId()
	if err != nil {
		if fileError, ok := err.(*fs.PathError); ok {
			common.CliExit(errors.New(fmt.Sprintf("Can't %s %s", fileError.Op, fileError.Path)))
		}
	}
	common.CliExit(err)
	printer.Print(b)
}

func listSavedEntities(p string) []string {
	files, err := ioutil.ReadDir(p)

	if err != nil {
		return []string{}
	}

	entityTypeTemplate := path.Base(p) + "/%v"
	var entities = make([]string, 0)

	for _, f := range files {
		if f.IsDir() {
			entities = append(entities, listSavedEntities(path.Join(p, f.Name()))...)
		} else if strings.HasSuffix(f.Name(), "json") && !strings.HasPrefix(f.Name(), auth.StrmCredsFilePrefix) {
			entities = append(entities, fmt.Sprintf(entityTypeTemplate, strings.Replace(f.Name(), ".json", "", -1)))
		}
	}

	return entities
}

func findConfigFile() string {
	files, err := ioutil.ReadDir(common.ConfigPath)
	common.CliExit(err)

	var configFilepath string

	for _, f := range files {
		if f.Name() == common.DefaultConfigFilename+common.DefaultConfigFileSuffix || f.Name() == common.DefaultConfigFilename+common.DefaultConfigFileSuffix {
			configFilepath = path.Join(common.ConfigPath, f.Name())
		}
	}

	return configFilepath
}

func showDefaultProject() {
	printer.Print("Default project: " + GetDefaultProject())
}
