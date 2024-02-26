package rabbitmq

import (
	"errors"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
)

const (
	delay = 3
)

func actionConsumer() {
	for {
		select {
		case eq := <-consumerChan:
			go consumer(eq)
		}
	}
}

// Consumer 队列监听函数
func consumer(action *ExchangeQueue) {
	var (
		msgs  <-chan amqp.Delivery
		errCh <-chan error
	)

create:
	msgs, errCh = createConsumer(action)
	for {
		select {
		case msg := <-msgs:
			action.composerFunc.RabbitComposerFunc(msg.Body)
		case err := <-errCh:
			log.WithField("err:", err).WithField("queue:", action.queue).Error("Consumer tries to reconnect")
			time.Sleep(delay * time.Second)
			goto create
		}
	}

}

func createConsumer(action *ExchangeQueue) (<-chan amqp.Delivery, <-chan error) {
	deliveries := make(chan amqp.Delivery)
	errCh := make(chan error, 1)

	go func() {
		conn, err := getRabbitMqConn()
		if err != nil {
			errCh <- err
			return
		}

		ch, err := conn.Channel()
		if err != nil {
			errCh <- err
			return
		}
		defer func() {
			defer func() {
				if err := recover(); err != nil {
					log.Printf("a panic about shutting down channel rabbitMQ:%s", err)
				}
			}()
			ch.Close()
		}()

		err = ch.Qos(1, 0, false)
		if err != nil {
			errCh <- err
			return
		}

		_, err = ch.QueueDeclare(action.getQueueName(), true, false, false, false, action.getQueueArgs())
		if err != nil {
			errCh <- err
			return
		}

		msgs, err := ch.Consume(action.getQueueName(), "", false, false, false, false, nil)
		if err != nil {
			errCh <- err
			return
		}
		log.WithField("queue:", action.queue).Info("Register a consumer")
		notifyErr := make(chan *amqp.Error, 1)
		for {
			select {
			case closeErr := <-ch.NotifyClose(notifyErr):
				if closeErr != nil {
					errCh <- err
					return
				}
			case msg := <-msgs:
				if msg.Acknowledger == nil {
					errCh <- errors.New("seems to have encountered an unknown problem")
					return
				}
				deliveries <- msg
				_ = msg.Ack(false)
			}
		}

	}()

	return deliveries, errCh
}
