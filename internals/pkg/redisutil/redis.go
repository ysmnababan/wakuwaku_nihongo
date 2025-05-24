package redisutil

import (
	"time"

	"github.com/gomodule/redigo/redis"

	"wakuwaku_nihongo/config"
)

type Redis struct {
	client *redis.Pool
}

func NewRedis(cfg config.RedisConfig) *Redis {

	pool := &redis.Pool{
		Dial: func() (redis.Conn, error) {
			if cfg.Password != "" {
				return redis.Dial("tcp", cfg.Address, redis.DialPassword(cfg.Password))
			}
			return redis.Dial("tcp", cfg.Address)
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			if time.Since(t) < time.Minute {
				return nil
			}
			_, err := c.Do("PING")
			return err
		},
		MaxIdle:         cfg.MaxIdle,
		MaxActive:       cfg.MaxActive,
		IdleTimeout:     time.Duration(cfg.IdleTimeout) * time.Second,
		Wait:            true,
		MaxConnLifetime: time.Duration(cfg.MaxConnLifeTime) * time.Second,
	}

	return &Redis{client: pool}
}

func (r *Redis) Do(command string, args ...interface{}) (interface{}, error) {
	conn := r.client.Get()
	defer func(conn redis.Conn) {
		_ = conn.Close()
	}(conn)
	return conn.Do(command, args...)
}

func (r *Redis) Get(key string) (string, error) {
	conn := r.client.Get()
	defer func(conn redis.Conn) {
		_ = conn.Close()
	}(conn)
	return redis.String(conn.Do("GET", key))
}

func (r *Redis) Set(key, value interface{}) (string, error) {
	conn := r.client.Get()
	defer func(conn redis.Conn) {
		_ = conn.Close()
	}(conn)
	return redis.String(conn.Do("SET", key, value))
}

func (r *Redis) SetEX(key string, value interface{}, expire float64) (string, error) {
	conn := r.client.Get()
	defer func(conn redis.Conn) {
		_ = conn.Close()
	}(conn)
	return redis.String(conn.Do("SETEX", key, expire, value))
}

func (r *Redis) TTL(key string) (int, error) {
	conn := r.client.Get()
	defer func(conn redis.Conn) {
		_ = conn.Close()
	}(conn)
	return redis.Int(conn.Do("TTL", key))
}

func (r *Redis) Del(key string) (int64, error) {
	conn := r.client.Get()
	defer func(conn redis.Conn) {
		_ = conn.Close()
	}(conn)
	return redis.Int64(conn.Do("DEL", key))
}

func (r *Redis) Expire(key string, ttl int) (int64, error) {
	conn := r.client.Get()
	defer func(conn redis.Conn) {
		_ = conn.Close()
	}(conn)
	return redis.Int64(conn.Do("EXPIRE", key, ttl))
}

func (r *Redis) Keys(pattern string) ([]string, error) {
	conn := r.client.Get()
	defer func(conn redis.Conn) {
		_ = conn.Close()
	}(conn)
	return redis.Strings(conn.Do("KEYS", pattern))
}

func (r *Redis) FlushAll() error {
	conn := r.client.Get()
	defer func(conn redis.Conn) {
		_ = conn.Close()
	}(conn)
	_, err := conn.Do("FLUSHALL")
	return err
}
