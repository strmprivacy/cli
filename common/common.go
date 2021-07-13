package common

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"gopkg.in/natefinch/lumberjack.v2"
	"strings"
)

// Auth properties
var BillingId string

// note that this can be overridden via the go build flags
var RootCommandName = "strm"

func InitLogging(configPath string) {
	log.SetOutput(&lumberjack.Logger{
		Filename:   configPath + "/" + RootCommandName + ".log",
		MaxSize:    1, // MB
		MaxBackups: 0,
	})
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
