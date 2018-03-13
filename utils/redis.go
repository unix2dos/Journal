package utils

import (
	"time"

	"github.com/garyburd/redigo/redis"
)

type RedisStore struct {
	redisPool *redis.Pool
}

func NewRedisStore(host string, password string) *RedisStore {
	var pool = &redis.Pool{
		MaxIdle:     200,
		IdleTimeout: 240 * time.Second,
		Dial: func() (redis.Conn, error) {
			// the redis protocol should probably be made sett-able
			c, err := redis.Dial("tcp", host)
			if err != nil {
				return nil, err
			}
			if len(password) > 0 {
				if _, err := c.Do("AUTH", password); err != nil {
					c.Close()
					return nil, err
				}
			} else {
				// check with PING
				if _, err := c.Do("PING"); err != nil {
					c.Close()
					return nil, err
				}
			}
			return c, err
		},
		// custom connection test method
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			if _, err := c.Do("PING"); err != nil {
				return err
			}
			return nil
		},
	}
	return &RedisStore{pool}
}

func (c *RedisStore) EXISTS(key string) (ok bool, err error) {
	conn := c.redisPool.Get()
	defer conn.Close()
	ok, err = redis.Bool(conn.Do("EXISTS", key))
	if err != nil {
		return
	}
	return
}

func (c *RedisStore) Set(key string, value interface{}) (err error) {
	conn := c.redisPool.Get()
	defer conn.Close()
	_, err = conn.Do("SET", key, value)
	if err != nil {
		return
	}
	return
}

func (c *RedisStore) Get(key string) (s string, err error) {
	conn := c.redisPool.Get()
	defer conn.Close()
	s, err = redis.String(conn.Do("GET", key))
	if err != nil {
		return
	}
	return
}

func (c *RedisStore) HMSet(key string, value interface{}) (err error) {
	conn := c.redisPool.Get()
	defer conn.Close()

	_, err = conn.Do("HMSET", redis.Args{}.Add(key).AddFlat(value)...)
	if err != nil {
		return
	}

	return
}
