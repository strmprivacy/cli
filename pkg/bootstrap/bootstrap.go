package bootstrap

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"os"
	"strings"
	"strmprivacy/strm/pkg/cmd"
	"strmprivacy/strm/pkg/common"
	"strmprivacy/strm/pkg/entity/account"
	"strmprivacy/strm/pkg/entity/batch_exporter"
	"strmprivacy/strm/pkg/entity/batch_job"
	"strmprivacy/strm/pkg/entity/data_connector"
	"strmprivacy/strm/pkg/entity/data_contract"
	"strmprivacy/strm/pkg/entity/data_subject"
	"strmprivacy/strm/pkg/entity/event_contract"
	"strmprivacy/strm/pkg/entity/installation"
	"strmprivacy/strm/pkg/entity/kafka_cluster"
	"strmprivacy/strm/pkg/entity/kafka_exporter"
	"strmprivacy/strm/pkg/entity/kafka_user"
	"strmprivacy/strm/pkg/entity/key_stream"
	"strmprivacy/strm/pkg/entity/keylinks"
	"strmprivacy/strm/pkg/entity/member"
	"strmprivacy/strm/pkg/entity/organization"
	"strmprivacy/strm/pkg/entity/policy"
	"strmprivacy/strm/pkg/entity/project"
	"strmprivacy/strm/pkg/entity/schema"
	"strmprivacy/strm/pkg/entity/schema_code"
	"strmprivacy/strm/pkg/entity/stream"
	"strmprivacy/strm/pkg/entity/usage"
	"strmprivacy/strm/pkg/util"
)

/*
*
these are the top level commands, i.e. the verbs.

Each verb sits in its own package, and will have subcommands for all the entity types
in STRM Privacy.
*/
func SetupVerbs(rootCmd *cobra.Command) {
	rootCmd.AddCommand(cmd.CreateCmd)
	rootCmd.AddCommand(cmd.GetCmd)
	rootCmd.AddCommand(cmd.DeleteCmd)
	rootCmd.AddCommand(cmd.ListCmd)
	rootCmd.AddCommand(cmd.CompletionCmd)
	rootCmd.AddCommand(cmd.SimulateCmd)
	rootCmd.AddCommand(cmd.ListenCmd)
	rootCmd.AddCommand(cmd.AuthCmd)
	rootCmd.AddCommand(cmd.VersionCmd)
	rootCmd.AddCommand(cmd.ContextCommand)
	rootCmd.AddCommand(cmd.ActivateCmd)
	rootCmd.AddCommand(cmd.ArchiveCmd)
	rootCmd.AddCommand(cmd.InviteCmd)
	rootCmd.AddCommand(cmd.ManageCmd)
}

func SetupServiceClients(accessToken *string) {
	clientConnection, ctx := common.SetupGrpc(common.ApiHost, accessToken)

	stream.SetupClient(clientConnection, ctx)
	kafka_exporter.SetupClient(clientConnection, ctx)
	batch_exporter.SetupClient(clientConnection, ctx)
	batch_job.SetupClient(clientConnection, ctx)
	data_connector.SetupClient(clientConnection, ctx)
	kafka_cluster.SetupClient(clientConnection, ctx)
	kafka_user.SetupClient(clientConnection, ctx)
	key_stream.SetupClient(clientConnection, ctx)
	schema.SetupClient(clientConnection, ctx)
	schema_code.SetupClient(clientConnection, ctx)
	event_contract.SetupClient(clientConnection, ctx)
	usage.SetupClient(clientConnection, ctx)
	installation.SetupClient(clientConnection, ctx)
	account.SetupClient(clientConnection, ctx)
	project.SetupClient(clientConnection, ctx)
	member.SetupClient(clientConnection, ctx)
	organization.SetupClient(clientConnection, ctx)
	data_subject.SetupClient(clientConnection, ctx)
	keylinks.SetupClient(clientConnection, ctx)
	data_contract.SetupClient(clientConnection, ctx)
	policy.SetupClient(clientConnection, ctx)
}

func ConfigPath() string {
	// if we set this environment variable, we work in a completely different configuration directory
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.
	// set the default configuration path
	configPathEnvVar := common.EnvPrefix + "_CONFIG_PATH"
	configPathEnv := os.Getenv(configPathEnvVar)
	defaultConfigPath := "~/.config/strmprivacy"

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
	// binds to an environment variable STRM_NUMBER. This helps
	// avoid conflicts.
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
		// keys with underscores, e.g. --favorite-color to STRM_FAVORITE_COLOR
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
