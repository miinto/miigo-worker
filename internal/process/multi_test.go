package process

import (
	"github.com/streadway/amqp"
	interfaces "miinto.com/miigo/worker/pkg"
	"testing"
)

func TestMultiModeProcess_CoreLoop(t *testing.T) {
	/*logger := &testSingleLogger{}
	primaryChannel := &testSingleChannel{}
	secondaryChannel := &testSingleChannel{}
	handler := &testSingleSuccessfullHandler{}

	setup := ProcessSetup{
		Channels: []miigo_channel.ChannelEntry{{
			QueueName:   "test-queue-name-0",
			ConsumerTag: "test-consumer-tag",
			AMQPChannel: primaryChannel,
		},{
			QueueName:   "test-queue-name-1",
			ConsumerTag: "test-consumer-tag",
			AMQPChannel: secondaryChannel,
		}},
		Handlers: map[string]interfaces.Handler{
			"Command\\BasicTestCommand": handler,
		},
		Logger: logger,
	}

	t.Error()

	deliveries := make([]<-chan amqp.Delivery, 0)
	for dIndex := 0; dIndex < 2; dIndex++ {
		delivery := make(chan amqp.Delivery, 2)
		for index := 0; index < 2; index++ {
			delivery <- amqp.Delivery{
				Body: []byte(`{"_type":"Command\\BasicTestCommand", "_data":{"index":` + strconv.Itoa(index) + `}}`),
			}
		}
		deliveries = append(deliveries, delivery)
		close(delivery)
	}

	p := &MultiModeProcess{}
	p.executeCoreLoop(deliveries, setup)*/
}

/*				FIXTURE TYPES			*/
type testMultiSuccessfullHandler struct {
	history []int
}
func (h *testMultiSuccessfullHandler) Handle(command interfaces.Command, logger interfaces.Logger) (bool, error) {
	h.history = append(h.history, int(command.GetPayload()["index"].(float64)))
	return true, nil
}
func (h *testMultiSuccessfullHandler) Validate(payload map[string]interface{}) (bool, error) {
	return true, nil
}

type testMultiFailedHandler struct {}
func (h *testMultiFailedHandler) Handle(command interfaces.Command, logger interfaces.Logger) (bool, error) {
	return true, nil
}
func (h *testMultiFailedHandler) Validate(payload map[string]interface{}) (bool, error) {
	return false, nil
}

type testMultiChannel struct {}
func (c *testMultiChannel) Qos(prefetchCount int, prefetchSize int, global bool) error {
	return nil
}
func (c *testMultiChannel) Consume(queue string, consumer string, autoAck bool, exclusive bool, noLocal bool, noWait bool, args amqp.Table) (<-chan amqp.Delivery, error) {
	return make(chan amqp.Delivery), nil
}

type testMultiLogger struct {}
func (l *testMultiLogger) SetMainPrefix(prefix string) {}
func (l *testMultiLogger) SetTempPrefix(prefix string) {}
func (l *testMultiLogger) LogLimited(body string) {}
func (l *testMultiLogger) LogVerbose(body string) {}