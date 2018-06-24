package sessionmanager

import (
	"fmt"
	"net/url"
	"strings"
	"sync"
	"sync/atomic"

	"gopkg.in/olahol/melody.v1"
)

// SessionManager handle the list of active sessions updated
type SessionManager struct {
	sessions      *map[string][]*melody.Session
	lock          *sync.Mutex
	lastSessionID *uint64
}

// NewSessionManager creates a new SessionManager
func NewSessionManager() *SessionManager {
	var result = SessionManager{}
	var m = make(map[string][]*melody.Session)
	result.sessions = &m
	result.lock = &sync.Mutex{}
	result.lastSessionID = new(uint64)

	return &result
}

// GetIDFromURL extract the matchID from the last element of the URL
func (manager *SessionManager) GetIDFromURL(URL *url.URL) string {
	sliced := strings.Split(URL.Path, "/")
	last := sliced[len(sliced)-1]
	return last
}

// GetSessions return the list of active sessions
func (manager *SessionManager) GetSessions(matchID string) []*melody.Session {
	return (*manager.sessions)[matchID]
}

// AddSession add a new session
func (manager *SessionManager) AddSession(session *melody.Session) {
	var matchID = manager.GetIDFromURL(session.Request.URL)
	session.Set("ID", manager.NextID())
	manager.lock.Lock()
	(*manager.sessions)[matchID] = append((*manager.sessions)[matchID], session)
	manager.lock.Unlock()
}

// RemoveSession removes a session from the current list
// TODO: use another structure to support several removal operations
// and regenerate the session array from it
func (manager *SessionManager) RemoveSession(session *melody.Session) {
	var matchID = manager.GetIDFromURL(session.Request.URL)
	var ID2Remove = manager.GetIDFromSession(session)

	manager.lock.Lock()
	var size = len((*manager.sessions)[matchID]) - 1
	var updatedSessions []*melody.Session

	if size > 0 {
		updatedSessions = make([]*melody.Session, size, size)
		for _, slice := range (*manager.sessions)[matchID] {
			var ID = manager.GetIDFromSession(slice)
			fmt.Printf("ID=%v ID2Remove=%v,", ID, ID2Remove)

			if ID != ID2Remove {
				updatedSessions = append(updatedSessions, slice)
			}
		}
	} else {
		updatedSessions = make([]*melody.Session, 0, 0)
	}

	(*manager.sessions)[matchID] = updatedSessions
	manager.lock.Unlock()
}

// GetIDFromSession returns the ID from the Session
func (manager *SessionManager) GetIDFromSession(session *melody.Session) uint64 {
	var result, _ = session.Get("ID")
	var pointer = result.(*uint64)
	return *pointer
}

// NextID calculates the next session ID
func (manager *SessionManager) NextID() *uint64 {
	atomic.AddUint64(manager.lastSessionID, 1)
	return manager.lastSessionID
}
