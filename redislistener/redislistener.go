package redislistener

import (
	"fmt"
	"sync"
	"time"

	"github.com/gomodule/redigo/redis"
)

func log(params ...interface{}) {
	fmt.Printf("[redislistener] %v \n", params)
}

// OnMessageHandler defines function to be executed when a new message arrives
type OnMessageHandler func(data []byte)

// RedisListener is the listener itself
type RedisListener struct {
	serverAddr   string
	subscribers  map[string][]OnMessageHandler
	connection   redis.Conn
	client       redis.PubSubConn
	wait         sync.WaitGroup
	readTimeout  time.Duration
	writeTimeout time.Duration
	ready        chan bool
	lastError    error
}

// NewRedisListener creates a new RedisListener
func NewRedisListener() *RedisListener {
	var result = RedisListener{}
	result.serverAddr = "localhost:6379"
	result.subscribers = make(map[string][]OnMessageHandler)
	result.readTimeout = 10 * time.Second
	result.writeTimeout = 10 * time.Second
	result.ready = make(chan bool)
	log("created", result.serverAddr)

	return &result
}

// Ready readonly channel
func (listener *RedisListener) Ready() <-chan bool {
	return listener.ready
}

// Subscribe adds OnMessageHandler for one channel
func (listener *RedisListener) Subscribe(channel string, handler OnMessageHandler) *RedisListener {
	log("subscribe", channel)
	listener.subscribers[channel] = append(listener.subscribers[channel], handler)
	return listener
}

// GetChannels returns the list of subscribed channels
func (listener *RedisListener) GetChannels() []string {
	channels := make([]string, len(listener.subscribers))

	i := 0
	for key := range listener.subscribers {
		channels[i] = key
		i++
	}
	return channels
}

// Connect establish the connection to Redis and listen to the subscribed channels
func (listener *RedisListener) Connect() *RedisListener {
	log("connect start")

	listener.connection, listener.lastError = redis.Dial("tcp", listener.serverAddr,
		redis.DialReadTimeout(listener.readTimeout),
		redis.DialWriteTimeout(listener.writeTimeout))
	defer listener.connection.Close()

	if listener.lastError != nil {
		listener.ready <- false
		log("connect error", listener.lastError)
		return listener
	}

	listener.client = redis.PubSubConn{Conn: listener.connection}
	listener.wait = sync.WaitGroup{}

	listener.lastError = listener.client.Subscribe(redis.Args{}.AddFlat(listener.GetChannels())...)
	defer listener.client.Close()

	if listener.lastError != nil {
		log("listen error", listener.lastError)
		listener.ready <- false
		return listener
	}

	go func() {
		for {
			listener.wait.Add(1)
			go listen(listener)
			go healhCheck(listener)
			listener.wait.Wait()
		}
	}()

	return listener
}

func listen(listener *RedisListener) {
	log("listen start")

	for {
		switch n := listener.client.Receive().(type) {
		case error:
			log("on error", n)

			listener.wait.Done()
			return
		case redis.Message:
			log("on message", n.Channel, n.Data)

			for _, handler := range listener.subscribers[n.Channel] {
				handler(n.Data)
			}
		case redis.Subscription:
			switch n.Count {
			// Notify application when all channels are subscribed.
			case len(listener.GetChannels()):
				log("subscriptions ready", n.Count)

				listener.ready <- true
			}
		}
	}
}

func healhCheck(listener *RedisListener) {
	log("healhCheck")

	ticker := time.NewTicker(listener.readTimeout)
	defer ticker.Stop()

	for listener.lastError == nil {
		select {
		case <-ticker.C:
			log("healhCheck ping")

			// TODO: implement reconnect
			var err = listener.client.Ping("")
			if err != nil {
				listener.wait.Done()
			}
		}
	}
}
