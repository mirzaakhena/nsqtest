package messagebroker

import (
	nsq "github.com/nsqio/go-nsq"
)

// Context is
type Context struct {
	Message []byte
}

// Consumer is
type Consumer struct {
	consumers []*nsq.Consumer
}

// ConsumerHandler is
type ConsumerHandler struct {
	Topic               string
	FunctionHandler     func(c *Context) error
	NumberOfConcurrency int
}

// NewConsumer is
func NewConsumer(consumerHandlers []ConsumerHandler) *Consumer {

	obj := new(Consumer)

	if obj.consumers == nil {
		obj.consumers = []*nsq.Consumer{}
	}

	nsqConfig := nsq.NewConfig()
	channel := "POS_Dispatcher_Channel"
	for _, ch := range consumerHandlers {

		if ch.FunctionHandler == nil {
			panic("FunctionHandler must not nil")
		}

		consumer, err := nsq.NewConsumer(ch.Topic, channel, nsqConfig)
		if err != nil {
			panic(err.Error())
		}

		// we set n number of consumer
		consumer.AddConcurrentHandlers(&internalHandlerStruct{
			handler: ch.FunctionHandler,
		}, ch.NumberOfConcurrency)

		obj.consumers = append(obj.consumers, consumer)
	}
	return obj
}

type internalHandlerStruct struct {
	handler func(c *Context) error
}

// HandleMessage is
func (mb *internalHandlerStruct) HandleMessage(m *nsq.Message) error {
	if mb.handler == nil {
		return nil
	}
	return mb.handler(&Context{m.Body})
}

// Run is
func (mb *Consumer) Run(nsqdURL string) {
	// wg := &sync.WaitGroup{}
	for _, cs := range mb.consumers {
		// wg.Add(1)

		// if err := cs.ConnectToNSQLookupd(nsqLookupdURL); err != nil {
		if err := cs.ConnectToNSQD(nsqdURL); err != nil {
			panic(err.Error())
		}
	}
	// wg.Wait()
}
