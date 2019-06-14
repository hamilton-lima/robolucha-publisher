package redis

import (
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/gomodule/redigo/redis"
	log "github.com/sirupsen/logrus"
	melody "gopkg.in/olahol/melody.v1"
)

// OnMessageHandler defines function to be executed when a new message arrives
type OnMessageHandler struct {
	Session *melody.Session
	Handler func(session *melody.Session, data []byte)
}

// Listener is the listener itself
type Listener struct {
	serverAddr        string
	subscribers       *Cache
	connection        redis.Conn
	client            redis.PubSubConn
	wait              sync.WaitGroup
	readTimeout       time.Duration
	writeTimeout      time.Duration
	healtCheckTimeout time.Duration
	lastError         error
	isInfo            bool
	isDebug           bool
}

// NewRedisListener creates a new RedisListener
func NewRedisListener() *Listener {

	host := os.Getenv("REDIS_HOST")
	port := os.Getenv("REDIS_PORT")
	serverAddr := fmt.Sprintf("%v:%v", host, port)

	var result = Listener{}
	result.serverAddr = serverAddr
	result.subscribers = MakeCache()
	result.healtCheckTimeout = time.Minute * 5
	result.readTimeout = result.healtCheckTimeout + (10 * time.Second)
	result.writeTimeout = 10 * time.Second
	result.isInfo = true
	result.isDebug = false

	log.WithFields(log.Fields{
		"serverAddr": serverAddr,
	}).Info("Connect to REDIS Configuration")

	return &result
}

// Subscribe adds OnMessageHandler for one channel
func (listener *Listener) Subscribe(channel string, handler *OnMessageHandler) *Listener {
	log.WithFields(log.Fields{
		"channel": channel,
	}).Info("Subscribe")

	exists := <-listener.subscribers.Exists(channel)
	if !exists {
		log.WithFields(log.Fields{
			"channel": channel,
		}).Info("First subscriber to channel")

		listener.lastError = listener.client.Subscribe(channel)

		if listener.lastError != nil {
			log.WithFields(log.Fields{
				"error":   listener.lastError,
				"channel": channel,
			}).Error("Subscribing to channel")
		}
	}

	listener.subscribers.Put(channel, handler)
	return listener
}

// UnSubscribeAll based on the session
func (listener *Listener) UnSubscribeAll(session *melody.Session) *Listener {
	listener.subscribers.RemoveAll(session)
	return listener
}

// SetDebugger change isDebug setting
func (listener *Listener) SetDebugger(isDebug bool) *Listener {
	log.WithFields(log.Fields{
		"isDebug": isDebug,
	}).Info("SetDebugger")

	listener.isDebug = isDebug
	return listener
}

// // GetChannels returns the list of subscribed channels
// func (listener *RedisListener) GetChannels() []string {
// 	channels := make([]string, len(listener.subscribers))

// 	i := 0
// 	for key := range listener.subscribers {
// 		channels[i] = key
// 		i++
// 	}
// 	return channels
// }

// Connect establish the connection to Redis and listen to the subscribed channels
func (listener *Listener) Connect() *Listener {
	log.Info("Connecting to REDIS")

	listener.connection, listener.lastError = redis.Dial("tcp", listener.serverAddr,
		redis.DialReadTimeout(listener.readTimeout),
		redis.DialWriteTimeout(listener.writeTimeout))

	if listener.lastError != nil {
		log.WithFields(log.Fields{
			"error": listener.lastError,
		}).Error("Error Connecting to REDIS")

		return listener
	}

	listener.client = redis.PubSubConn{Conn: listener.connection}

	go func() {
		// TODO: will allow retries, remove for now
		for {
			defer listener.client.Unsubscribe()
			defer listener.connection.Close()

			log.Info("Connected to REDIS, listening messages")
			listener.wait.Add(1)
			go listen(listener)
			go healhCheck(listener)
			listener.wait.Wait()
		}
	}()

	return listener
}

func listen(listener *Listener) {
	log.Info("REDIS listener start")

	for {
		log.Debug("Waiting for REDIS")
		var input interface{} = listener.client.Receive()

		switch input.(type) {
		case error:
			log.WithFields(log.Fields{
				"message": input,
			}).Info("onError")

			listener.wait.Done()
			return
		case redis.Message:
			var message = input.(redis.Message)

			if log.IsLevelEnabled(log.DebugLevel) {
				var data = string(message.Data)

				log.WithFields(log.Fields{
					"channel": message.Channel,
					"message": data,
					"length":  len(data),
				}).Debug("onMessage")
			}

			for onMessageHandler := range listener.subscribers.Get(message.Channel) {
				onMessageHandler.Handler(onMessageHandler.Session, message.Data)
			}

		case redis.Subscription:
			var subscription = input.(redis.Subscription)
			log.WithFields(log.Fields{
				"subscription": subscription,
			}).Debug("subscription")

		default:
			log.WithFields(log.Fields{
				"input": input,
			}).Debug("something else")
		}
	}
}

func healhCheck(listener *Listener) {
	log.Info("healhCheck")

	ticker := time.NewTicker(listener.healtCheckTimeout)
	defer ticker.Stop()

	for listener.lastError == nil {
		select {
		case <-ticker.C:
			log.Info("healthCheck ping")

			// TODO: implement reconnect
			var err = listener.client.Ping("")
			if err != nil {
				panic(1)
			}
		}
	}
}
