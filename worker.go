package worker

import (
	"errors"
	"math/rand"
	"github.com/miinto/miigo-worker/internal/channel"
	"github.com/miinto/miigo-worker/internal/process"
	"github.com/miinto/miigo-worker/pkg"
	"time"
)

type Worker struct {
	handlers map[string]interfaces.Handler
	channels []channel.ChannelEntry
	logger interfaces.Logger
}

func NewWorkerService() *Worker {
	return &Worker{}
}

func (w *Worker) Start() error {
	if len(w.channels) == 0 {
		return errors.New("miigo-worker - no channels to listen on, start aborted")
	}

	if len(w.handlers) == 0 {
		return errors.New("miigo-worker - no handlers registered, start aborted")
	}

	if w.logger == nil {
		return errors.New("miigo-worker - no logger registered, start aborted")
	}

	rand.Seed(time.Now().Unix())
	if len(w.channels) == 1 {
		p := &process.SingleModeProcess{}
		p.Start(process.ProcessSetup{
			Channels: w.channels,
			Handlers: w.handlers,
			Logger: w.logger,
		})
	} else {
		p := &process.MultiModeProcess{}
		p.Start(process.ProcessSetup{
			Channels: w.channels,
			Handlers: w.handlers,
			Logger: w.logger,
		})
	}

	return nil
}

func (w *Worker) RegisterChannel(queueName string, consumerTag string, chann interfaces.Channel) {
	if w.channels == nil {
		w.channels = make([]channel.ChannelEntry, 0)
	}

	ch := channel.ChannelEntry{
		QueueName: queueName,
		ConsumerTag: consumerTag,
		AMQPChannel: chann,
	}
	w.channels = append(w.channels, ch)
}

func (w *Worker) RegisterHandler(hLabel string, hStruct interfaces.Handler) error {
	if w.handlers == nil {
		w.handlers = make(map[string]interfaces.Handler, 0)
	}

	if _, ok := w.handlers[hLabel]; ok {
		return errors.New("miigo-worker - handler label already in use")
	}

	w.handlers[hLabel] = hStruct
	return nil
}

func (w *Worker) RegisterLogger(logger interfaces.Logger) {
	w.logger = logger
}