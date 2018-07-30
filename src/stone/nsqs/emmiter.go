package nsqs

import (
	"encoding/json"
	"log"

	nsq "github.com/nsqio/go-nsq"
)

// Emitter exposes a interface for emitting and listening for events.
type Emitter interface {
	Emit(topic string, payload interface{}) error
	EmitAsync(topic string, payload interface{}) error
}

type eventEmitter struct {
	*nsq.Producer
	address string
}

// NewEmitter returns a new eventEmitter configured with the
// variables from the config parameter, or returning an non-nil err
// if an error occurred while creating nsq producer.
func NewEmitter(ec EmitterConfig) (emitter Emitter, err error) {
	config := newEmitterConfig(ec)

	address := ec.Address
	if len(address) == 0 {
		address = "localhost:4150"
	}

	producer, err := nsq.NewProducer(address, config)
	if err != nil {
		return
	}
	addNsqStopable(producer)

	emitter = &eventEmitter{producer, address}

	return
}

// Emit emits a message to a specific topic using nsq producer, returning
// an error if encoding payload fails or if an error occurred while publishing
// the message.
func (ee eventEmitter) Emit(topic string, payload interface{}) (err error) {
	if len(topic) == 0 {
		err = ErrTopicRequired
		return
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return
	}

	err = ee.Publish(topic, body)

	return
}

// Emit emits a message to a specific topic using nsq producer, but does not wait for
// the response from `nsqd`. Returns an error if encoding payload fails and
// logs to console if an error occurred while publishing the message.
func (ee eventEmitter) EmitAsync(topic string, payload interface{}) (err error) {
	if len(topic) == 0 {
		err = ErrTopicRequired
		return
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return
	}

	responseChan := make(chan *nsq.ProducerTransaction, 1)

	if err = ee.PublishAsync(topic, body, responseChan, ""); err != nil {
		return
	}

	go func(responseChan chan *nsq.ProducerTransaction) {
		for {
			select {
			case trans := <-responseChan:
				if trans.Error != nil {
					log.Fatalf(trans.Error.Error())
				}
			}
		}
	}(responseChan)

	return
}
