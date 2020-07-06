package config

import (
	"encoding/json"
	"fmt"
	"github.com/spf13/viper"
	"github.com/streadway/amqp"
	"log"
	"strings"
)

type Consumer struct {
	conn *amqp.Connection
	channel *amqp.Channel
	tag string
	done chan error
}

type UpdateToken struct {
	Type string `json:"type"`
	Timestamp int `json:"timestamp"`
	OriginService string `json:"originService"`
	DestinationService string `json:"destinationService"`
	Id string `json:"id"`
}

func NewConsumer(amqpUri, exchange, exchangeType, queue, key, ctag string) error {
	c := &Consumer{
		conn:    nil,
		channel: nil,
		tag:     ctag,
		done:    make(chan error),
	}

	var err error

	log.Printf("dialing %s", amqpUri)
	c.conn, err = amqp.Dial(amqpUri)
	if err != nil {
		return fmt.Errorf("Dial: %s", err)
	}

	log.Printf("got Connection, getting Channel")
	c.channel, err = c.conn.Channel()
	if err != nil {
		return fmt.Errorf("Channel: %s", err)
	}

	log.Printf("got Channel, declaring Exchange (%s)", exchange)
	if err = c.channel.ExchangeDeclare(
		exchange,
		exchangeType,
		true,
		false,
		false,
		false,
		nil,
		); err != nil {
		return fmt.Errorf("Exchange Declare: %s", err)
	}

	log.Printf("declared Exchange, declaring Queue (%s)", queue)
	state, err := c.channel.QueueDeclare(
		queue,
		false,
		false,
		false,
		false,
		nil)
	if err != nil {
		return fmt.Errorf("Queue Declare: %s", err)
	}

	log.Printf("declared Queue (%d messages, %d consumers), bidding to Exchange (key '%s')",
		state.Messages, state.Consumers, key)

	if err = c.channel.QueueBind(
		queue,
		key,
		exchange,
		false,
		nil); err != nil {
		return fmt.Errorf("Queue Bind: %s", err)
	}

	log.Printf("Queue bound to Exchange, starting Consume (consumer tag '%s')", c.tag)
	deliveries, err := c.channel.Consume(
		queue,
		c.tag,
		false,
		false,
		false,
		false,
		nil,
		)
	if err != nil {
		return fmt.Errorf("Queue Consume: %s", err)
	}

	go handle(deliveries, c.done)

	return nil

}

func handle(deliveries <-chan amqp.Delivery, done chan error) {
	for d:= range deliveries {
		log.Printf("got %dB consumer: [%v] delivery: [%v] routingkey: [%v] %s",
			len(d.Body), d.ConsumerTag, d.DeliveryTag, d.RoutingKey, d.Body)
		handleRefreshEvent(d.Body, d.ConsumerTag)
		d.Ack(false)
	}
	log.Printf("handle: deliveries channel closed")
	done <- nil
}



func handleRefreshEvent(body []byte, consumerTag string) {
	updateToken:= &UpdateToken{}
	err := json.Unmarshal(body, updateToken)
	if err != nil {
		log.Printf("Problem parsing UpdateToken: %v", err.Error())
	} else {
		if strings.Contains(updateToken.DestinationService, consumerTag) {
			log.Println("Reloading Viper config from Spring Cloud Config server")

			LoadConfigurationFromBranch(
				viper.GetString("configServerUrl"),
				consumerTag,
				viper.GetString("profile"),
				viper.GetString("configBranch"),
				)
		}
	}
}

func StartListener(appName string, amqpServer string, exchangeName string) {
	err := NewConsumer(amqpServer, exchangeName,"topic", "config-event-queue", exchangeName, appName)
	if err != nil {
		log.Fatalf("%s", err)
	}

	log.Printf("running forever")
	select {}
}
