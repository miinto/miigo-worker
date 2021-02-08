package handler

import (
	"miinto.com/miigo/worker/pkg/command"
)

type HandlerEntry struct {
	Handler func(command *command.Command) (bool, error)
}