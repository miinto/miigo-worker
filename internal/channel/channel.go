package channel

import interfaces "miinto.com/miigo/worker/pkg"

type ChannelEntry struct {
	QueueName string
	ConsumerTag string
	AMQPChannel interfaces.Channel
}