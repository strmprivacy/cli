package cmd

import (
	"context"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"google.golang.org/grpc"
	"os"
	"streammachine.io/strm/auth"
	"streammachine.io/strm/egress"
	"streammachine.io/strm/entity"
	"streammachine.io/strm/entity/batch_exporter"
	"streammachine.io/strm/entity/event_contract"
	"streammachine.io/strm/entity/kafka_cluster"
	"streammachine.io/strm/entity/kafka_exporter"
	"streammachine.io/strm/entity/kafka_user"
	"streammachine.io/strm/entity/key_stream"
	schema "streammachine.io/strm/entity/schema"
	"streammachine.io/strm/entity/sink"
	"streammachine.io/strm/entity/stream"
	"streammachine.io/strm/sims"
	"streammachine.io/strm/utils"
	"strings"

	"github.com/spf13/viper"
)

// note that this can be overridden via the go build flags
var CommandName = "strm"

var cfgPath string

const (
	apiHostFlag           = "api-host"
	defaultConfigFilename = "strm"

	// The environment variable prefix of all environment variables bound to our command line flags.
	// For example, --api-host is bound to STRM_API_HOST
	envPrefix = "STRM"
)

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   CommandName,
	Short: fmt.Sprintf("Stream Machine CLI %s", Version),
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		// You can bind cobra and viper in a few locations,
		// but PersistencePreRunE on the root command works well
		r := initializeConfig(cmd)
		auth.ConfigPath = cfgPath
		utils.ConfigPath = cfgPath

		auth.RootCommandName = CommandName

		if cmd.Parent() != AuthCmd && cmd != CompletionCmd && cmd != VersionCmd && cmd.Name() != "help" {
			apiAuthUrl := utils.GetStringAndErr(cmd.Flags(), auth.ApiAuthUrlFlag)

			authClient := &auth.Auth{Uri: apiAuthUrl}
			authClient.LoadLogin()
			token, billingId := authClient.GetToken(true)
			authClient.StoreLogin()

			apiHost := utils.GetStringAndErr(cmd.Flags(), apiHostFlag)

			clientConnection, ctx := entity.SetupGrpc(apiHost, token)
			sims.SetBillingId(billingId)
			egress.BillingId = billingId
			setupServiceClients(clientConnection, ctx)
			setBillingIdInHandlers(billingId)
		}

		return r
	},
}

func setupServiceClients(clientConnection *grpc.ClientConn, ctx context.Context) {
	stream.SetupClient(clientConnection, ctx)
	kafka_exporter.SetupClient(clientConnection, ctx)
	batch_exporter.SetupClient(clientConnection, ctx)
	sink.SetupClient(clientConnection, ctx)
	kafka_cluster.SetupClient(clientConnection, ctx)
	kafka_user.SetupClient(clientConnection, ctx)
	key_stream.SetupClient(clientConnection, ctx)
	schema.SetupClient(clientConnection, ctx)
	event_contract.SetupClient(clientConnection, ctx)
}

func setBillingIdInHandlers(billingId string) {
	batch_exporter.BillingId = billingId
	kafka_cluster.BillingId = billingId
	kafka_exporter.BillingId = billingId
	kafka_user.BillingId = billingId
	key_stream.BillingId = billingId
	sink.BillingId = billingId
	stream.BillingId = billingId
	schema.BillingId = billingId
	event_contract.BillingId = billingId
}

/**
these are the top level commands, i.e. the verbs.

Each verb sits in its own package, and will have subcommands for all the entity types
in Stream Machine.
*/
func setupVerbs() {
	RootCmd.AddCommand(CreateCmd)
	RootCmd.AddCommand(GetCmd)
	RootCmd.AddCommand(DeleteCmd)
	RootCmd.AddCommand(ListCmd)
	RootCmd.AddCommand(CompletionCmd)
	RootCmd.AddCommand(SimCmd)
	RootCmd.AddCommand(EgressCmd)
	RootCmd.AddCommand(AuthCmd)
	RootCmd.AddCommand(VersionCmd)
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	cobra.CheckErr(RootCmd.Execute())
}

func init() {

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.
	// set the default configuration path
	var err error
	cfgPath, err = utils.ExpandTilde("~/.config/stream-machine")
	cobra.CheckErr(err)

	RootCmd.PersistentFlags().String(apiHostFlag, "apis.streammachine.io:443", "API host and port")
	RootCmd.PersistentFlags().String(auth.EventAuthHostFlag, "https://auth.strm.services", "Security Token Service for events")
	RootCmd.PersistentFlags().String(auth.ApiAuthUrlFlag, "https://api.streammachine.io/v1", "Auth URL for user logins")
	RootCmd.PersistentFlags().StringVar(&auth.TokenFile, "token-file", "",
		"Token file that contains an access token (default is $HOME/.config/stream-machine/strm-creds-<api-auth-host>.json)")
	RootCmd.PersistentFlags().String(egress.UrlFlag, "wss://out.strm.services/ws", "Websocket to receive events from")
	setupVerbs()
}

func initializeConfig(cmd *cobra.Command) error {
	v := viper.New()

	// Set the base name of the config file, without the file extension.
	v.SetConfigName(defaultConfigFilename)

	// Set as many paths as you like where viper should look for the
	// config file.

	// if we set this environment variable, we work in a completely different configuration directory
	e := os.Getenv(envPrefix + "_CONFIG_PATH")
	if len(e) != 0 {
		p, err := utils.ExpandTilde(e)
		cobra.CheckErr(err)
		v.AddConfigPath(p)
		cfgPath = p
	} else {
		v.AddConfigPath(cfgPath)
	}

	// Attempt to read the config file, gracefully ignoring errors
	// caused by a config file not being found. Return an error
	// if we cannot parse the config file.
	if err := v.ReadInConfig(); err != nil {
		// It's okay if there isn't a config file
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return err
		}
	}

	// When we bind flags to environment variables expect that the
	// environment variables are prefixed, e.g. a flag like --number
	// binds to an environment variable STING_NUMBER. This helps
	// avoid conflicts.
	// you could set STRM_BILLINGID for instance
	v.SetEnvPrefix(envPrefix)

	// Bind to environment variables
	// Works great for simple config names, but needs help for names
	// like --favorite-color which we fix in the bindFlags function
	v.AutomaticEnv()

	// Bind the current command's flags to viper
	bindFlags(cmd, v)

	return nil
}

// Bind each cobra flag to its associated viper configuration (config file and environment variable)
func bindFlags(cmd *cobra.Command, v *viper.Viper) {
	cmd.Flags().VisitAll(func(f *pflag.Flag) {
		// Environment variables can't have dashes in them, so bind them to their equivalent
		// keys with underscores, e.g. --favorite-color to STING_FAVORITE_COLOR
		if strings.Contains(f.Name, "-") {
			envVarSuffix := strings.ToUpper(strings.ReplaceAll(f.Name, "-", "_"))
			err := v.BindEnv(f.Name, fmt.Sprintf("%s_%s", envPrefix, envVarSuffix))
			cobra.CheckErr(err)
		}

		// Apply the viper config value to the flag when the flag is not set and viper has a value
		if !f.Changed && v.IsSet(f.Name) {
			val := v.Get(f.Name)
			err := cmd.Flags().Set(f.Name, fmt.Sprintf("%v", val))
			cobra.CheckErr(err)
		}
	})
}
