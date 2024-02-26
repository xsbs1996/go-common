package rabbitmq

import (
	"encoding/json"
	"errors"

	"github.com/streadway/amqp"
)

type QueueType = int8

const (
	Normal QueueType = iota + 1 // 普通队列
)

type Event struct {
	QueueName    string
	Body         interface{}           // 内容
	ComposerFunc ComposerFuncInterface // 消费者处理方法
}

func Publish(queueType QueueType, event *Event) error {
	var (
		err error
		eq  *ExchangeQueue
	)
	switch queueType {
	case Normal:
		err, eq = publishNormal(event)
	default:
		err = errors.New("wrong queue type")
	}
	if err == nil && eq != nil {
		go AddConsumer(eq)
	}
	return err
}

// Publish 生产普通消息
func publishNormal(event *Event) (error, *ExchangeQueue) {
	var (
		eq = &ExchangeQueue{exchangeNameNormal, event.QueueName, event.ComposerFunc}
	)

	conn, err := getRabbitMqConn()
	if err != nil {
		return err, nil
	}

	ch, err := conn.Channel()
	if err != nil {
		return err, nil
	}
	defer ch.Close()

	ch, err = declareQueuesExchange(ch, eq)
	if err != nil {
		return err, nil
	}

	//消息体
	body, err := json.Marshal(event.Body)
	if err != nil {
		return err, nil
	}

	p := &amqp.Publishing{
		DeliveryMode: amqp.Persistent,
		ContentType:  "text/plain",
		Body:         body,
	}
	err = amqpPublish(ch, p, eq)
	if err != nil {
		return err, nil
	}

	return nil, eq
}

// GetExchangeQueue 获取消费队列参数
func GetExchangeQueue(queueType QueueType, event *Event) *ExchangeQueue {
	var (
		eq *ExchangeQueue
	)
	switch queueType {
	case Normal:
		eq = &ExchangeQueue{exchangeNameNormal, event.QueueName, event.ComposerFunc}
	default:
		panic("wrong queue type")
	}
	return eq
}

// AddConsumer 增加消费者
func AddConsumer(eq *ExchangeQueue) {
	_, ok := consumerMap.Load(eq.queue)
	if ok {
		return
	}
	consumerMap.Store(eq.queue, true)
	consumerChan <- eq
	return
}
