package interfaces

import "github.com/streadway/amqp"

type Channel interface {
	Qos(prefetchCount int, prefetchSize int, global bool) error
	Consume(queue string, consumer string, autoAck bool, exclusive bool, noLocal bool, noWait bool, args amqp.Table) (<-chan amqp.Delivery, error)
}