package context

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"
	"strmprivacy/strm/pkg/common"
	"strmprivacy/strm/pkg/util"

	"github.com/jedib0t/go-pretty/v6/list"
	"github.com/spf13/cobra"
	"github.com/strmprivacy/api-definitions-go/v2/api/account/v1"
)

var printer util.Printer

func configurePrinter(command *cobra.Command) util.Printer {
	outputFormat := util.GetStringAndErr(command.Flags(), common.OutputFormatFlag)

	p := availablePrinters()[outputFormat+command.Name()]

	if p == nil {
		var allowedValues string

		switch command.Name() {
		case entityInfoCommandName:
			allowedValues = common.ContextOutputFormatFlagAllowedValuesText
		case configCommandName:
			allowedValues = common.ConfigOutputFormatFlagAllowedValuesText
		case accountCommandName:
			allowedValues = common.ConfigOutputFormatFlagAllowedValuesText
		case projectCommandName:
			allowedValues = common.ProjectOutputFormatFlagAllowedValuesText
		}

		common.CliExit(errors.New(fmt.Sprintf("Output format '%v' is not supported for '%v'. Allowed values: %v", command.CommandPath(), outputFormat, allowedValues)))
	}

	return p
}

func availablePrinters() map[string]util.Printer {
	return map[string]util.Printer{
		common.OutputFormatJsonRaw + entityInfoCommandName:  jsonRawPrinter{},
		common.OutputFormatJson + entityInfoCommandName:     jsonPrettyPrinter{},
		common.OutputFormatFilepath + entityInfoCommandName: filepathPrinter{},
		common.OutputFormatPlain + configCommandName:        configPlainPrinter{},
		common.OutputFormatJson + configCommandName:         configJsonPrinter{},
		common.OutputFormatJson + accountCommandName:        util.ProtoMessageJsonPrettyPrinter{},
		common.OutputFormatJsonRaw + accountCommandName:     util.ProtoMessageJsonRawPrinter{},
		common.OutputFormatPlain + accountCommandName:       accountPlainPrinter{},
		common.OutputFormatPlain + projectCommandName:       projectPrinter{},
	}
}

type filepathPrinter struct{}
type jsonRawPrinter struct{}
type jsonPrettyPrinter struct{}
type configPlainPrinter struct{}
type accountPlainPrinter struct{}
type configJsonPrinter struct{}
type projectPrinter struct{}

func (p filepathPrinter) Print(data interface{}) {
	entity, _ := (data).(savedEntity)

	fmt.Println(entity.Path)
}

func (p jsonRawPrinter) Print(data interface{}) {
	entity, _ := (data).(savedEntity)
	compactJson := util.CompactJson(entity.Contents)
	fmt.Println(string(compactJson.Bytes()))
}

func (p jsonPrettyPrinter) Print(data interface{}) {
	entity, _ := (data).(savedEntity)
	prettyJson := util.PrettifyJson(util.CompactJson(entity.Contents))
	fmt.Println(string(prettyJson.Bytes()))
}

func (p accountPlainPrinter) Print(data interface{}) {
	entity, _ := (data).(*account.GetAccountDetailsResponse)
	fmt.Println(fmt.Sprintf("Max Input Streams: %v", entity.MaxInputStreams))
	fmt.Println(fmt.Sprintf("Handle (used for contracts and schemas): %v", entity.Handle))
	fmt.Println(fmt.Sprintf("Subscription type: %v", entity.Subscription))
}

func (p configPlainPrinter) Print(data interface{}) {
	config, _ := (data).(configuration)

	fmt.Println(fmt.Sprintf("Configuration directory: %v", config.ConfigPath))
	fmt.Println(fmt.Sprintf("Configuration file: %v", config.ConfigFilepath))
	fmt.Println(fmt.Sprintf("Configuration file contents: \n\n    %v", strings.ReplaceAll(config.Contents, "\n", "\n    ")))

	if len(config.SavedEntities) > 0 {
		l := list.NewWriter()
		l.SetOutputMirror(os.Stdout)

		for _, entity := range config.SavedEntities {
			l.AppendItem(entity)
		}

		fmt.Println("Saved entities:")
		l.Render()
	}
}

func (p configJsonPrinter) Print(data interface{}) {
	entity, _ := (data).(configuration)
	b, _ := json.Marshal(entity)
	rawJson := util.CompactJson(b)
	fmt.Println(string(rawJson.Bytes()))
}

func (p projectPrinter) Print(data interface{}) {
	fmt.Println(data)
}
