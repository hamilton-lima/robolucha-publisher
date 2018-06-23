package redislistener

import (
	"fmt"
	"sync"
	"time"

	"github.com/gomodule/redigo/redis"
)

func (listener *RedisListener) log(params ...interface{}) {
	if listener.verbose {
		fmt.Print("[redislistener] ")
		for _, line := range params {
			fmt.Printf("%v ", line)
		}
		fmt.Print("\n")
	}
}

// OnMessageHandler defines function to be executed when a new message arrives
type OnMessageHandler func(data []byte)

// RedisListener is the listener itself
type RedisListener struct {
	serverAddr        string
	subscribers       map[string][]OnMessageHandler
	connection        redis.Conn
	client            redis.PubSubConn
	wait              sync.WaitGroup
	readTimeout       time.Duration
	writeTimeout      time.Duration
	healtCheckTimeout time.Duration
	ready             chan bool
	lastError         error
	verbose           bool
}

// NewRedisListener creates a new RedisListener
func NewRedisListener() *RedisListener {
	var result = RedisListener{}
	result.serverAddr = "localhost:6379"
	result.subscribers = make(map[string][]OnMessageHandler)
	result.healtCheckTimeout = time.Minute
	result.readTimeout = result.healtCheckTimeout + (10 * time.Second)
	result.writeTimeout = 10 * time.Second
	result.verbose = true

	result.ready = make(chan bool)
	result.log("created", result.serverAddr)

	return &result
}

// Ready readonly channel
func (listener *RedisListener) Ready() <-chan bool {
	return listener.ready
}

// Subscribe adds OnMessageHandler for one channel
func (listener *RedisListener) Subscribe(channel string, handler OnMessageHandler) *RedisListener {
	listener.log("subscribe", channel)
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
	listener.log("connect start")

	listener.connection, listener.lastError = redis.Dial("tcp", listener.serverAddr,
		redis.DialReadTimeout(listener.readTimeout),
		redis.DialWriteTimeout(listener.writeTimeout))

	if listener.lastError != nil {
		listener.ready <- false
		listener.log("connect error", listener.lastError)
		return listener
	}

	listener.client = redis.PubSubConn{Conn: listener.connection}
	listener.lastError = listener.client.Subscribe(redis.Args{}.AddFlat(listener.GetChannels())...)

	if listener.lastError != nil {
		listener.log("listen error", listener.lastError)
		listener.ready <- false
		return listener
	}

	go func() {
		// TODO: will allow retries, remove for now
		for {
			defer listener.client.Unsubscribe()
			defer listener.connection.Close()

			listener.log("start to listen")
			listener.wait.Add(1)
			go listen(listener)
			go healhCheck(listener)
			listener.wait.Wait()
		}
	}()

	return listener
}

func listen(listener *RedisListener) {
	listener.log("listen start")

	for {
		listener.log("before receive")
		var input interface{} = listener.client.Receive()
		listener.log("before the switch", input)

		switch input.(type) {
		case error:
			listener.log("on error", input)

			listener.wait.Done()
			return
		case redis.Message:
			var message = input.(redis.Message)
			listener.log("on message", message.Channel, message.Data)

			for _, handler := range listener.subscribers[message.Channel] {
				handler(message.Data)
			}
		case redis.Subscription:
			var subscription = input.(redis.Subscription)
			listener.log("subscription", subscription)
			switch subscription.Count {
			// Notify application when all channels are subscribed.
			case len(listener.GetChannels()):
				listener.log("subscriptions ready", subscription.Count)

				listener.ready <- true
			}
		default:
			listener.log("something else", input)
		}
	}
}

func healhCheck(listener *RedisListener) {
	listener.log("healhCheck")

	ticker := time.NewTicker(listener.healtCheckTimeout)
	defer ticker.Stop()

	for listener.lastError == nil {
		select {
		case <-ticker.C:
			listener.log("healhCheck ping")

			// TODO: implement reconnect
			var err = listener.client.Ping("")
			if err != nil {
				listener.wait.Done()
			}
		}
	}
}