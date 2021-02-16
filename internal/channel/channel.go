package channel

import interfaces "github.com/miinto/miigo-worker/pkg"

type ChannelEntry struct {
	QueueName string
	ConsumerTag string
	AMQPChannel interfaces.Channel
}