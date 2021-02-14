package process

import (
	"fmt"
	"github.com/streadway/amqp"
	"reflect"
)

type poolEntry struct {
	name string
	pool []amqp.Delivery
}

type MultiModeProcess struct {}

func (p *MultiModeProcess) Start(setup ProcessSetup) error {
	setup.Logger.LogLimited(fmt.Sprintf("Starting multi channel mode [%d] with %d handlers ...", len(setup.Channels), len(setup.Handlers)))
	for index, _ := range setup.Handlers {
		setup.Logger.LogLimited("Registered handler: ["+index+"]")
	}
	for _, val := range setup.Channels {
		setup.Logger.LogLimited("Listening on queue: ["+val.QueueName+"]")
	}

	var del <-chan amqp.Delivery
	deliveries := make([]<-chan amqp.Delivery, 0)
	for _, ch := range setup.Channels {
		_ = ch.AMQPChannel.Qos(2, 0,false)
		del, _ = ch.AMQPChannel.Consume(ch.QueueName, ch.ConsumerTag,false, false, false, false, nil)

		deliveries = append(deliveries, del)
	}

	loop := make(chan bool)
	go p.executeCoreLoop(deliveries, setup)
	<-loop

	return nil
}

func (p *MultiModeProcess) executeCoreLoop(deliveries []<-chan amqp.Delivery,setup ProcessSetup) {
	var cases []reflect.SelectCase
	messagePool := make([]poolEntry, 0)

	for _, cfgE := range setup.Channels {
		messagePool = append(messagePool, poolEntry{
			cfgE.QueueName,
			make([]amqp.Delivery, 0),
		})
	}

	cases = make([]reflect.SelectCase, len(deliveries) + 1)
	for index, dCh := range deliveries {
		cases[index] = reflect.SelectCase{
			Dir:  reflect.SelectRecv,
			Chan: reflect.ValueOf(dCh),
		}
	}
	cases[len(cases)-1] = reflect.SelectCase{Dir: reflect.SelectDefault}

	for {
		var cmd amqp.Delivery
		index, cmdVal, _ := reflect.Select(cases)
		// default select case - no messages available in any channels
		if index == (len(cases)-1) {
			for queue, _ := range messagePool {
				if len(messagePool[queue].pool) > 0 {
					cmd = messagePool[queue].pool[0]
					messagePool[queue].pool = messagePool[queue].pool[1:]

					result, err := handleIncommingCommand(cmd, setup)
					if err != nil {
						setup.Logger.LogLimited(err.Error())
						cmd.Ack(false)
						break
					}

					if result == true {
						cmd.Ack(false)
					} else {
						cmd.Nack(false, false)
					}
					break
				}
			}
		// a channel was selected - add the message to the message pool container
		} else {
			messagePool[index].pool = append(messagePool[index].pool, cmdVal.Interface().(amqp.Delivery))
		}
	}
}