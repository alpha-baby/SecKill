package main

import (
	"SecKill/SecProxy/service"
	"github.com/astaxie/beego"
	"github.com/garyburd/redigo/redis"
	"go.etcd.io/etcd/clientv3"
)

var (
	RedisPoll  *redis.Pool      // redis 连接池
	EtcdClient *clientv3.Client // etcd 连接
)

func initSec() (err error) {
	// 初始化日志
	err = initLog()
	if err != nil {
		return err
	}

	// 初始化redis
	err = initRedis()
	if err != nil {
		return err
	}

	// 初始化etcd
	err = initEtcd()
	if err != nil {
		return err
	}

	// 加载秒杀配置
	err = loadSecConf()
	if err != nil {
		return err
	}

	err = service.InitService(SecKillConfig)
	if err != nil {
		return err
	}
	initSecProduct()

	beego.Debug("init sec success")
	return nil
}
