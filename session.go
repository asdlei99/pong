package pong

import (
	"net/http"
	"fmt"
)

// SessionIO define a interface to handle Session's read and write
// pong provide a in memory session manager as default stand by SessionManager interface
// you can define yourself's SessionManager like store session to Redis,MongoDB,File
type SessionIO interface {
	// NewSession should generate a sessionId which is unique compare to existent,and return this sessionId
	// this sessionId string will store in browser by cookies,so the sessionId string should compatible with cookies value rule
	NewSession() (sessionId string)
	// Destory should do operation to remove an session's data in store by give sessionId
	Destory(sessionId string) error
	// Reset should update the give old sessionId to a new id,but the value should be the same
	Reset(oldSessionId string) (newSessionId string, err error)
	// return whether this sessionId is existent in store
	Has(sessionId string) bool
	// read the whole value point to the give sessionId
	Read(sessionId string) (wholeValue map[string]interface{})
	// update the sessionId's value to store
	// the give value just has changed part not all of the value point to sessionId
	Write(sessionId string, changes map[string]interface{}) error
}

type Session struct {
	pong  *Pong
	id    string
	store map[string]interface{}
}

// get the value by name from this session
func (s *Session) Get(name string) interface{} {
	return s.store[name]
}

// set a value with name to this session
// can be used to overwrite old value
func (s *Session) Set(changes map[string]interface{}) error {
	for key, value := range changes {
		s.store[key] = value
	}
	return s.pong.sessionManager.Write(s.id, changes)
}

// update old sessionId with new one
// this will update sessionId store in browser's cookie and session manager's store
func (c *Context) ResetSession() error {
	sessionManager := c.pong.sessionManager
	newId, err := sessionManager.Reset(c.Session.id)
	if err != nil {
		return err
	}
	c.Session.id = newId
	//update sessionID in cookies
	c.Response.Cookie(&http.Cookie{
		HttpOnly: true,
		Name:     SessionCookiesName,
		Value:    newId,
	})
	return nil
}

// remove sessionId
// this will remove sessionId store in browser's cookie and session manager's store
func (c *Context) DestorySession() error {
	err := c.pong.sessionManager.Destory(c.Session.id)
	if err != nil {
		return err
	}
	c.Session = nil
	//delete sessionID in cookies
	c.Response.Cookie(&http.Cookie{
		Name:   SessionCookiesName,
		MaxAge: -1,
	})
	return nil
}

// if you want you HTTP session call this
// EnableSession will use memory to store session data
// EnableSession will read sessionId from request cookies value by name SessionCookiesName default is "SESSIONID" as sessionId
// EnableSession will cause performance drop compare to not use Session
func (pong *Pong) EnableSession(sessionManager SessionIO) {
	if pong.sessionManager != nil {
		fmt.Errorf("sessionManager %v has been set, don't call EnableSession more than once", pong.sessionManager)
		return
	}
	pong.sessionManager = sessionManager
	pong.Root.Middleware(func(c *Context) {
		c.Session = &Session{
			pong:  c.pong,
			store: make(map[string]interface{}),
		}
		sCookie, err := c.Request.HTTPRequest.Cookie(SessionCookiesName)
		if err == nil {
			c.Session.id = sCookie.Value
			if c.pong.sessionManager.Has(c.Session.id) {
				c.Session.store = c.pong.sessionManager.Read(c.Session.id)
			} else {
				goto noSessionID
			}
		} else {
			goto noSessionID
		}
		return
		noSessionID:
		{
			c.Session.id = c.pong.sessionManager.NewSession()
			c.Response.Cookie(&http.Cookie{
				HttpOnly: true,
				Name:     SessionCookiesName,
				Value:    c.Session.id,
			})
		}

	})
}
