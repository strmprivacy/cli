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
	"streammachine.io/strm/pkg/constants"
	"streammachine.io/strm/pkg/egress"
	"streammachine.io/strm/pkg/util"
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

		common.AuthHost = util.GetStringAndErr(cmd.Flags(), auth.ApiAuthUrlFlag)
		common.ApiHost = util.GetStringAndErr(cmd.Flags(), apiHostFlag)

		auth.SetupAuth(common.AuthHost)
		if auth.Auth.LoadLogin() == nil {
			bootstrap.SetupServiceClients(auth.Auth.GetToken())
		}

		return err
	}
}

func init() {
	constants.ConfigPath = bootstrap.ConfigPath()
	common.InitLogging()

	persistentFlags := RootCmd.PersistentFlags()
	persistentFlags.String(apiHostFlag, "apis.streammachine.io:443", "API host and port")
	persistentFlags.String(auth.EventAuthHostFlag, "https://auth.strm.services", "Security Token Service for events")
	persistentFlags.String(auth.ApiAuthUrlFlag, "https://api.streammachine.io/v1", "Auth URL for user logins")
	persistentFlags.StringVar(&auth.TokenFile, "token-file", "",
		"Token file that contains an access token (default is $HOME/.config/stream-machine/strm-creds-<api-auth-host>.json)")
	persistentFlags.String(egress.UrlFlag, "wss://out.strm.services/ws", "Websocket to receive events from")
	persistentFlags.StringP(constants.OutputFormatFlag, constants.OutputFormatFlagShort, constants.OutputFormatTable, fmt.Sprintf("Output format [%v]", constants.OutputFormatFlagAllowedValuesText))

	err := RootCmd.RegisterFlagCompletionFunc(constants.OutputFormatFlag, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return constants.OutputFormatFlagAllowedValues, cobra.ShellCompDirectiveNoFileComp
	})

	common.CliExit(err)
	bootstrap.SetupVerbs(RootCmd)
}
