package process

import (
	"errors"
	"fmt"
	"github.com/streadway/amqp"
	"math/rand"
	"miinto.com/miigo/worker/internal/channel"
	"miinto.com/miigo/worker/internal/command"
	"miinto.com/miigo/worker/pkg"
	"time"
)

type ProcessSetup struct {
	Channels []channel.ChannelEntry
	Handlers map[string]interfaces.Handler
	Logger interfaces.Logger
}

func handleIncommingCommand(d amqp.Delivery, setup ProcessSetup) (bool,error) {
	cmd, err := command.NewGenericCommand(string(d.Body))
	setup.Logger.SetTempPrefix(getHID())

	if err != nil {
		return false, errors.New("ERROR: Invalid command received (NOT JSON) ["+err.Error()+"]")
	}

	if cmd.GetType() == "" {
		return false, errors.New("ERROR: Invalid command received (Not Miigo format)")
	}

	if hE,ok := setup.Handlers[cmd.GetType()]; ok {
		setup.Logger.LogLimited(fmt.Sprintf("Received command [" + cmd.GetType() + "] [" + string(d.Body) + "]"))

		result, err := hE.Validate(cmd.GetPayload())
		if (result == true) {
			setup.Logger.LogLimited("Command validation successful - execution going forward.")
		} else {
			setup.Logger.LogLimited("Command validation failed - execution halted and skipped.")
			return result, err
		}

		start := float64(time.Now().UnixNano())
		result, err = hE.Handle(cmd, setup.Logger)
		end := float64(time.Now().UnixNano())

		setup.Logger.LogLimited(fmt.Sprintf("Command completed with result [%v]. Exec time [%v]", result, (end / float64(time.Second) - start / float64(time.Second))))

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