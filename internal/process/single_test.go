package process

import (
	"github.com/streadway/amqp"
	miigo_channel "miinto.com/miigo/worker/internal/channel"
	interfaces "miinto.com/miigo/worker/pkg"
	"reflect"
	"strconv"
	"testing"
)

func TestSingleModeProcess_CoreLoop(t *testing.T) {
	logger := &testSingleLogger{}
	channel := &testSingleChannel{}
	handler := &testSingleSuccessfullHandler{}

	setup := ProcessSetup{
		Channels: []miigo_channel.ChannelEntry{{
			QueueName:   "test-queue-name",
			ConsumerTag: "test-consumer-tag",
			AMQPChannel: channel,
		}},
		Handlers: map[string]interfaces.Handler{
			"Command\\BasicTestCommand": handler,
		},
		Logger: logger,
	}

	deliveries := make(chan amqp.Delivery, 3)
	for index := 0; index < 3; index++ {
		deliveries <- amqp.Delivery{
			Body: []byte(`{"_type":"Command\\BasicTestCommand", "_data":{"index":`+strconv.Itoa(index)+`}}`),
		}
	}
	close(deliveries)

	p := &SingleModeProcess{}
	p.executeCoreLoop(deliveries, setup)

	if !reflect.DeepEqual(handler.history, []int{0,1,2}) {
		t.Error("handlers called incorrectly")
	}
}

/*				FIXTURE TYPES			*/
type testSingleSuccessfullHandler struct {
	history []int
}
func (h *testSingleSuccessfullHandler) Handle(command interfaces.Command, logger interfaces.Logger) (bool, error) {
	if h.history == nil {
		h.history = make([]int, 0)
	}
	h.history = append(h.history, int(command.GetPayload()["index"].(float64)))
	return true, nil
}
func (h *testSingleSuccessfullHandler) Validate(payload map[string]interface{}) (bool, error) {
	return true, nil
}

type testSingleFailedHandler struct {}
func (h *testSingleFailedHandler) Handle(command interfaces.Command, logger interfaces.Logger) (bool, error) {
	return true, nil
}
func (h *testSingleFailedHandler) Validate(payload map[string]interface{}) (bool, error) {
	return false, nil
}

type testSingleChannel struct {}
func (c *testSingleChannel) Qos(prefetchCount int, prefetchSize int, global bool) error {
	return nil
}
func (c *testSingleChannel) Consume(queue string, consumer string, autoAck bool, exclusive bool, noLocal bool, noWait bool, args amqp.Table) (<-chan amqp.Delivery, error) {
	return make(chan amqp.Delivery), nil
}

type testSingleLogger struct {}
func (l *testSingleLogger) SetMainPrefix(prefix string) {}
func (l *testSingleLogger) SetTempPrefix(prefix string) {}
func (l *testSingleLogger) LogLimited(body string) {}
func (l *testSingleLogger) LogVerbose(body string) {}