package kafkaconsumer

import (
	"context"
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/spf13/cobra"
	"golang.org/x/oauth2/clientcredentials"
	"math/rand"
	"os"
	"os/signal"
	"streammachine.io/strm/pkg/common"
	"streammachine.io/strm/pkg/entity/kafka_exporter"
	"streammachine.io/strm/pkg/util"
	"syscall"
	"time"
)

const (
	KafkaBrokerFlag   = "kafka-broker"
	GroupIdFlag       = "group-id"
	SslCaLocationFlag = "ssl-ca-location"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func Run(cmd *cobra.Command, kafkaExporterName *string) {
	flags := cmd.Flags()
	clientId := util.GetStringAndErr(flags, common.ClientIdFlag)
	clientSecret := util.GetStringAndErr(flags, common.ClientSecretFlag)
	brokers := util.GetStringAndErr(flags, KafkaBrokerFlag)
	kafkaExporter := kafka_exporter.Get(kafkaExporterName).KafkaExporter

	// TODO this needs to be changed. The client secret shouldn't even be
	// TODO visible on the Kafka exporters. Part of the Authorization revamp.
	clientId = kafkaExporter.Users[0].ClientId
	clientSecret = kafkaExporter.Users[0].ClientSecret
	topic := kafkaExporter.Target.Topic
	groupId := util.GetStringAndErr(flags, GroupIdFlag)
	if len(groupId) == 0 {
		groupId = fmt.Sprintf("random-%d", rand.Int())
	}

	sslCaLocation := util.GetStringAndErr(flags, SslCaLocationFlag)

	configMap := kafka.ConfigMap{
		"bootstrap.servers":       brokers,
		"security.protocol":       "SASL_SSL",
		"sasl.mechanisms":         "OAUTHBEARER",
		"group.id":                groupId,
		"socket.keepalive.enable": "true",
		"log.connection.close":    "false",
	}
	if len(sslCaLocation) > 0 {
		_ = configMap.SetKey("ssl.ca.location", sslCaLocation)
	}

	consumer, err := kafka.NewConsumer(&configMap)
	common.CliExit(err)

	clientConfig := &clientcredentials.Config{
		ClientID:     clientId,
		ClientSecret: clientSecret,
		TokenURL:     common.EventAuthHost + "/token",
	}

	refreshToken(clientConfig, consumer)
	err = consumer.SubscribeTopics([]string{topic}, nil)
	common.CliExit(err)

	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)

	// Process messages
	// librdkafka shows `AllBrokersDown` messages for a simple tcp disconnect.
	// we're only acting on it if we have 2 in a row in the poll loop
	hadError := false
	run := true
	for run == true {
		select {
		case sig := <-sigchan:
			fmt.Printf("Caught signal %v: terminating\n", sig)
			run = false
		default:
			ev := consumer.Poll(100)
			if ev == nil {
				continue
			}
			switch e := ev.(type) {
			case *kafka.Message:
				fmt.Println(string(e.Value))
				hadError = false
			case kafka.OffsetsCommitted:
			case kafka.OAuthBearerTokenRefresh:
				refreshToken(clientConfig, consumer)
			case kafka.Error:
				if hadError {
					common.CliExit(e)
				}
				//fmt.Println("Error", e)
				hadError = true
			}
		}
	}
	_ = consumer.Close()
}

func refreshToken(config *clientcredentials.Config, consumer *kafka.Consumer) {
	token, err := config.Token(context.Background())
	common.CliExit(err)
	err = consumer.SetOAuthBearerToken(kafka.OAuthBearerToken{
		TokenValue: token.AccessToken,
		Expiration: token.Expiry,
	})
	common.CliExit(err)
	//_, _ = fmt.Fprintf(os.Stderr, "Token refreshed until %s", gostradamus.DateTimeFromTime(token.Expiry).IsoFormat())
}
