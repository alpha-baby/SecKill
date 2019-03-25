package service

import (
	"errors"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/garyburd/redigo/redis"
	"strconv"
	"time"
)

var (
	secKillServer *SecKillServer
)

func InitService(serviceConf *SecKillServer) (err error) {
	secKillServer = serviceConf
	beego.Debug("init service success")

	secKillServer.secLimitMgr = &SecLimitMgr{
		UserLimitMap: make(map[int]*SecLimit, 10000),
		IpLimitMap:   make(map[string]*SecLimit, 10000),
	}

	// 从redis 中读取黑名单信息
	err = loadBlackList()
	if err != nil {
		return fmt.Errorf("load black list error:%v", err)
	}
	beego.Debug("init load black list success")

	// 初始化两个redis连接池
	{
		err = initProxy2LayerRedis()
		if err != nil {
			return errors.New(fmt.Sprintf("init Proxy 2 Layer Redis pool error:%v", err))
		}
		err = initLayer2ProxyRedis()
		if err != nil {
			return errors.New(fmt.Sprintf("init Layer 2 Proxy Redis pool error:%v", err))
		}
	}
	beego.Debug("init Proxy2Layer Redis pool success")

	secKillServer.SecReqChan = make(chan *SecRequest, secKillServer.SecReqChanSize)
	initRedisProcessFunc()

	return
}

func initProxy2LayerRedis() (err error) {
	secKillServer.Proxy2LayerRedisPool = &redis.Pool{
		MaxIdle:     secKillServer.RedisProxy2LayerConf.RedisMaxIdle,
		MaxActive:   secKillServer.RedisProxy2LayerConf.RedisMaxActive,
		IdleTimeout: time.Duration(secKillServer.RedisProxy2LayerConf.RedisIdleTimeout) * time.Second,
		Dial: func() (conn redis.Conn, e error) {
			return redis.Dial("tcp", secKillServer.RedisProxy2LayerConf.RedisAddr)
		},
	}

	conn := secKillServer.Proxy2LayerRedisPool.Get()
	defer conn.Close()
	_, err = conn.Do("ping")
	if err != nil {
		return err
	}
	return nil
}

func initLayer2ProxyRedis() (err error) {
	secKillServer.Layer2ProxyRedisPool = &redis.Pool{
		MaxIdle:     secKillServer.RedisLayer2ProxyConf.RedisMaxIdle,
		MaxActive:   secKillServer.RedisLayer2ProxyConf.RedisMaxActive,
		IdleTimeout: time.Duration(secKillServer.RedisLayer2ProxyConf.RedisIdleTimeout) * time.Second,
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", secKillServer.RedisLayer2ProxyConf.RedisAddr)
		},
	}

	conn := secKillServer.Layer2ProxyRedisPool.Get()
	defer conn.Close()

	_, err = conn.Do("ping")
	if err != nil {
		logs.Error("ping redis failed, err:%v", err)
		return
	}

	return
}

func loadBlackList() (err error) {
	err = initBlackRedis()
	if err != nil {
		beego.Error("init black redis error")
		return err
	}

	conn := secKillServer.BlackRedisPool.Get()
	defer conn.Close()

	// id black list
	idList, errDo := redis.Strings(conn.Do("hgetall", "idblacklist"))
	if errDo != nil {
		beego.Warn(fmt.Sprintf("load id black list from redis fail: err is %v", errDo))
		return errDo
	}
	for _, idStr := range idList {
		id, errStrconv := strconv.Atoi(idStr)
		if errStrconv != nil {
			beego.Warn("strconv the black id read from redis")
			continue
		}
		secKillServer.IdBlackMap[id] = true
	}

	// ip black list
	ipList, errDo := redis.Strings(conn.Do("hgetall", "ipblacklist"))
	if errDo != nil {
		beego.Warn(fmt.Sprintf("load ip black list from redis fail: err is %v", errDo))
		return errDo
	}
	for _, ip := range ipList {
		secKillServer.IpBlackMap[ip] = true
	}
	//go SyncIpBlackList()
	//go SyncIdBlackList()
	return nil
}

func initBlackRedis() (err error) {
	secKillServer.BlackRedisPool = &redis.Pool{
		MaxIdle:     secKillServer.RedisBlackConfig.RedisMaxIdle,
		MaxActive:   secKillServer.RedisBlackConfig.RedisMaxActive,
		IdleTimeout: time.Duration(secKillServer.RedisBlackConfig.RedisIdleTimeout) * time.Second,
		Dial: func() (conn redis.Conn, e error) {
			return redis.Dial("tcp", secKillServer.RedisBlackConfig.RedisAddr)
		},
	}

	conn := secKillServer.BlackRedisPool.Get()
	defer conn.Close()
	_, err = conn.Do("ping")
	if err != nil {
		return err
	}
	return nil
}

func initRedisProcessFunc() {
	for i := 0; i < secKillServer.WriteProxy2LayerGoroutineNum; i++ {
		go writeHandle()
	}

	for i := 0; i < secKillServer.ReadProxy2LayerGoroutineNum; i++ {
		go readHandle()
	}
}
