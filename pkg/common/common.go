package common

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"gopkg.in/natefinch/lumberjack.v2"
	"runtime"
	"strings"
)

// Auth properties
var BillingId = ""

// note that this can be overridden via the go build flags
var RootCommandName = "strm"

func InitLogging(configPath string) {
	log.SetLevel(log.TraceLevel)
	log.SetOutput(&lumberjack.Logger{
		Filename:   configPath + "/" + RootCommandName + ".log",
		MaxSize:    1, // MB
		MaxBackups: 0,
	})
	log.Info(fmt.Sprintf("Config path is set to: %v", configPath))
}

func CliExit(msg interface{}) {
	if msg != nil {
		_, file, line, _ := runtime.Caller(1)
		log.WithFields(log.Fields{"file": file, "line": line}).Error(msg)
		cobra.CheckErr(msg)
	}
}

func MissingIdTokenError() {
	CliExit(fmt.Sprintf("No login information found. Use: `%v auth login` first.", RootCommandName))
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

func BillingIdIsMissing() bool {
	return len(strings.TrimSpace(BillingId)) == 0
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
