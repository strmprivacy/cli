package bootstrap

import (
	"context"
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"strings"
	"strmprivacy/strm/pkg/cmd"
	"strmprivacy/strm/pkg/common"
	"strmprivacy/strm/pkg/entity/account"
	"strmprivacy/strm/pkg/entity/batch_exporter"
	"strmprivacy/strm/pkg/entity/batch_job"
	"strmprivacy/strm/pkg/entity/data_connector"
	"strmprivacy/strm/pkg/entity/data_contract"
	"strmprivacy/strm/pkg/entity/data_subject"
	"strmprivacy/strm/pkg/entity/installation"
	"strmprivacy/strm/pkg/entity/kafka_cluster"
	"strmprivacy/strm/pkg/entity/kafka_exporter"
	"strmprivacy/strm/pkg/entity/kafka_user"
	"strmprivacy/strm/pkg/entity/key_stream"
	"strmprivacy/strm/pkg/entity/keylinks"
	"strmprivacy/strm/pkg/entity/organization"
	"strmprivacy/strm/pkg/entity/policy"
	"strmprivacy/strm/pkg/entity/project"
	"strmprivacy/strm/pkg/entity/schema_code"
	"strmprivacy/strm/pkg/entity/stream"
	"strmprivacy/strm/pkg/entity/usage"
	"strmprivacy/strm/pkg/entity/user"
	"strmprivacy/strm/pkg/logs"
	"strmprivacy/strm/pkg/monitor"
	"strmprivacy/strm/pkg/user_project"
)

const (
	cliVersionHeader = "strm-cli-version"
	zedTokenHeader   = "strm-zed-token"
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
	rootCmd.AddCommand(cmd.ReviewCmd)
	rootCmd.AddCommand(cmd.ApproveCmd)
	rootCmd.AddCommand(cmd.ActivateCmd)
	rootCmd.AddCommand(cmd.ArchiveCmd)
	rootCmd.AddCommand(cmd.InviteCmd)
	rootCmd.AddCommand(cmd.ManageCmd)
	rootCmd.AddCommand(cmd.MonitorCmd)
	rootCmd.AddCommand(cmd.LogsCmd)
	rootCmd.AddCommand(cmd.UpdateCmd)
	rootCmd.AddCommand(cmd.EvaluateCmd)
}

func SetupServiceClients(accessToken *string, zedToken *string) {
	clientConnection, ctx := SetupGrpc(common.ApiHost, accessToken, zedToken)

	stream.SetupClient(clientConnection, ctx)
	kafka_exporter.SetupClient(clientConnection, ctx)
	batch_exporter.SetupClient(clientConnection, ctx)
	batch_job.SetupClient(clientConnection, ctx)
	data_connector.SetupClient(clientConnection, ctx)
	kafka_cluster.SetupClient(clientConnection, ctx)
	kafka_user.SetupClient(clientConnection, ctx)
	key_stream.SetupClient(clientConnection, ctx)
	schema_code.SetupClient(clientConnection, ctx)
	usage.SetupClient(clientConnection, ctx)
	installation.SetupClient(clientConnection, ctx)
	account.SetupClient(clientConnection, ctx)
	project.SetupClient(clientConnection, ctx)
	user.SetupClient(clientConnection, ctx)
	organization.SetupClient(clientConnection, ctx)
	data_subject.SetupClient(clientConnection, ctx)
	keylinks.SetupClient(clientConnection, ctx)
	data_contract.SetupClient(clientConnection, ctx)
	policy.SetupClient(clientConnection, ctx)
	monitor.SetupClient(clientConnection, ctx)
	logs.SetupClient(clientConnection, ctx)
}

func InitializeConfig(cmd *cobra.Command) error {
	viperConfig := viper.New()

	// Set the base name of the config file, without the file extension.
	viperConfig.SetConfigName(common.DefaultConfigFilename)

	// Set as many paths as you like where viper should look for the
	// config file.
	viperConfig.AddConfigPath(common.ConfigPath())

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

func SetupGrpc(host string, token *string, zedToken *string) (*grpc.ClientConn, context.Context) {

	var err error
	var creds grpc.DialOption

	if strings.Contains(host, ":50051") {
		creds = grpc.WithTransportCredentials(insecure.NewCredentials())
	} else {
		creds = grpc.WithTransportCredentials(credentials.NewClientTLSFromCert(nil, ""))
	}

	clientConnection, err := grpc.Dial(host, creds, grpc.WithUnaryInterceptor(clientInterceptor))
	common.CliExit(err)

	var mdMap = map[string]string{cliVersionHeader: common.Version}

	if token != nil {
		mdMap["authorization"] = "Bearer " + *token
	}
	if zedToken != nil {
		mdMap[zedTokenHeader] = *zedToken
	}

	return clientConnection, metadata.NewOutgoingContext(context.Background(), metadata.New(mdMap))
}

func clientInterceptor(
	ctx context.Context,
	method string,
	req interface{},
	reply interface{},
	cc *grpc.ClientConn,
	invoker grpc.UnaryInvoker,
	opts ...grpc.CallOption,
) error {
	zedToken := user_project.GetZedToken()

	if zedToken != nil {
		ctx = metadata.AppendToOutgoingContext(ctx, zedTokenHeader, *zedToken)
	}

	var header metadata.MD
	opts = append(opts, grpc.Header(&header))
	err := invoker(ctx, method, req, reply, cc, opts...)

	zedTokenValue := header.Get(zedTokenHeader)

	if len(zedTokenValue) > 0 {
		user_project.SetZedToken(zedTokenValue[0])
	}

	return err
}
