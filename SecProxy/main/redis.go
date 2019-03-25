package main

import (
	"github.com/garyburd/redigo/redis"
	"time"
)

func initRedis() (err error) {
	RedisPoll = &redis.Pool{
		MaxIdle:     SecKillConfig.RedisBlackConfig.RedisMaxIdle,
		MaxActive:   SecKillConfig.RedisBlackConfig.RedisMaxActive,
		IdleTimeout: time.Duration(SecKillConfig.RedisBlackConfig.RedisIdleTimeout) * time.Second,
		Dial: func() (conn redis.Conn, e error) {
			return redis.Dial("tcp", SecKillConfig.RedisBlackConfig.RedisAddr)
		},
	}

	conn := RedisPoll.Get()
	defer conn.Close()
	_, err = conn.Do("ping")
	if err != nil {
		return err
	}
	return nil
}
