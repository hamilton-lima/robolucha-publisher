package redis_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	redis "gitlab.com/robolucha/robolucha-publisher/redis"
	melody "gopkg.in/olahol/melody.v1"
)

func mockHandler(data string) *redis.OnMessageHandler {
	mockSession := melody.Session{Keys: make(map[string]interface{})}
	mockSession.Keys["data"] = data
	mockHandler := redis.OnMessageHandler{Session: &mockSession}
	return &mockHandler
}

func getValue(handler *redis.OnMessageHandler) string {
	return handler.Session.Keys["data"].(string)
}

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
			cache.Put("foo", mockHandler("123"))
			handler := <-cache.Get("foo")
			Expect(getValue(handler)).To(Equal("123"))
		})

		It("should return a single element", func() {
			cache.Put("bar", mockHandler("123"))

			Expect(<-cache.Exists("foo")).To(BeFalse())
			Expect(<-cache.Exists("bar")).To(BeTrue())
		})

		It("should remove all from multiple channels", func() {
			handler := mockHandler("123")
			handler2 := mockHandler("456")
			cache.Put("foo", handler)
			cache.Put("foo", handler2)
			cache.Put("bar", handler)
			cache.Put("bar", handler2)

			cache.RemoveAll(handler.Session)
			// Remove all will not affect other data
			Expect(<-cache.Exists("foo")).To(BeTrue())
			Expect(<-cache.Exists("bar")).To(BeTrue())

			// Only handler 2 should be present on both channels
			for handler := range cache.Get("foo") {
				Expect(getValue(handler)).To(Equal("456"))
			}

			for handler := range cache.Get("bar") {
				Expect(getValue(handler)).To(Equal("456"))
			}
		})

	})
})
