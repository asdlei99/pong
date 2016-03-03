package redis_session

import (
	"io"
	"encoding/base64"
	"crypto/rand"
	"github.com/gwuhaolin/pong"
	"gopkg.in/redis.v3"
)

//redis sessionManager
type redisSessionManager struct {
	pong.SessionIO
	redisClient *redis.Client
}

func New(options *redis.Options) pong.SessionIO {
	return &redisSessionManager{
		redisClient:redis.NewClient(options),
	}
}

func (manager *redisSessionManager) NewSession() (sessionId string) {
	bs := make([]byte, 8)
	io.ReadFull(rand.Reader, bs)
	sessionId = base64.URLEncoding.EncodeToString(bs)
	if manager.Has(sessionId) {
		return manager.NewSession()
	}
	return sessionId
}

func (manager *redisSessionManager) Destory(sessionId string) error {
	return manager.redisClient.Del(sessionId).Err()
}

func (manager *redisSessionManager) Reset(oldSessionId string) (newSessionId string, err error) {
	newSessionId = manager.NewSession()
	err = manager.redisClient.Rename(oldSessionId, newSessionId).Err()
	return
}

func (manager *redisSessionManager) Has(sessionId string) bool {
	has, _ := manager.redisClient.Exists(sessionId).Result()
	return has
}

func (manager *redisSessionManager) Read(sessionId string) (wholeValue map[string]interface{}) {
	//TODO
	manager.redisClient.HMGet(sessionId).Result()
	return
}

func (manager *redisSessionManager) Write(sessionId string, changes map[string]interface{}) error {
	//TODO
	//for k, v := range changes {
		//manager.redisClient.HMSet(sessionId, k, v)
	//}
	return nil
}