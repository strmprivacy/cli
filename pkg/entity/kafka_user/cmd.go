package kafka_user

import (
	"github.com/spf13/cobra"
	"strmprivacy/strm/pkg/entity/kafka_exporter"
	"strmprivacy/strm/pkg/util"
)

const (
	saveFlag = "save"
)

var longDoc = util.LongDocsUsage(`
A Kafka User is a user on a Kafka Exporter, that can be used for authentication when connecting to a Kafka Exporter. By
default, every Kafka Exporter gets one Kafka User upon creation, but these can be added/removed later.

In the current data model, the user does not have a assignable name; it is assigned upon creation.
`)

func DeleteCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "kafka-user (name ...)",
		Short: "Delete one or more Kafka Users",
		Long: util.DedentTrim(`
			Delete one or more Kafka Users by 'name'
			Names are randomly assigned during creation and cannot be chosen.
		`),
		DisableAutoGenTag: true,
		PreRun: func(cmd *cobra.Command, args []string) {
			printer = configurePrinter(cmd)
		},
		Run: func(cmd *cobra.Command, args []string) {
			for _, arg := range args {
				del(&arg, cmd)
			}
		},
		Args:              cobra.MinimumNArgs(1),
		ValidArgsFunction: namesCompletion,
	}
}

func GetCmd() *cobra.Command {
	return &cobra.Command{
		Use:               "kafka-user (name)",
		Short:             "Get Kafka User",
		Long:              longDoc,
		DisableAutoGenTag: true,
		PreRun: func(cmd *cobra.Command, args []string) {
			printer = configurePrinter(cmd)
		},
		Run: func(cmd *cobra.Command, args []string) {
			get(&args[0], cmd)
		},
		Args:              cobra.ExactArgs(1), // the stream name
		ValidArgsFunction: namesCompletion,
	}
}

func ListCmd() *cobra.Command {
	return &cobra.Command{
		Use:               "kafka-users (kafka-exporter-name)",
		Short:             "List Kafka Users",
		Long:              longDoc,
		DisableAutoGenTag: true,
		PreRun: func(cmd *cobra.Command, args []string) {
			printer = configurePrinter(cmd)
		},
		Run: func(cmd *cobra.Command, args []string) {
			list(&args[0], cmd)
		},
		Args:              cobra.ExactArgs(1), // the kafka exporter name
		ValidArgsFunction: kafka_exporter.NamesCompletion,
	}
}

func CreateCmd() *cobra.Command {
	kafkaUser := &cobra.Command{
		Use:               "kafka-user (exporter-name)",
		Short:             "Create a Kafka User on a Kafka Exporter",
		Long:              longDoc,
		DisableAutoGenTag: true,
		PreRun: func(cmd *cobra.Command, args []string) {
			printer = configurePrinter(cmd)
		},
		Run: func(cmd *cobra.Command, args []string) {
			streamName := &args[0]
			create(streamName, cmd)
		},
		Args:              cobra.ExactArgs(1), // the kafka-exporter name
		ValidArgsFunction: kafka_exporter.NamesCompletion,
	}
	flags := kafkaUser.Flags()
	flags.Bool(saveFlag, false, "save the result in the config directory")

	return kafkaUser
}
