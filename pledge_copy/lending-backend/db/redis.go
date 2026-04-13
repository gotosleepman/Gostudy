package db

import (
	"encoding/json"
	"fmt"
	"time"

	"lending-copy/config"
	"lending-copy/log"

	"github.com/gomodule/redigo/redis"
)

func InitRedis() *redis.Pool {
	log.Logger.Info("Init Redis")
	redisConf := config.Config.Redis
	RedisConn = &redis.Pool{
		MaxIdle:     10,
		MaxActive:   0,
		Wait:        true,
		IdleTimeout: 180 * time.Second,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", fmt.Sprintf("%s:%s", redisConf.Address, redisConf.Port))
			if err != nil {
				return nil, err
			}
			if redisConf.Password != "" {
				if _, err = c.Do("auth", redisConf.Password); err != nil {
					_ = c.Close()
					return nil, err
				}
			}
			if _, err = c.Do("select", redisConf.Db); err != nil {
				_ = c.Close()
				return nil, err
			}
			return c, nil
		},
	}
	_ = RedisConn.Get().Close()
	return RedisConn
}

func RedisSet(key string, data interface{}, aliveSeconds int) error {
	conn := RedisConn.Get()
	defer func() { _ = conn.Close() }()
	value, err := json.Marshal(data)
	if err != nil {
		return err
	}
	if aliveSeconds > 0 {
		_, err = conn.Do("set", key, value, "EX", aliveSeconds)
	} else {
		_, err = conn.Do("set", key, value)
	}
	return err
}

func RedisGet(key string) ([]byte, error) {
	conn := RedisConn.Get()
	defer func() { _ = conn.Close() }()
	return redis.Bytes(conn.Do("get", key))
}

func RedisFlushDB() error {
	conn := RedisConn.Get()
	defer func() { _ = conn.Close() }()
	_, err := conn.Do("flushdb")
	return err
}

func RedisSetString(key string, data string, aliveSeconds int) error {
	conn := RedisConn.Get()
	defer func() { _ = conn.Close() }()
	var err error
	if aliveSeconds > 0 {
		_, err = conn.Do("set", key, data, "EX", aliveSeconds)
	} else {
		_, err = conn.Do("set", key, data)
	}
	return err
}
