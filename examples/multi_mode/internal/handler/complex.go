package handler

import (
	interfaces "miinto.com/miigo/worker/pkg"
	"time"
)

type ComplexCommandHandler struct {}

func (h *ComplexCommandHandler) Handle(command interfaces.Command, logger interfaces.Logger) (bool,error) {
	logger.Log(
		"Complex command handler ... "+command.GetPayload()["foo"].(string)+" :: "+command.GetPayload()["bar"].(map[string]interface{})["foo"].(string),
		"LIMITED",
	)

	time.Sleep(10*time.Millisecond)
	return true, nil
}

func (h *ComplexCommandHandler) Validate(payload map[string]interface{}) (bool,error) {
	return true, nil
}