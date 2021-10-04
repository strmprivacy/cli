package main

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"streammachine.io/strm/pkg/auth"
	"streammachine.io/strm/pkg/bootstrap"
	"streammachine.io/strm/pkg/cmd"
	"streammachine.io/strm/pkg/common"
	"streammachine.io/strm/pkg/util"
	"streammachine.io/strm/pkg/web_socket"
)

const (
	apiHostFlag = "api-host"
)

func main() {
	common.CliExit(RootCmd.Execute())
}

var RootCmd = &cobra.Command{
	Use:               common.RootCommandName,
	Short:             fmt.Sprintf("Stream Machine CLI %s", cmd.Version),
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
	persistentFlags.String(apiHostFlag, "apis.streammachine.io:443", "API host and port")
	persistentFlags.String(auth.EventsAuthUrlFlag, "https://auth.strm.services", "Event authentication host")
	persistentFlags.String(auth.ApiAuthUrlFlag, "https://accounts.streammachine.io", "User authentication host")
	persistentFlags.StringVar(&auth.TokenFile, "token-file", "",
		"Token file that contains an access token (default is $HOME/.config/stream-machine/strm-creds-<api-auth-url>.json)")
	persistentFlags.String(web_socket.WebSocketUrl, "wss://out.strm.services/ws", "Websocket to receive events from")
	persistentFlags.StringP(common.OutputFormatFlag, common.OutputFormatFlagShort, common.OutputFormatTable, fmt.Sprintf("Output format [%v]", common.OutputFormatFlagAllowedValuesText))

	err := RootCmd.RegisterFlagCompletionFunc(common.OutputFormatFlag, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return common.OutputFormatFlagAllowedValues, cobra.ShellCompDirectiveNoFileComp
	})

	common.CliExit(err)
	bootstrap.SetupVerbs(RootCmd)
}
