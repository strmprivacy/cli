package common

import (
	"context"
	"errors"
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"gopkg.in/natefinch/lumberjack.v2"
	"io/ioutil"
	"os"
	"path"
	"runtime"
	"strings"
)

var RootCommandName = "strm"

const activeProjectFilename = "active_project"

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
			fmt.Fprintln(os.Stderr, fmt.Sprintf(`Error code = %s
Details = %s`, (*st).Code(), (*st).Message()))
		} else {
			fmt.Fprintln(os.Stderr, err)
		}

		os.Exit(1)
	}
}

func MissingIdTokenError() {
	CliExit(errors.New(fmt.Sprintf("No login information found. Use: `%v auth login` first.", RootCommandName)))
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

func GetActiveProject() string {
	activeProjectFilePath := path.Join(ConfigPath, activeProjectFilename)

	bytes, err := ioutil.ReadFile(activeProjectFilePath)
	CliExit(err)
	activeProject := string(bytes)
	log.Infoln("Current active project is: " + activeProject)
	return activeProject
}
