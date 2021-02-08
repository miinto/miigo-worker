package channel

import "github.com/streadway/amqp"

type ChannelEntry struct {
	QueueName string
	ConsumerTag string
	AMQPChannel *amqp.Channel
}