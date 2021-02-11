package handler

import (
	"miinto.com/miigo/worker/pkg"
	"time"
)

type BasicCommandHandler struct {}

func (h *BasicCommandHandler) Handle(command interfaces.Command, logger interfaces.Logger) (bool,error) {
	logger.Log(
		"Basic command handler ... "+command.GetPayload()["foo"].(string),
		"LIMITED",
	)

	time.Sleep(10*time.Millisecond)
	return true, nil
}

func (h *BasicCommandHandler) Validate(payload map[string]interface{}) (bool,error) {
	return true, nil
}