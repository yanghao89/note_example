package redisstore

import (
	"time"

	"github.com/garyburd/redigo/redis"
)

type RedisStore struct {
	RedisPool *redis.Pool
}

type Conf struct {
	RedisHost string
	RedisDB   string
	RedisPwd  string
	Timeout   int64
	MaxIdle   int
	MaxActive int
}

//NewConn 链接Redis
func (r *RedisStore) NewConn(conf *Conf) *redis.Pool {
	return &redis.Pool{
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", conf.RedisHost)
			if err != nil {
				return nil, err
			}
			if conf.RedisPwd != "" {
				if _, err := c.Do("AUTH", conf.RedisPwd); err != nil {
					return nil, err
				}
			}
			if conf.RedisDB != "" {
				if _, err := c.Do("SELECT", conf.RedisDB); err != nil {
					return nil, err
				}
			}
			timeOut := time.Duration(conf.Timeout) * time.Second
			redis.DialConnectTimeout(timeOut)
			redis.DialReadTimeout(timeOut)
			redis.DialWriteTimeout(timeOut)
			return c, nil
		},
		MaxIdle:     conf.MaxIdle,
		MaxActive:   conf.MaxActive,
		IdleTimeout: 1 * time.Second,
		Wait:        true,
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
	}
}

//Do 执行redis命令
func (r *RedisStore) Do(commandName string, args ...interface{}) (interface{}, error) {
	c := r.RedisPool.Get()
	defer c.Close()
	return c.Do(commandName, args...)
}
