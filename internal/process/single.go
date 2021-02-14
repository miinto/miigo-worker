package process

import (
	"fmt"
	"github.com/streadway/amqp"
)

type SingleModeProcess struct {}

func (p *SingleModeProcess) Start(setup ProcessSetup) error {
	setup.Logger.LogLimited(fmt.Sprintf("Starting single channel mode [queue name: %v] with %d handlers ...", setup.Channels[0].QueueName, len(setup.Handlers)))
	for index, _ := range setup.Handlers {
		setup.Logger.LogLimited("Registered handler: ["+index+"]")
	}

	ch := setup.Channels[0]
	_ = ch.AMQPChannel.Qos(1, 0,false)
	deliveries, _ := ch.AMQPChannel.Consume(ch.QueueName, ch.ConsumerTag,false, false, false, false, nil)

	loop := make(chan bool)
	go p.executeCoreLoop(deliveries, setup)
	<-loop

	return nil
}

func (p *SingleModeProcess) executeCoreLoop(deliveries <-chan amqp.Delivery, setup ProcessSetup) {
	var result bool
	var err error
	for d := range deliveries {
		result, err = handleIncommingCommand(d, setup)
		if err != nil {
			setup.Logger.LogLimited(err.Error())
			d.Ack(false)
			continue
		}

		if result == true {
			d.Ack(false)
		} else {
			d.Nack(false, true)
		}
	}
}