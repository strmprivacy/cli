package kafkaconsumer

import (
	"context"
	"errors"
	"fmt"
	"github.com/Shopify/sarama"
	"github.com/spf13/cobra"
	"log"
	"os"
	"os/signal"
	"strings"
	"strmprivacy/strm/pkg/auth"
	"strmprivacy/strm/pkg/common"
	"strmprivacy/strm/pkg/entity/kafka_exporter"
	"strmprivacy/strm/pkg/util"
	"sync"
	"syscall"
)

const (
	KafkaBootstrapHostFlag = "kafka-bootstrap-hosts"
	GroupIdFlag            = "group-id"
)

func Run(cmd *cobra.Command, kafkaExporterName *string) {
	flags := cmd.Flags()
	clientId := util.GetStringAndErr(flags, common.ClientIdFlag)
	clientSecret := util.GetStringAndErr(flags, common.ClientSecretFlag)
	bootstrapBrokers := util.GetStringAndErr(flags, KafkaBootstrapHostFlag)
	kafkaExporter := kafka_exporter.Get(kafkaExporterName, cmd).KafkaExporter

	// TODO this needs to be changed. The client secret shouldn't even be
	// TODO visible on the Kafka exporters. Part of the Authorization revamp.
	clientId = kafkaExporter.Users[0].ClientId
	clientSecret = kafkaExporter.Users[0].ClientSecret
	topic := kafkaExporter.Target.Topic
	groupId := util.GetStringAndErr(flags, GroupIdFlag)
	if len(groupId) == 0 {
		common.CliExit(errors.New(fmt.Sprintf("Please set a Kafka Consumer group id with --%v", GroupIdFlag)))
	}

	//sarama.Logger = log.New(os.Stdout, "[sarama] ", log.LstdFlags)

	config := sarama.NewConfig()
	config.Consumer.Group.Rebalance.Strategy = sarama.BalanceStrategyRoundRobin
	config.Net.SASL.Enable = true
	config.Net.SASL.Mechanism = sarama.SASLTypeOAuth
	config.Net.TLS.Enable = true
	config.Consumer.Return.Errors = true
	config.Version = sarama.MaxVersion
	// TODO needs to be tested, but should work
	tokenUrl := fmt.Sprintf("%v/auth/realms/streams/protocol/openid-connect/token", common.ApiAuthHost)
	config.Net.SASL.TokenProvider = auth.NewTokenProvider(clientId, clientSecret, tokenUrl)

	consumer := Consumer{
		ready: make(chan bool),
	}

	ctx, cancel := context.WithCancel(context.Background())
	client, err := sarama.NewConsumerGroup(strings.Split(bootstrapBrokers, ","), groupId, config)
	if err != nil {
		common.CliExit(errors.New(fmt.Sprintf("Error creating consumer group client: %v", err)))
	}

	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			// `Consume` should be called inside an infinite loop, when a
			// server-side rebalance happens, the consumer session will need to be
			// recreated to get the new claims
			if err := client.Consume(ctx, []string{topic}, &consumer); err != nil {
				common.CliExit(errors.New(fmt.Sprintf("Error from consumer: %v", err)))
			}
			// check if context was cancelled, signaling that the consumer should stop
			if ctx.Err() != nil {
				return
			}
			consumer.ready = make(chan bool)
		}
	}()

	<-consumer.ready // Await till the consumer has been set up

	sigterm := make(chan os.Signal, 1)
	signal.Notify(sigterm, syscall.SIGINT, syscall.SIGTERM)
	select {
	case <-ctx.Done():
		log.Println("Shutting down Kafka Consumer")
	case <-sigterm:
		log.Println("Shutting down Kafka Consumer")
	}
	cancel()
	wg.Wait()
	if err = client.Close(); err != nil {
		common.CliExit(errors.New(fmt.Sprintf("Error closing client: %v", err)))
	}

}

// Consumer represents a Sarama consumer group consumer
type Consumer struct {
	ready chan bool
}

// Setup is run at the beginning of a new session, before ConsumeClaim
func (consumer *Consumer) Setup(sarama.ConsumerGroupSession) error {
	// Mark the consumer as ready
	close(consumer.ready)
	return nil
}

// Cleanup is run at the end of a session, once all ConsumeClaim goroutines have exited
func (consumer *Consumer) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

// ConsumeClaim must start a consumer loop of ConsumerGroupClaim's Messages().
func (consumer *Consumer) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for message := range claim.Messages() {
		fmt.Println(string(message.Value))
		session.MarkMessage(message, "")
	}

	return nil
}
