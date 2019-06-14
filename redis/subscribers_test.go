package redis

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Subscribers cache", func() {

	Context("when adding data to the cache", func() {
		var cache *Cache

		BeforeEach(func() {
			cache = MakeCache()
		})

		It("should be empty", func() {
			exist := cache.Exists("foo")
			Expect(exist).To(BeFalse())
		})

	})

})
