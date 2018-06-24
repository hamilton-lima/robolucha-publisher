package sessionmanager

import (
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
	var lastSessionID uint64
	result.lastSessionID = &lastSessionID

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
	session.Set("ID", manager.nextID())
	manager.lock.Lock()
	(*manager.sessions)[matchID] = append((*manager.sessions)[matchID], session)
	manager.lock.Unlock()
}

// RemoveSession removes a session from the current list
func (manager *SessionManager) RemoveSession(session *melody.Session) {
	// TODO: FIX THIS

	var matchID = manager.GetIDFromURL(session.Request.URL)
	var ID2Remove, _ = session.Get("ID")

	manager.lock.Lock()
	var size = len((*manager.sessions)[matchID]) - 1
	var updatedSessions = make([]*melody.Session, size, size)
	var pos = 0
	for _, slice := range (*manager.sessions)[matchID] {
		var ID, _ = slice.Get("ID")
		if *ID.(*uint64) != *ID2Remove.(*uint64) {
			updatedSessions[pos] = slice
			pos++
		}
	}

	(*manager.sessions)[matchID] = updatedSessions

	manager.lock.Unlock()
}

// nextID calculates the next session ID
func (manager *SessionManager) nextID() *uint64 {
	atomic.AddUint64(manager.lastSessionID, 1)
	return manager.lastSessionID
}
