package context

import (
	"encoding/json"
	"fmt"
	"github.com/jedib0t/go-pretty/v6/list"
	"github.com/spf13/cobra"
	"os"
	"strmprivacy/strm/pkg/common"
	"strmprivacy/strm/pkg/util"
	"strings"
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
		case billingIdInfoCommandName:
			allowedValues = common.ConfigOutputFormatFlagAllowedValuesText
		}

		common.CliExit(fmt.Sprintf("Output format '%v' is not supported for '%v'. Allowed values: %v", command.CommandPath(), outputFormat, allowedValues))
	}

	return p
}

func availablePrinters() map[string]util.Printer {
	return map[string]util.Printer{
		common.OutputFormatJsonRaw + entityInfoCommandName:  jsonRawPrinter{},
		common.OutputFormatJson + entityInfoCommandName:     jsonPrettyPrinter{},
		common.OutputFormatFilepath + entityInfoCommandName: filepathPrinter{},
		common.OutputFormatPlain + configCommandName:        plainPrinter{},
		common.OutputFormatJson + configCommandName:         configJsonPrinter{},
		common.OutputFormatPlain + billingIdInfoCommandName: billingIdPrinter{},
	}
}

type filepathPrinter struct{}
type jsonRawPrinter struct{}
type jsonPrettyPrinter struct{}
type plainPrinter struct{}
type configJsonPrinter struct{}
type billingIdPrinter struct{}

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

func (p plainPrinter) Print(data interface{}) {
	config, _ := (data).(configuration)

	fmt.Println(fmt.Sprintf("Configuration directory: %v", config.ConfigPath))
	fmt.Println(fmt.Sprintf("Configuration file: %v", config.ConfigFilepath))
	fmt.Println(fmt.Sprintf("Configuration file contents: \n\n    %v", strings.ReplaceAll(string(config.Contents), "\n", "\n    ")))

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

func (p billingIdPrinter) Print(data interface{}) {
	fmt.Println(data)
}
