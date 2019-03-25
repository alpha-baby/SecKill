package main

import (
	"fmt"
	"strings"

	"SecKill/SecProxy/service"
	"github.com/astaxie/beego"
	"github.com/pkg/errors"
)

var (
	SecKillConfig = &service.SecKillServer{
		SecProductInfoConfMap: make(map[int]*service.SecProductInfoConf, 1024),
	}
)

func initConfig() (err error) {

	//  读取redis相关的配置
	err = readRedisConfig()
	if err != nil {
		return err
	}

	// 读取etcd 配置信息
	err = readEtcdConfig()
	if err != nil {
		return err
	}

	// 读取日志配置
	SecKillConfig.LogPath = beego.AppConfig.String("log_path")
	if len(SecKillConfig.LogPath) == 0 {
		return errors.New("init config log_path err")
	}

	SecKillConfig.LogLevel = beego.AppConfig.String("log_level")
	if len(SecKillConfig.LogPath) == 0 {
		return errors.New("init config log_level err")
	}

	// 读取cookie 密钥
	SecKillConfig.CookieSecretKey = beego.AppConfig.String("cookie_secretkey")
	if len(SecKillConfig.CookieSecretKey) == 0 {
		return errors.New("init config cookie_secretkey err")
	}

	// 频率控制阈值
	SecKillConfig.UserSecAccessLimit, err = beego.AppConfig.Int("user_sec_access_limit")
	if err != nil {
		return errors.New("init config user_sec_access_limit err")
	}

	SecKillConfig.ReferWhiteList = strings.Split(beego.AppConfig.String("refer_whitelist"), ",")
	if len(SecKillConfig.ReferWhiteList) == 0 {
		return errors.New("init config refer_whitelist err")
	}

	SecKillConfig.IpSecAccessLimit, err = beego.AppConfig.Int("ip_sec_access_limit")
	if err != nil {
		return errors.New("init config ip_sec_access_limit err")
	}

	SecKillConfig.WriteProxy2LayerGoroutineNum, err = beego.AppConfig.Int("write_layer2proxy_goroutine_num")
	if err != nil {
		return errors.New("init config write_layer2proxy_goroutine_num err")
	}

	SecKillConfig.ReadProxy2LayerGoroutineNum, err = beego.AppConfig.Int("read_layer2proxy_goroutine_num")
	if err != nil {
		return errors.New("init config read_layer2proxy_goroutine_num err")
	}

	SecKillConfig.SecReqChanSize, err = beego.AppConfig.Int("sec_req_chan_size")
	if err != nil {
		return errors.New("init config sec_req_chan_size err")
	}

	return nil
}

// 读取redis 配置
func readRedisConfig() (err error) {
	SecKillConfig.WriteProxy2LayerGoroutineNum, err = beego.AppConfig.Int("write_proxy2layer_goroutine_num")
	if err != nil {
		return errors.New(fmt.Sprintf("init config write_proxy2layer_goroutine_num err is %v", err))
	}

	SecKillConfig.ReadProxy2LayerGoroutineNum, err = beego.AppConfig.Int("read_layer2proxy_goroutine_num")
	if err != nil {
		return errors.New(fmt.Sprintf("init config read_layer2proxy_goroutine_num err is %v", err))
	}
	// ;redis黑名单相关配置
	{
		SecKillConfig.RedisBlackConfig.RedisAddr = beego.AppConfig.String("redis_black_addr")
		if len(SecKillConfig.RedisBlackConfig.RedisAddr) == 0 {
			return errors.New("init config redis_black_addr err")
		}

		SecKillConfig.RedisBlackConfig.RedisMaxIdle, err = beego.AppConfig.Int("redis_black_idle")
		if err != nil {
			return errors.New(fmt.Sprintf("init config redis_black_idle err is %v", err))
		}

		SecKillConfig.RedisBlackConfig.RedisMaxActive, err = beego.AppConfig.Int("redis_black_active")
		if err != nil {
			return errors.New(fmt.Sprintf("init config redis_black_active err is %v", err))
		}

		SecKillConfig.RedisBlackConfig.RedisIdleTimeout, err = beego.AppConfig.Int("redis_black_idle_timeout")
		if err != nil {
			return errors.New(fmt.Sprintf("init config redis_black_idle_timeout err is %v", err))
		}
	}
	// ;redis 接入层->业务逻辑层
	{
		SecKillConfig.RedisProxy2LayerConf.RedisAddr = beego.AppConfig.String("redis_proxy2layer_addr")
		if len(SecKillConfig.RedisProxy2LayerConf.RedisAddr) == 0 {
			return errors.New("init config redis_proxy2layer_addr err")
		}

		SecKillConfig.RedisProxy2LayerConf.RedisMaxIdle, err = beego.AppConfig.Int("redis_proxy2layer_idle")
		if err != nil {
			return errors.New(fmt.Sprintf("init config redis_proxy2layer_idle err is %v", err))
		}

		SecKillConfig.RedisProxy2LayerConf.RedisMaxActive, err = beego.AppConfig.Int("redis_proxy2layer_active")
		if err != nil {
			return errors.New(fmt.Sprintf("init config redis_proxy2layer_active err is %v", err))
		}

		SecKillConfig.RedisProxy2LayerConf.RedisIdleTimeout, err = beego.AppConfig.Int("redis_proxy2layer_idle_timeout")
		if err != nil {
			return errors.New(fmt.Sprintf("init config redis_proxy2layer_idle_timeout  err is %v", err))
		}
	}

	return nil
}

// 读取etcd相关的配置
func readEtcdConfig() (err error) {

	SecKillConfig.EtcdConfig.EtcdAddr = beego.AppConfig.String("etcd_addr")
	if len(SecKillConfig.EtcdConfig.EtcdAddr) == 0 {
		return errors.New("init config etcd_addr err")
	}

	SecKillConfig.EtcdConfig.EtcdTimeout, err = beego.AppConfig.Int("etcd_timeout")
	if err != nil {
		return errors.New("init config etcd_timeout err")
	}

	SecKillConfig.EtcdConfig.EtcdSecKeyPrefix = beego.AppConfig.String("etcd_sec_key_prefix")
	if len(SecKillConfig.EtcdConfig.EtcdSecKeyPrefix) == 0 {
		return errors.New("init config etcd_sec_key_prefix err")
	} else {

		if strings.HasSuffix(SecKillConfig.EtcdConfig.EtcdSecKeyPrefix, "/") == false {
			SecKillConfig.EtcdConfig.EtcdSecKeyPrefix += "/"
		}
	}

	SecKillConfig.EtcdConfig.EtcdProductKey = beego.AppConfig.String("etcd_product_key")
	if len(SecKillConfig.EtcdConfig.EtcdProductKey) == 0 {
		return errors.New("init config etcd_product_key err")
	} else {
		// 把前缀拼接上
		SecKillConfig.EtcdConfig.EtcdProductKey = fmt.Sprintf("%s%s", SecKillConfig.EtcdConfig.EtcdSecKeyPrefix, SecKillConfig.EtcdConfig.EtcdProductKey)
	}

	SecKillConfig.EtcdConfig.EtcdBlackListKey = beego.AppConfig.String("etcd_black_list_key")
	if len(SecKillConfig.EtcdConfig.EtcdBlackListKey) == 0 {
		return errors.New("init config etcd_black_list_key err")
	} else {
		// 把前缀拼接上
		SecKillConfig.EtcdConfig.EtcdBlackListKey = fmt.Sprintf("%s%s", SecKillConfig.EtcdConfig.EtcdSecKeyPrefix, SecKillConfig.EtcdConfig.EtcdBlackListKey)
	}

	return nil
}
