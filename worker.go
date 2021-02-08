package worker

import (
	"errors"
	"miinto.com/miigo/worker/internal/channel"
	"miinto.com/miigo/worker/internal/handler"
	"miinto.com/miigo/worker/internal/process"
	"miinto.com/miigo/worker/pkg/command"
)

type Worker struct {
	handlers map[string]handler.HandlerEntry
	channels []channel.ChannelEntry
}

func NewWorkerService() *Worker {
	return &Worker{}
}

func (w *Worker) Start() error {
	if len(w.channels) == 0 {
		return errors.New("miigo-worker - no channels to listen on, start aborted")
	}

	if len(w.channels) == 1 {
		process.StartSingleMode(w.channels, w.handlers)
	} else {
		process.StartMultiMode(w.channels, w.handlers)
	}

	return nil
}

func (w *Worker) RegisterChannel(ch channel.ChannelEntry) {
	if w.channels == nil {
		w.channels = make([]channel.ChannelEntry, 0)
	}

	w.channels = append(w.channels, ch)
}

func (w *Worker) RegisterHandler(hLabel string, hFunction func(command *command.Command) (bool, error)) error {
	if w.handlers == nil {
		w.handlers = make(map[string]handler.HandlerEntry, 0)
	}

	if _, ok := w.handlers[hLabel]; ok {
		return errors.New("miigo-worker - handler label already in use")
	}

	w.handlers[hLabel] = handler.HandlerEntry{
		Handler: hFunction,
	}
	return nil
}