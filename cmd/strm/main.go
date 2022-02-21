package main

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/cobra/doc"
	"github.com/spf13/pflag"
	"path"
	"path/filepath"
	"strings"
	"strmprivacy/strm/pkg/auth"
	"strmprivacy/strm/pkg/bootstrap"
	"strmprivacy/strm/pkg/cmd"
	"strmprivacy/strm/pkg/common"
	"strmprivacy/strm/pkg/kafkaconsumer"
	"strmprivacy/strm/pkg/util"
	"strmprivacy/strm/pkg/web_socket"
)

const (
	apiHostFlag      = "api-host"
	generateDocsFlag = "generate-docs"
)

func main() {
	flags := RootCmd.Flags()
	flags.Bool(generateDocsFlag, false, "generate docs")
	err := flags.MarkHidden("generate-docs")

	if err != nil {
		return
	}

	err = RootCmd.Execute()
	if err != nil {
		common.CliExit(err)
	}

	const fmTemplate = `---
title: "%s"
---
`

	linkHandler := func(name string) string {
		pathArray := strings.Split(name, "strm_")
		name = pathArray[len(pathArray)-1]
		return "/cli-reference/" + strings.Replace(name, "_", "/", -1)
	}

	filePrepender := func(filename string) string {
		fmt.Println(filepath.Base(filename))
		pathArray := strings.Split(filename, "_")
		name := pathArray[len(pathArray)-1]
		base := strings.TrimSuffix(name, path.Ext(name))
		return fmt.Sprintf(fmTemplate, strings.Replace(base, "_", " ", -1))
	}

	if util.GetBoolAndErr(flags, generateDocsFlag) {
		err := doc.GenMarkdownTreeCustom(RootCmd, "./generated_docs", filePrepender, linkHandler)
		if err != nil {
			common.CliExit(err)
		}
	}
}

var RootCmd = &cobra.Command{
	Use:               common.RootCommandName,
	Short:             fmt.Sprintf("STRM Privacy CLI %s", cmd.Version),
	PersistentPreRunE: rootCmdPreRun(),
}

func rootCmdPreRun() func(cmd *cobra.Command, args []string) error {
	return func(cmd *cobra.Command, args []string) error {
		util.CreateConfigDirAndFileIfNotExists()
		err := bootstrap.InitializeConfig(cmd)

		log.Infoln(fmt.Sprintf("Executing command: %v", cmd.CommandPath()))
		cmd.Flags().Visit(func(flag *pflag.Flag) {
			log.Infoln(fmt.Sprintf("flag %v=%v", flag.Name, flag.Value))
		})

		common.ApiHost = util.GetStringAndErr(cmd.Flags(), apiHostFlag)
		common.ApiAuthHost = util.GetStringAndErr(cmd.Flags(), auth.ApiAuthUrlFlag)
		common.EventAuthHost = util.GetStringAndErr(cmd.Flags(), auth.EventsAuthUrlFlag)

		if auth.Auth.LoadLogin() == nil {
			bootstrap.SetupServiceClients(auth.Auth.GetToken())
		}

		return err
	}
}

func init() {
	common.ConfigPath = bootstrap.ConfigPath()
	common.InitLogging()

	persistentFlags := RootCmd.PersistentFlags()
	persistentFlags.String(apiHostFlag, "api.strmprivacy.io:443", "API host and port")
	persistentFlags.String(auth.EventsAuthUrlFlag, "https://sts.strmprivacy.io", "Event authentication host")
	persistentFlags.String(auth.ApiAuthUrlFlag, "https://accounts.strmprivacy.io", "User authentication host")
	persistentFlags.StringVar(&auth.TokenFile, "token-file", "",
		"Token file that contains an access token (default is $HOME/.config/strmprivacy/credentials-<api-auth-url>.json)")
	persistentFlags.String(web_socket.WebSocketUrl, "wss://websocket.strmprivacy.io/ws", "Websocket to receive events from")
	persistentFlags.String(kafkaconsumer.KafkaBootstrapHostFlag, "export-bootstrap.kafka.strmprivacy.io:9092", "Kafka bootstrap brokers, separated by comma")
	persistentFlags.StringP(common.OutputFormatFlag, common.OutputFormatFlagShort, common.OutputFormatTable, fmt.Sprintf("Output format [%v]", common.OutputFormatFlagAllowedValuesText))

	err := RootCmd.RegisterFlagCompletionFunc(common.OutputFormatFlag, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return common.OutputFormatFlagAllowedValues, cobra.ShellCompDirectiveNoFileComp
	})

	common.CliExit(err)
	bootstrap.SetupVerbs(RootCmd)
}
