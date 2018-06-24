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
	sessions      *map[string]map[uint64]*melody.Session
	lock          *sync.Mutex
	lastSessionID *uint64
}

// NewSessionManager creates a new SessionManager
func NewSessionManager() *SessionManager {
	var result = SessionManager{}
	var m = make(map[string]map[uint64]*melody.Session)
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

// Broadcast send message to the current sessions of the matchID
func (manager *SessionManager) Broadcast(matchID string, message []byte) {
	manager.lock.Lock()
	for key := range (*manager.sessions)[matchID] {
		var session = (*manager.sessions)[matchID][key]
		go session.Write(message)
	}
	manager.lock.Unlock()
}

// GetSessions should be used only for tests
func (manager *SessionManager) GetSessions(matchID string) []*melody.Session {
	var result []*melody.Session
	manager.lock.Lock()
	for key := range (*manager.sessions)[matchID] {
		var session = (*manager.sessions)[matchID][key]
		fmt.Printf("getsession %v, %p \n", matchID, session)
		result = append(result, session)
	}
	manager.lock.Unlock()
	fmt.Printf("getsession %v \n", result)
	return result
}

// AddSession add a new session
func (manager *SessionManager) AddSession(session *melody.Session) {
	var matchID = manager.GetIDFromURL(session.Request.URL)
	manager.lock.Lock()
	var ID = manager.NextID()
	session.Set("ID", ID)

	if (*manager.sessions)[matchID] == nil {
		(*manager.sessions)[matchID] = make(map[uint64]*melody.Session)
	}

	(*manager.sessions)[matchID][ID] = session
	fmt.Printf("after adding to session %v \n", (*manager.sessions)[matchID])
	manager.lock.Unlock()
}

// RemoveSession removes a session from the current list
func (manager *SessionManager) RemoveSession(session *melody.Session) {
	var matchID = manager.GetIDFromURL(session.Request.URL)
	var ID2Remove = manager.GetIDFromSession(session)

	manager.lock.Lock()
	delete((*manager.sessions)[matchID], ID2Remove)
	manager.lock.Unlock()
}

// GetIDFromSession returns the ID from the Session
func (manager *SessionManager) GetIDFromSession(session *melody.Session) uint64 {
	var result, _ = session.Get("ID")
	if result == nil {
		return uint64(0)
	}

	return result.(uint64)
}

// NextID calculates the next session ID
func (manager *SessionManager) NextID() uint64 {
	atomic.AddUint64(manager.lastSessionID, 1)
	return *manager.lastSessionID
}
