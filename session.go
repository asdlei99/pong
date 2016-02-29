package pong

import (
	"io"
	"crypto/rand"
	"encoding/base64"
)

type SessionManager interface {
	NewSession() (sessionId string)
	Destory(sessionId string)
	Has(sessionId string) bool
	Read(sessionId string) (value map[string]interface{})
	Write(sessionId string, value map[string]interface{})
}

type Session struct {
	pong          *Pong
	id            string
	store         map[string]interface{}
	hasChangeFlag []string
}

func (s *Session)Get(name string) interface{} {
	return s.store[name]
}

func (s *Session)Set(name string, value interface{}) {
	s.store[name] = value
	s.hasChangeFlag = append(s.hasChangeFlag, name)
}

func (s *Session)Reset() {
	sessionManager := s.pong.SessionManager
	newId := sessionManager.NewSession()
	sessionManager.Write(newId, s.store)
	oldId := s.id
	s.id = newId
	sessionManager.Destory(oldId)
}

func (s *Session)Destory() {
	s.pong.SessionManager.Destory(s.id)
}

//default in memory sessionManager
type memorySessionManager struct {
	SessionManager
	store map[string]map[string]interface{}
}

func (manager *memorySessionManager)NewSession() (sessionId string) {
	bs := make([]byte, 32)
	io.ReadFull(rand.Reader, bs)
	sessionId = base64.URLEncoding.EncodeToString(bs)
	manager.store[sessionId] = make(map[string]interface{})
	return sessionId
}

func (manager *memorySessionManager)Destory(sessionId string) {
	delete(manager.store, sessionId)
}

func (manager *memorySessionManager)Has(sessionId string) bool {
	return manager.store[sessionId] != nil
}

func (manager *memorySessionManager)Read(sessionId string) map[string]interface{} {
	return manager.store[sessionId]
}

func (manager *memorySessionManager)Write(sessionId string, value map[string]interface{}) {
	manager.store[sessionId] = value
}

