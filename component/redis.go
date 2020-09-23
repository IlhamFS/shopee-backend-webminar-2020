package component

import (
	"github.com/gomodule/redigo/redis"
)

type redisCon struct {
	redis redis.Conn
}

type Redis interface {
	Get(key string) (result string, err error)
	Set(key string, value string) (err error)
}

func NewRedis(redis redis.Conn) Redis {
	return &redisCon{
		redis: redis,
	}
}

func (r *redisCon) Get(key string) (result string, err error) {
	return redis.String(r.redis.Do("GET", key))
}

func (r *redisCon) Set(key, value string) (err error) {
	_, err = r.redis.Do("SET", key, value)
	return
}

func InitializeRedis() (redisCon Redis, err error) {
	c, err := redis.Dial("tcp", "localhost:6379")
	if err != nil {
		return nil, err
	}
	return NewRedis(c), nil
}
