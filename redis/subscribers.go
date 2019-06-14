package redis

import melody "gopkg.in/olahol/melody.v1"

type operationType int

const (
	get = iota
	put
	removeAll
	exists
)

type request struct {
	operation operationType
	key       string
	value     *OnMessageHandler
	out       chan *OnMessageHandler
	exists    chan bool
	session   *melody.Session
}

// Cache defines cache structure
type Cache struct {
	data    map[string][]*OnMessageHandler
	channel chan request
}

// MakeCache builds a new Cache
func MakeCache() *Cache {
	result := Cache{
		data:    make(map[string][]*OnMessageHandler),
		channel: make(chan request),
	}
	go result.listen()
	return &result
}

// Get returns a channel for the keys's data
func (c *Cache) Get(key string) chan *OnMessageHandler {
	request := request{
		operation: get,
		key:       key,
		out:       make(chan *OnMessageHandler),
	}
	c.channel <- request
	return request.out
}

// Exists returns a channel for the keys's data
func (c *Cache) Exists(key string) chan bool {
	request := request{
		operation: get,
		key:       key,
		exists:    make(chan bool),
	}
	c.channel <- request
	return request.exists
}

// Put stores data for the key
func (c *Cache) Put(key string, value *OnMessageHandler) {
	request := request{
		operation: put,
		key:       key,
		value:     value,
	}
	c.channel <- request
}

// RemoveAll stores data for the key
func (c *Cache) RemoveAll(session *melody.Session) {
	request := request{
		operation: removeAll,
		session:   session,
	}
	c.channel <- request
}

func (c *Cache) listen() {
	for {
		request := <-c.channel
		switch request.operation {
		case get:
			c.send(request)
		case put:
			c.data[request.key] = append(c.data[request.key], request.value)
		case removeAll:
			c.removeAll(request.session)
		case exists:
			request.exists <- len(c.data[request.key]) > 0
		}
	}
}

func (c *Cache) send(r request) {
	for _, line := range c.data[r.key] {
		r.out <- line
	}
	close(r.out)
}

func (c *Cache) removeAll(session *melody.Session) {
	for key, list := range c.data {
		result := make([]*OnMessageHandler, 0)
		for _, handler := range list {
			if handler.Session != session {
				result = append(result, handler)
			}
		}
		c.data[key] = result
	}
}
