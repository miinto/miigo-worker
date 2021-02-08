package process

import (
	"errors"
	"fmt"
	"github.com/streadway/amqp"
	"miinto.com/miigo/worker/internal/channel"
	"miinto.com/miigo/worker/internal/handler"
	"miinto.com/miigo/worker/pkg/command"
	"reflect"
	"time"
)

func StartSingleMode(channels []channel.ChannelEntry, handlers map[string] handler.HandlerEntry) error {
	fmt.Printf("miigo-worker - starting single channel mode [%d] with %d handlers ...\n", len(channels), len(handlers))

	ch := channels[0]
	_ = ch.AMQPChannel.Qos(1, 0,false)
	delivery, _ := ch.AMQPChannel.Consume(ch.QueueName, ch.ConsumerTag,false, false, false, false, nil)

	loop := make(chan bool)
	go func() {
		var result bool
		var err error
		for d := range delivery {
			result, err = handleIncommingCommand(d, handlers)
			if err != nil {
				fmt.Println(err.Error())
				d.Ack(false)
				continue
			}

			if result == true {
				d.Ack(false)
			} else {
				d.Nack(false, true)
			}
		}
	}()
	<-loop

	return nil
}

type poolEntry struct {
	name string
	pool []amqp.Delivery
}

func StartMultiMode(channels []channel.ChannelEntry, handlers map[string] handler.HandlerEntry) error {
	fmt.Printf("miigo-worker - starting multi channel mode [%d] with %d handlers ...\n", len(channels), len(handlers))

	var del <-chan amqp.Delivery
	delChans := make([]<-chan amqp.Delivery, 0)
	for _, ch := range channels {
		_ = ch.AMQPChannel.Qos(2, 0,false)
		del, _ = ch.AMQPChannel.Consume(ch.QueueName, ch.ConsumerTag,false, false, false, false, nil)

		delChans = append(delChans, del)
	}

	loop := make(chan bool)
	var cases []reflect.SelectCase
	go func() {
		messagePool := make([]poolEntry, 0)

		for _, cfgE := range channels {
			messagePool = append(messagePool, poolEntry{
				cfgE.QueueName,
				make([]amqp.Delivery, 0),
			})
		}

		cases = make([]reflect.SelectCase, len(delChans) + 1)
		for index, dCh := range delChans {
			cases[index] = reflect.SelectCase{
				Dir:  reflect.SelectRecv,
				Chan: reflect.ValueOf(dCh),
			}
		}
		cases[len(cases)-1] = reflect.SelectCase{Dir: reflect.SelectDefault}

		for {
			index, command, _ := reflect.Select(cases)
			if index == (len(cases)-1) {
				execute(messagePool, handlers)
			} else {
				fmt.Println(messagePool[index].name)
				messagePool[index].pool = append(messagePool[index].pool, command.Interface().(amqp.Delivery))
			}
		}
	}()
	<-loop

	return nil
}

func execute (messagePool []poolEntry, handlers map[string] handler.HandlerEntry) {
	var cmd amqp.Delivery
	for queue, _ := range messagePool {
		if len(messagePool[queue].pool) > 0 {
			cmd = messagePool[queue].pool[0]
			messagePool[queue].pool = messagePool[queue].pool[1:]
			handleIncommingCommand(cmd, handlers)
			cmd.Ack(false)
			//time.Sleep(10*time.Millisecond)
			break
		}
	}
}

func handleIncommingCommand(d amqp.Delivery, handlers map[string] handler.HandlerEntry) (bool,error) {
	var cmd *command.Command
	var err error
	cmd, err = command.CreateFromJson(string(d.Body))

	if err != nil {
		return false, errors.New("miigo-worker - invalid command received (NOT JSON) ["+err.Error()+"]")
	}

	if cmd.Type == "" {
		return false, errors.New("invalid command received (Not Maleficarum format)")
	}

	if hE,ok := handlers[cmd.Type]; ok {
		fmt.Println("miigo-worker - received command [" + cmd.Type + "] [" + string(d.Body) + "]")
		start := float64(time.Now().UnixNano())
		result, err := hE.Handler(cmd)
		end := float64(time.Now().UnixNano())
		fmt.Printf("miigo-worker - command completed with result [%v]. Exec time [%v]\n", result, (end / float64(time.Second) - start / float64(time.Second)))
		return result, err
	} else {
		return false, errors.New("miigo-worker - invalid command received (Handler not registered) [" + cmd.Type + "]")
	}
}