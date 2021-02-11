package process

import (
	"errors"
	"fmt"
	"github.com/streadway/amqp"
	"math/rand"
	"miinto.com/miigo/worker/internal/channel"
	"miinto.com/miigo/worker/pkg"
	"miinto.com/miigo/worker/internal/command"
	"reflect"
	"time"
)

type poolEntry struct {
	name string
	pool []amqp.Delivery
}

type MasterProcessSetup struct {
	Channels []channel.ChannelEntry
	Handlers map[string]interfaces.Handler
	Logger interfaces.Logger
}

func StartSingleMode(setup MasterProcessSetup) error {
	setup.Logger.Log(fmt.Sprintf("Starting single channel mode [%v] with %d handlers ...", setup.Channels[0].QueueName, len(setup.Handlers)),"LIMITED")

	ch := setup.Channels[0]
	_ = ch.AMQPChannel.Qos(1, 0,false)
	delivery, _ := ch.AMQPChannel.Consume(ch.QueueName, ch.ConsumerTag,false, false, false, false, nil)

	loop := make(chan bool)
	go func() {
		var result bool
		var err error
		for d := range delivery {
			result, err = handleIncommingCommand(d, setup)
			if err != nil {
				setup.Logger.Log(err.Error(),"LIMITED")
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

func StartMultiMode(setup MasterProcessSetup) error {
	setup.Logger.Log(fmt.Sprintf("Starting multi channel mode [%d] with %d handlers ...", len(setup.Channels), len(setup.Handlers)), "LIMITED")

	var del <-chan amqp.Delivery
	delChans := make([]<-chan amqp.Delivery, 0)
	for _, ch := range setup.Channels {
		_ = ch.AMQPChannel.Qos(2, 0,false)
		del, _ = ch.AMQPChannel.Consume(ch.QueueName, ch.ConsumerTag,false, false, false, false, nil)

		delChans = append(delChans, del)
	}

	loop := make(chan bool)
	var cases []reflect.SelectCase
	go func() {
		messagePool := make([]poolEntry, 0)

		for _, cfgE := range setup.Channels {
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
							setup.Logger.Log(err.Error(),"LIMITED")
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
	}()
	<-loop

	return nil
}

func handleIncommingCommand(d amqp.Delivery, setup MasterProcessSetup) (bool,error) {
	cmd, err := command.NewGenericCommand(string(d.Body))
	setup.Logger.SetTempPrefix(getHID())

	if err != nil {
		return false, errors.New("ERROR: Invalid command received (NOT JSON) ["+err.Error()+"]")
	}

	if cmd.GetType() == "" {
		return false, errors.New("ERROR: Invalid command received (Not Maleficarum format)")
	}

	if hE,ok := setup.Handlers[cmd.GetType()]; ok {
		setup.Logger.Log(fmt.Sprintf("Received command [" + cmd.GetType() + "] [" + string(d.Body) + "]"),"LIMITED")

		result, err := hE.Validate(cmd.GetPayload())
		if (result == true) {
			setup.Logger.Log("Command validation successful - execution going forward.","LIMITED")
		} else {
			setup.Logger.Log("Command validation failed - execution halted and skipped.","LIMITED")
			return result, err
		}

		start := float64(time.Now().UnixNano())
		result, err = hE.Handle(cmd, setup.Logger)
		end := float64(time.Now().UnixNano())

		setup.Logger.Log(
			fmt.Sprintf("Command completed with result [%v]. Exec time [%v]", result, (end / float64(time.Second) - start / float64(time.Second))),
			"LIMITED")

		return result, err
	} else {
		return false, errors.New("ERROR: Invalid command received (Handler not registered) [" + cmd.GetType() + "]")
	}
}

func getHID() string {
	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

	s := make([]rune, 16)
	for i := range s {
		s[i] = letters[rand.Intn(len(letters))]
	}
	return "HID-"+string(s)
}