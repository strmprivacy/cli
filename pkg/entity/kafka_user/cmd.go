package kafka_user

import (
	"github.com/spf13/cobra"
	"strmprivacy/strm/pkg/entity/kafka_exporter"
)

const (
	saveFlag    = "save"
	projectName = "project"
)

var longDoc = `A Kafka User is a user on a Kafka Exporter, that can be used for authentication when connecting to a Kafka Exporter. By
default, every Kafka Exporter gets one Kafka User upon creation, but these can be added/removed later.

In the current data model, the user does not have a assignable name; it is assigned upon creation. It’s still very low
level. See the end of this page for an example.

### Usage`

func DeleteCmd() *cobra.Command {
	return &cobra.Command{
		Use:               "kafka-user [name ...]",
		Short:             "Delete one or more Kafka users",
		Long:              longDoc,
		DisableAutoGenTag: true,
		PreRun: func(cmd *cobra.Command, args []string) {
			printer = configurePrinter(cmd)
		},
		Run: func(cmd *cobra.Command, args []string) {
			for i, _ := range args {
				del(&args[i])
			}
		},
		Args:              cobra.MinimumNArgs(1),
		ValidArgsFunction: namesCompletion,
	}
}

func GetCmd() *cobra.Command {
	return &cobra.Command{
		Use:               "kafka-user [name]",
		Short:             "Get Kafka user",
		Long:              longDoc,
		DisableAutoGenTag: true,
		PreRun: func(cmd *cobra.Command, args []string) {
			printer = configurePrinter(cmd)
		},
		Run: func(cmd *cobra.Command, args []string) {
			get(&args[0])
		},
		Args:              cobra.ExactArgs(1), // the stream name
		ValidArgsFunction: namesCompletion,
	}
}

func ListCmd() *cobra.Command {
	return &cobra.Command{
		Use:               "kafka-users [kafka-exporter-name]",
		Short:             "List Kafka users",
		Long:              longDoc,
		DisableAutoGenTag: true,
		PreRun: func(cmd *cobra.Command, args []string) {
			printer = configurePrinter(cmd)
		},
		Run: func(cmd *cobra.Command, args []string) {
			list(&args[0])
		},
		Args:              cobra.ExactArgs(1), // the kafka exporter name
		ValidArgsFunction: kafka_exporter.NamesCompletion,
	}
}

func CreateCmd() *cobra.Command {
	kafkaUser := &cobra.Command{
		Use:               "kafka-user [exporter-name]",
		Short:             "Create a Kafka user on a Kafka exporter",
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
	flags.String(projectName, "", `Project name to create resource in`)

	return kafkaUser
}
