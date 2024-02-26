package rabbitmq

import (
	"fmt"
	"github.com/streadway/amqp"
	"sync"
)

var consumerMap = new(sync.Map)
var consumerChan = make(chan *ExchangeQueue, 10)

func init() {
	go actionConsumer()
}

type ComposerFuncInterface interface {
	RabbitComposerFunc([]byte)
}

const (
	exchangeNameNormal = "" // 默认队列交换机
)

var (
	prefix string // 前缀
)

type ExchangeQueue struct {
	exchange     string                // 交换机名
	queue        string                // 队列名
	composerFunc ComposerFuncInterface // 消费者接口
}

// SetPrefix 设置前缀
func SetPrefix(s string) {
	prefix = s
	return
}

func (e *ExchangeQueue) getQueueArgs() amqp.Table {
	return nil
}

func (e *ExchangeQueue) getExchangeArgs() amqp.Table {
	return nil
}

func (e *ExchangeQueue) setExchangeName() {
	e.exchange = ""
	return
}

func (e *ExchangeQueue) getQueueName() string {
	return fmt.Sprintf("%s_%s", prefix, e.queue)
}

func (e *ExchangeQueue) getExchangeType() string {
	return amqp.ExchangeDirect
}
