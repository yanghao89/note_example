package redisstore

import (
	"sync"
)

var (
	dbStore    *RedisStore
	storeMutex sync.Mutex
)

func Init() {
	_ = CreateStore()
}

func CreateStore() *RedisStore {
	if dbStore == nil {
		storeMutex.Lock()
		defer storeMutex.Unlock()
		if dbStore == nil {
			redisStore := &RedisStore{}
			redisStore.RedisPool = redisStore.NewConn(&Conf{
				RedisHost: "127.0.0.1:6379",
				RedisDB:   "0",
				RedisPwd:  "",
				Timeout:   20,
				MaxActive: 0,
				MaxIdle:   2,
			})
			dbStore = redisStore
		}
	}
	return dbStore
}
