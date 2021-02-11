package main

import (
	"fmt"
	"github.com/streadway/amqp"
	"miinto.com/miigo/worker"
	"miinto.com/miigo/worker/examples/multi_mode/internal/handler"
	"miinto.com/miigo/worker/examples/multi_mode/internal/logger"
	"miinto.com/miigo/worker/internal/channel"
)

func main() {
	w := worker.NewWorkerService()
	w.RegisterHandler("Command\\BasicCommand", &handler.BasicCommandHandler{})
	w.RegisterHandler("Command\\ComplexCommand", &handler.ComplexCommandHandler{})

	logger := &logger.Logger{}
	logger.SetMainPrefix("miigo-worker-singlemode-example")
	logger.SetTempPrefix("STARTUP MODE")
	w.RegisterLogger(logger)

	con, _ := amqp.Dial("amqp://miinto:miinto@localhost:5672/")
	defer con.Close()

	ch, _ := con.Channel()
	defer ch.Close()

	w.RegisterChannel(channel.ChannelEntry{
		QueueName: "go-generic-0-0",
		ConsumerTag: "miigo-worker-alpha",
		AMQPChannel: ch,
	})

	con, _ = amqp.Dial("amqp://miinto:miinto@localhost:5672/")
	defer con.Close()

	ch, _ = con.Channel()
	defer ch.Close()

	w.RegisterChannel(channel.ChannelEntry{
		QueueName: "go-generic-1-0",
		ConsumerTag: "miigo-worker-alpha",
		AMQPChannel: ch,
	})

	err := w.Start()
	if err != nil {
		fmt.Println(err.Error())
	}
}