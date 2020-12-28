package lock

import (
	"example/pkg/redisstore"

	"github.com/garyburd/redigo/redis"
)

type DistributedLock struct {
	lockKey string `json:"lock_key"`
	timeout int    `json:"timeout"`
}

func NewDistributedLock() *DistributedLock {
	return &DistributedLock{
		lockKey: "lock:test",
		timeout: 30,
	}
}

func (d *DistributedLock) TryLock() (bool, error) {
	_, err := redis.String(redisstore.CreateStore().Do("SET", d.lockKey, 1, "EX", d.timeout, "NX"))
	if err != nil {
		if err == redis.ErrNil {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func (d *DistributedLock) UnLock() error {
	_, err := redisstore.CreateStore().Do("DEL", d.lockKey)
	return err
}
