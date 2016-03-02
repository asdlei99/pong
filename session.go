package pong

import (
	"crypto/rand"
	"encoding/base64"
	"io"
	"net/http"
)

// SessionManager define a interface to handle Session's read and write
// pong provide a in memory session manager as default stand by SessionManager interface
// you can define yourself's SessionManager like store session to Redis,MongoDB,File
type SessionManager interface {
	// NewSession should generate a sessionId which is unique compare to existent,and return this sessionId
	// this sessionId string will store in browser by cookies,so the sessionId string should compatible with cookies value rule
	NewSession() (sessionId string)
	// NewSession should do operation to remove an session's data in store by give sessionId
	Destory(sessionId string)
	// Reset should update the give old sessionId to a new id,but the value should be the same
	Reset(oldSessionId string) (newSessionId string)
	// return whether this sessionId is existent in store
	Has(sessionId string) bool
	// read the value point to the give sessionId
	Read(sessionId string) (value map[string]interface{})
	// update the sessionId's value to store
	// the give value just has changed part not all of the value point to sessionId
	Write(sessionId string, value map[string]interface{})
}

type Session struct {
	pong               *Pong
	id                 string
	idHasChange        bool
	store              map[string]interface{}
	hasChangeValueFlag []string
}

// get the value by name
func (s *Session) Get(name string) interface{} {
	return s.store[name]
}

// set a value with name
// can be used to overwrite old value
func (s *Session) Set(name string, value interface{}) {
	s.store[name] = value
	s.hasChangeValueFlag = append(s.hasChangeValueFlag, name)
}

// update old sessionId with new one
// this will update sessionId store in browser's cookie and session manager's store
func (s *Session) Reset() {
	sessionManager := s.pong.SessionManager
	s.id = sessionManager.Reset(s.id)
	s.idHasChange = true
}

// remove sessionId
// this will remove sessionId store in browser's cookie and session manager's store
func (s *Session) Destory() {
	s.pong.SessionManager.Destory(s.id)
	s.id = ""
	s.idHasChange = true
}

//default in memory sessionManager
type memorySessionManager struct {
	SessionManager
	store map[string]map[string]interface{}
}

func (manager *memorySessionManager) NewSession() (sessionId string) {
	bs := make([]byte, 8)
	io.ReadFull(rand.Reader, bs)
	sessionId = base64.URLEncoding.EncodeToString(bs)
	if manager.Has(sessionId) {
		return manager.NewSession()
	}
	manager.store[sessionId] = make(map[string]interface{})
	return sessionId
}

func (manager *memorySessionManager) Destory(sessionId string) {
	delete(manager.store, sessionId)
}

func (manager *memorySessionManager) Reset(oldSessionId string) (newSessionId string) {
	newSessionId = manager.NewSession()
	manager.store[newSessionId] = manager.store[oldSessionId]
	delete(manager.store, oldSessionId)
	return
}

func (manager *memorySessionManager) Has(sessionId string) bool {
	return manager.store[sessionId] != nil
}

func (manager *memorySessionManager) Read(sessionId string) map[string]interface{} {
	return manager.store[sessionId]
}

func (manager *memorySessionManager) Write(sessionId string, value map[string]interface{}) {
	manager.store[sessionId] = value
}

func (pong *Pong) EnableSession() {
	if pong.SessionManager == nil {
		pong.SessionManager = &memorySessionManager{
			store: make(map[string]map[string]interface{}),
		}
	}
	pong.Root.Middleware(func(c *Context) {
		c.Session = &Session{
			pong:  c.pong,
			store: make(map[string]interface{}),
		}
		sCookie, err := c.Request.HTTPRequest.Cookie(SessionCookiesName)
		if err == nil {
			c.Session.id = sCookie.Value
			if c.pong.SessionManager.Has(c.Session.id) {
				c.Session.store = c.pong.SessionManager.Read(c.Session.id)
			} else {
				goto noSessionID
			}
		} else {
			goto noSessionID
		}
		return
	noSessionID:
		{
			c.Session.id = c.pong.SessionManager.NewSession()
			c.Response.Cookie(&http.Cookie{
				HttpOnly: true,
				Name:     SessionCookiesName,
				Value:    c.Session.id,
			})
		}

	})
	pong.TailMiddleware(func(c *Context) {
		if c.Session.idHasChange {
			if len(c.Session.id) > 0 {
				//update sessionID in cookies
				c.Response.Cookie(&http.Cookie{
					HttpOnly: true,
					Name:     SessionCookiesName,
					Value:    c.Session.id,
				})
			} else {
				//delete sessionID in cookies
				c.Response.Cookie(&http.Cookie{
					Name:   SessionCookiesName,
					MaxAge: -1,
				})
			}
		}
		if len(c.Session.hasChangeValueFlag) > 0 {
			change := make(map[string]interface{})
			for _, name := range c.Session.hasChangeValueFlag {
				change[name] = c.Session.store[name]
			}
			c.pong.SessionManager.Write(c.Session.id, change)
		}
	})
}
