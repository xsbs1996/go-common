package redisdb

import (
	"fmt"
	"strconv"
	"sync"
	"time"

	"github.com/go-redis/redis/v7"
)

type RedisConf struct {
	Ip       string `json:"ip" required:"true"`
	Port     string `json:"port" required:"true"`
	DB       string `json:"db" default:"0"`
	Password string `json:"password"`
}

var (
	RedisClient *redis.Client
	mx          sync.Mutex
	once        sync.Once
)

var c *RedisConf

// InitRedisDB 初始化redis
func InitRedisDB(conf *RedisConf) {
	once.Do(func() {
		c = conf
		_, err := GetRedisClient()
		if err != nil {
			panic(err)
		}
	})
}

func initClient(c *RedisConf) {
	redisDB, _ := strconv.Atoi(c.DB)
	RedisClient = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", c.Ip, c.Port),
		Password: c.Password,
		DB:       redisDB,
	})
}

// 心跳断线重连
func heartbeat() {
	tick := time.NewTicker(3 * time.Second)
	defer tick.Stop()

	for {
		select {
		case <-tick.C:
			err := RedisClient.Ping().Err()
			if err != nil {
				//重连
				initClient(c)
			}
		}
	}
}

// GetRedisClient 获取redis链接
func GetRedisClient() (*redis.Client, error) {
	mx.Lock()
	defer mx.Unlock()
	if RedisClient == nil {
		initClient(c)
		go heartbeat()
	}
	err := RedisClient.Ping().Err()
	return RedisClient, err
}
