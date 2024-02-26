package rabbitmq

import (
	"crypto/tls"
	"fmt"
	"sync"
	"time"

	"github.com/streadway/amqp"
)

var (
	rabbitMqConn *amqp.Connection
	mx           = new(sync.Mutex)
	once         = new(sync.Once)
)

type RabbitConf struct {
	Method   string `json:"method" default:""`
	Ip       string `json:"ip" default:""`
	Port     string `json:"port" default:""`
	Username string `json:"username" default:""`
	Password string `json:"password" default:""`
}

var c *RabbitConf

// InitRabbitMq 初始化连接
func InitRabbitMq(rabbitConf *RabbitConf) (*amqp.Connection, error) {
	mx.Lock()
	defer mx.Unlock()
	if rabbitMqConn != nil && !rabbitMqConn.IsClosed() {
		return rabbitMqConn, nil
	}

	amqUrl := fmt.Sprintf("%s://%s:%s@%s:%s", rabbitConf.Method, rabbitConf.Username, rabbitConf.Password, rabbitConf.Ip, rabbitConf.Port)

	var err error
	if rabbitConf.Method == "amqps" {
		rabbitMqConn, err = amqp.DialTLS(amqUrl, &tls.Config{
			InsecureSkipVerify: true,
		})
	} else {
		rabbitMqConn, err = amqp.Dial(amqUrl)
	}

	if err != nil {
		return nil, err
	}
	c = rabbitConf
	return rabbitMqConn, nil
}

// 心跳断线重连
func heartbeat() {
	tick := time.NewTicker(time.Second * 3)
	defer tick.Stop()
	for {
		select {
		case <-tick.C:
			if rabbitMqConn == nil || rabbitMqConn.IsClosed() {
				rabbitMqConn, _ = InitRabbitMq(c)
			}
		}
	}
}

// GetRabbitMqConn 获取连接
func getRabbitMqConn() (*amqp.Connection, error) {
	var err error
	if rabbitMqConn == nil || rabbitMqConn.IsClosed() {
		rabbitMqConn, err = InitRabbitMq(c)
		if err != nil {
			return nil, err
		}
		go once.Do(heartbeat)
	}
	return rabbitMqConn, nil
}

// 声明交换机和队列并绑定
func declareQueuesExchange(ch *amqp.Channel, eq *ExchangeQueue) (*amqp.Channel, error) {
	var err error
	eq.setExchangeName()

	_, err = ch.QueueDeclare(eq.getQueueName(), true, false, false, false, eq.getQueueArgs())
	if err != nil {
		return ch, err
	}
	return ch, nil
}

func amqpPublish(ch *amqp.Channel, p *amqp.Publishing, eq *ExchangeQueue) error {
	return ch.Publish(eq.exchange, eq.getQueueName(), false, false, *p)
}
