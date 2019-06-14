package session

import (
	"net/http"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	melody "gopkg.in/olahol/melody.v1"
)

var _ = Describe("Sessionmanager", func() {

	Context("when creating the session manager", func() {
		var manager *SessionManager

		BeforeEach(func() {
			manager = NewSessionManager()
		})

		It("should have empty list of sessions", func() {
			var sessions = manager.GetSessions("")
			Expect(len(sessions)).To(Equal(0))
		})

		It("should have empty list of sessions for different matchIDs", func() {
			var sessions = manager.GetSessions("32")
			Expect(len(sessions)).To(Equal(0))

			sessions = manager.GetSessions("42")
			Expect(len(sessions)).To(Equal(0))
		})

		var globalMatchID = "42"

		var createSessionWithRequest = func() *melody.Session {
			var session = &melody.Session{}
			req, _ := http.NewRequest("GET", "/match/"+globalMatchID, nil)
			session.Request = req
			return session
		}

		It("should return the proper matchID from the URL", func() {
			var session = createSessionWithRequest()
			var matchID = manager.GetIDFromURL(session.Request.URL)
			Expect(matchID).To(Equal(globalMatchID))
		})

		It("should have 1 element after adding to the list", func() {
			manager.AddSession(createSessionWithRequest())
			var sessions = manager.GetSessions(globalMatchID)
			Expect(len(sessions)).To(Equal(1))
		})

		It("should have 2 element after adding to the list", func() {
			manager.AddSession(createSessionWithRequest())
			manager.AddSession(createSessionWithRequest())

			var sessions = manager.GetSessions(globalMatchID)
			Expect(len(sessions)).To(Equal(2))
		})

		It("should have 1 element after adding 2 and removing 1 from the list", func() {
			var session1 = createSessionWithRequest()
			var session2 = createSessionWithRequest()
			manager.AddSession(session1)
			manager.AddSession(session2)
			manager.RemoveSession(session1)

			var sessions = manager.GetSessions(globalMatchID)
			Expect(len(sessions)).To(Equal(1))
		})

		It("should have 0 element after adding 2 and removing 2 from the list", func() {
			var session1 = createSessionWithRequest()
			var session2 = createSessionWithRequest()
			manager.AddSession(session1)
			manager.AddSession(session2)
			manager.RemoveSession(session1)
			manager.RemoveSession(session2)

			var sessions = manager.GetSessions(globalMatchID)
			Expect(len(sessions)).To(Equal(0))
		})

		It("nextID should generate values in sequence", func() {
			Expect(manager.NextID()).To(Equal(uint64(1)))
			Expect(manager.NextID()).To(Equal(uint64(2)))
			Expect(manager.NextID()).To(Equal(uint64(3)))
		})

		It("should have proper ID in the sessions after adding", func() {
			var session1 = createSessionWithRequest()
			var session2 = createSessionWithRequest()
			var session3 = createSessionWithRequest()

			Expect(manager.GetIDFromSession(session1)).To(Equal(uint64(0)))
			Expect(manager.GetIDFromSession(session2)).To(Equal(uint64(0)))
			Expect(manager.GetIDFromSession(session3)).To(Equal(uint64(0)))

			manager.AddSession(session1)
			manager.AddSession(session2)
			manager.AddSession(session3)

			Expect(manager.GetIDFromSession(session1)).To(Equal(uint64(1)))
			Expect(manager.GetIDFromSession(session2)).To(Equal(uint64(2)))
			Expect(manager.GetIDFromSession(session3)).To(Equal(uint64(3)))
		})

		It("should have proper ID when setting to Session", func() {
			var session1 = createSessionWithRequest()
			var ID = manager.NextID()
			session1.Set("ID", ID)
			var readID = manager.GetIDFromSession(session1)
			Expect(readID).To(Equal(ID))
		})

		It("should have valid elements after adding and removing from the list", func() {
			var session1 = createSessionWithRequest()
			var session2 = createSessionWithRequest()
			manager.AddSession(session1)
			manager.AddSession(session2)
			manager.RemoveSession(session1)

			var sessions = manager.GetSessions(globalMatchID)
			var session = sessions[0]
			Expect(session).ShouldNot(BeNil())

		})

	})

})
