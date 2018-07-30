package nsqs

import (
	"errors"

	nsq "github.com/nsqio/go-nsq"
)

var (
	// ErrTopicRequired is returned when topic is not passed as parameter.
	ErrTopicRequired = errors.New("topic is mandatory")
	// ErrHandlerRequired is returned when handler is not passed as parameter.
	ErrHandlerRequired = errors.New("handler is mandatory")
	// ErrChannelRequired is returned when channel is not passed as parameter in bus.ListenerConfig.
	ErrChannelRequired = errors.New("channel is mandatory")
)

// On listen to a message from a specific topic using nsq consumer, returns
// an error if topic and channel not passed or if an error occurred while creating
// nsq consumer.
func On(lc ListenerConfig, handler nsq.Handler) (err error) {
	if len(lc.Topic) == 0 {
		err = ErrTopicRequired
		return
	}

	if len(lc.Channel) == 0 {
		err = ErrChannelRequired
		return
	}

	if handler == nil {
		err = ErrHandlerRequired
		return
	}

	if len(lc.Lookup) == 0 {
		lc.Lookup = []string{"localhost:4161"}
	}

	if lc.HandlerConcurrency == 0 {
		lc.HandlerConcurrency = 1
	}

	config := newListenerConfig(lc)
	consumer, err := nsq.NewConsumer(lc.Topic, lc.Channel, config)
	if err != nil {
		return
	}
	addNsqStopable(consumer)

	consumer.AddConcurrentHandlers(handler, lc.HandlerConcurrency)
	err = consumer.ConnectToNSQLookupds(lc.Lookup)

	return
}
