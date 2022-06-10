package common

import (
	"context"
	"errors"
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"gopkg.in/natefinch/lumberjack.v2"
	"os"
	"runtime"
	"strings"
)

var RootCommandName = "strm"

var ApiAuthHost string
var ApiHost string
var EventAuthHost string

var ProjectId string

func SetupGrpc(host string, token *string) (*grpc.ClientConn, context.Context) {

	var err error
	var creds grpc.DialOption

	if strings.Contains(host, ":50051") {
		creds = grpc.WithTransportCredentials(insecure.NewCredentials())
	} else {
		creds = grpc.WithTransportCredentials(credentials.NewClientTLSFromCert(nil, ""))
	}

	clientConnection, err := grpc.Dial(host, creds)
	CliExit(err)

	var md metadata.MD
	if token != nil {
		md = metadata.New(map[string]string{"authorization": "Bearer " + *token, "strm-cli-version": Version})
	} else {
		md = metadata.New(map[string]string{"strm-cli-version": Version})
	}

	ctx := metadata.NewOutgoingContext(context.Background(), md)
	return clientConnection, ctx
}

func InitLogging() {
	log.SetLevel(log.TraceLevel)
	log.SetOutput(&lumberjack.Logger{
		Filename:   LogFileName(),
		MaxSize:    1, // MB
		MaxBackups: 0,
	})
	log.Info(fmt.Sprintf("Config path is set to: %v", ConfigPath))
}

func LogFileName() string {
	return ConfigPath + "/" + RootCommandName + ".log"
}

func CliExit(err error) {
	if err != nil {
		_, file, line, _ := runtime.Caller(1)
		log.WithFields(log.Fields{"file": file, "line": line}).Error(err)

		st, ok := status.FromError(err)

		if ok {
			switch (*st).Code() {
			case codes.FailedPrecondition:
				fmt.Fprintln(os.Stderr, "A precondition failed for this command:", (*st).Message())
			default:
				fmt.Fprintln(os.Stderr, st)
			}
		} else {
			fmt.Fprintln(os.Stderr, err)
		}

		os.Exit(1)
	}
}

func MissingIdTokenError() {
	CliExit(errors.New(fmt.Sprintf("No login information found. Use: `%v auth login` first.", RootCommandName)))
}

func MissingBillingIdCompletionError(commandPath string) ([]string, cobra.ShellCompDirective) {
	log.Infoln(fmt.Sprintf("Called '%v' without login info", commandPath))
	cobra.CompErrorln(fmt.Sprintf("No login information found. Use: `%v auth login` first.", RootCommandName))

	return nil, cobra.ShellCompDirectiveNoFileComp
}

func GrpcRequestCompletionError(err error) ([]string, cobra.ShellCompDirective) {
	errorMessage := fmt.Sprintf("%v", err)
	log.Errorln(errorMessage)
	cobra.CompErrorln(errorMessage)

	return nil, cobra.ShellCompDirectiveNoFileComp
}

func NoFilesEmptyCompletion(cmd *cobra.Command, args []string, complete string) ([]string, cobra.ShellCompDirective) {
	return nil, cobra.ShellCompDirectiveNoFileComp
}

func MarkRequiredFlags(cmd *cobra.Command, flagNames ...string) {
	for _, flag := range flagNames {
		err := cmd.MarkFlagRequired(flag)
		CliExit(err)
	}
}
