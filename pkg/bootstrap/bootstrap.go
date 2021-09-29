package bootstrap

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"os"
	"streammachine.io/strm/pkg/cmd"
	"streammachine.io/strm/pkg/common"
	"streammachine.io/strm/pkg/entity"
	"streammachine.io/strm/pkg/entity/batch_exporter"
	"streammachine.io/strm/pkg/entity/event_contract"
	"streammachine.io/strm/pkg/entity/kafka_cluster"
	"streammachine.io/strm/pkg/entity/kafka_exporter"
	"streammachine.io/strm/pkg/entity/kafka_user"
	"streammachine.io/strm/pkg/entity/key_stream"
	"streammachine.io/strm/pkg/entity/schema"
	"streammachine.io/strm/pkg/entity/schema_code"
	"streammachine.io/strm/pkg/entity/sink"
	"streammachine.io/strm/pkg/entity/stream"
	"streammachine.io/strm/pkg/entity/usage"
	"streammachine.io/strm/pkg/util"
	"strings"
)

/**
these are the top level commands, i.e. the verbs.

Each verb sits in its own package, and will have subcommands for all the entity types
in Stream Machine.
*/
func SetupVerbs(rootCmd *cobra.Command) {
	rootCmd.AddCommand(cmd.CreateCmd)
	rootCmd.AddCommand(cmd.GetCmd)
	rootCmd.AddCommand(cmd.DeleteCmd)
	rootCmd.AddCommand(cmd.ListCmd)
	rootCmd.AddCommand(cmd.CompletionCmd)
	rootCmd.AddCommand(cmd.SimCmd)
	rootCmd.AddCommand(cmd.EgressCmd)
	rootCmd.AddCommand(cmd.AuthCmd)
	rootCmd.AddCommand(cmd.VersionCmd)
	rootCmd.AddCommand(cmd.ContextCommand)
}

func SetupServiceClients(accessToken *string) {
	clientConnection, ctx := entity.SetupGrpc(common.ApiHost, accessToken)

	stream.SetupClient(clientConnection, ctx)
	kafka_exporter.SetupClient(clientConnection, ctx)
	batch_exporter.SetupClient(clientConnection, ctx)
	sink.SetupClient(clientConnection, ctx)
	kafka_cluster.SetupClient(clientConnection, ctx)
	kafka_user.SetupClient(clientConnection, ctx)
	key_stream.SetupClient(clientConnection, ctx)
	schema.SetupClient(clientConnection, ctx)
	schema_code.SetupClient(clientConnection, ctx)
	event_contract.SetupClient(clientConnection, ctx)
	usage.SetupClient(clientConnection, ctx)
}

func ConfigPath() string {
	// if we set this environment variable, we work in a completely different configuration directory
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.
	// set the default configuration path
	configPathEnvVar := common.EnvPrefix + "_CONFIG_PATH"
	configPathEnv := os.Getenv(configPathEnvVar)
	defaultConfigPath := "~/.config/stream-machine"

	var err error

	var configPath string
	if len(configPathEnv) != 0 {
		log.Debugln("Value for " + configPathEnvVar + " found in environment: " + configPathEnv)
		configPath, err = util.ExpandTilde(configPathEnv)
	} else {
		log.Debugln("No value for " + configPathEnvVar + " found. Falling back to default: " + defaultConfigPath)
		configPath, err = util.ExpandTilde(defaultConfigPath)
	}

	common.CliExit(err)

	return configPath
}

func InitializeConfig(cmd *cobra.Command) error {
	viperConfig := viper.New()

	// Set the base name of the config file, without the file extension.
	viperConfig.SetConfigName(common.DefaultConfigFilename)

	// Set as many paths as you like where viper should look for the
	// config file.
	viperConfig.AddConfigPath(common.ConfigPath)

	// Attempt to read the config file, gracefully ignoring errors
	// caused by a config file not being found. Return an error
	// if we cannot parse the config file.
	if err := viperConfig.ReadInConfig(); err != nil {
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
	viperConfig.SetEnvPrefix(common.EnvPrefix)

	// Bind to environment variables
	// Works great for simple config names, but needs help for names
	// like --favorite-color which we fix in the bindFlags function
	viperConfig.AutomaticEnv()

	// Bind the current command's flags to viper
	bindFlags(cmd, viperConfig)

	log.Infoln(fmt.Sprintf("Viper configuration: %v", viperConfig.AllSettings()))

	return nil
}

// Bind each cobra flag to its associated viper configuration (config file and environment variable)
func bindFlags(cmd *cobra.Command, v *viper.Viper) {
	cmd.Flags().VisitAll(func(f *pflag.Flag) {
		// Environment variables can't have dashes in them, so bind them to their equivalent
		// keys with underscores, e.g. --favorite-color to STING_FAVORITE_COLOR
		if strings.Contains(f.Name, "-") {
			envVarSuffix := strings.ToUpper(strings.ReplaceAll(f.Name, "-", "_"))
			err := v.BindEnv(f.Name, fmt.Sprintf("%s_%s", common.EnvPrefix, envVarSuffix))
			common.CliExit(err)
		}

		// Apply the viper config value to the flag when the flag is not set and viper has a value
		if !f.Changed && v.IsSet(f.Name) {
			val := v.Get(f.Name)
			err := cmd.Flags().Set(f.Name, fmt.Sprintf("%v", val))
			common.CliExit(err)
		}
	})
}
