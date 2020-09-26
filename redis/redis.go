package redis

import (
	"errors"
	"github.com/gomodule/redigo/redis"
	"time"
)

var (
	ErrNX                   = errors.New("key already exist")
	ErrInsufficientArgument = errors.New("wrong number of arguments")
)

type redisCon struct {
	redis redis.Conn
}

type RedisInterface interface {
	Get(key string) (result string, err error)
	Set(key string, value string) (err error)
	SetNX(key, value string, ttl time.Duration) error
	Del(key ...string) error
}

func NewRedis(redis redis.Conn) RedisInterface {
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

func (r *redisCon) Del(key ...string) error {
	if len(key) == 0 {
		return ErrInsufficientArgument
	}
	_, err := r.redis.Do("DEL", redis.Args{}.AddFlat(key)...)
	return err
}

func (r *redisCon) SetNX(key, value string, ttl time.Duration) error {
	var (
		err   error
		reply interface{}
	)

	if 0 >= ttl {
		reply, err = r.redis.Do("SET", key, value, "NX")
	} else {
		reply, err = r.redis.Do("SET", key, value, "EX", int64(ttl.Seconds()), "NX")
	}

	if nil != err {
		return err
	}
	if nil == reply {
		return ErrNX
	}

	return err
}

func InitializeRedis() (redisCon RedisInterface, err error) {
	c, err := redis.Dial("tcp", "localhost:6379")
	if err != nil {
		return nil, err
	}
	return NewRedis(c), nil
}
