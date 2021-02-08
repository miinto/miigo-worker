package main

import (
	"encoding/json"
	"github.com/streadway/amqp"
	"miinto.com/miigo/worker"
	"miinto.com/miigo/worker/internal/channel"
	"miinto.com/miigo/worker/pkg/command"
	"time"
)

type basicCommandPayload struct {
	Foo string					`json:"foo"`
}

func main() {
	w := worker.NewWorkerService()
	w.RegisterHandler("Command\\Test\\HelloWorld", handleBasicCommand)

	con, _ := amqp.Dial("amqp://miinto:miinto@localhost:5672/")
	defer con.Close()

	ch, _ := con.Channel()
	defer ch.Close()

	w.RegisterChannel(channel.ChannelEntry{
		QueueName: "go-generic-1-0",
		ConsumerTag: "miigo-worker-alpha",
		AMQPChannel: ch,
	})

	con, _ = amqp.Dial("amqp://miinto:miinto@localhost:5672/")
	defer con.Close()

	ch, _ = con.Channel()
	defer ch.Close()

	w.RegisterChannel(channel.ChannelEntry{
		QueueName: "go-generic-0-0",
		ConsumerTag: "miigo-worker-alpha",
		AMQPChannel: ch,
	})

	w.Start()
}

func handleBasicCommand(cmd *command.Command) (bool,error) {
	var payload basicCommandPayload
	_ = json.Unmarshal([]byte(cmd.Payload), &payload)

	time.Sleep(10*time.Millisecond)
	return true, nil
}