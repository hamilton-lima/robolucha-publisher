package redis_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	redis "gitlab.com/robolucha/robolucha-publisher/redis"
	melody "gopkg.in/olahol/melody.v1"
)

var _ = Describe("Redis", func() {

	Context("when adding data to the cache", func() {
		var cache *redis.Cache

		BeforeEach(func() {
			cache = redis.MakeCache()
		})

		It("should be empty", func() {
			exist := <-cache.Exists("foo")
			Expect(exist).To(BeFalse())
		})

		It("should return a single element", func() {
			mockSession := melody.Session{Keys: make(map[string]interface{})}
			mockSession.Keys["some"] = "data"
			mockHandler := redis.OnMessageHandler{Session: &mockSession}
			cache.Put("foo", &mockHandler)
			handler := <-cache.Get("foo")
			Expect(handler.Session.Keys["some"]).To(Equal("data"))
		})

	})
})
