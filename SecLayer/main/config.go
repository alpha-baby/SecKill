package main

import (
	"SecKill/SecLayer/service"
	"github.com/astaxie/beego/config"
	"github.com/pkg/errors"
)

var (
	appConfig *service.SecLayerConf
	conf      config.Configer
)

func initConfig(confType, fileName string) (err error) {
	conf, err = config.NewConfig(confType, fileName)
	if err != nil {
		return
	}

	appConfig = &service.SecLayerConf{}
	// logs 读取日志相关配置
	{
		//  logs::log_level
		appConfig.LogLevel = conf.String("logs::log_level")
		// 容错性处理
		if len(appConfig.LogLevel) == 0 {
			appConfig.LogLevel = "debug"
		}

		appConfig.LogPath = conf.String("logs::log_path")
		if len(appConfig.LogPath) == 0 {
			err = errors.New("init config logs::log_path error")
			return
		}
	}

	err = loadServiceConfig()
	if err != nil {
		return err
	}

	// 读取redis相关配置
	err = loadRedisConfig()
	if err != nil {
		return err
	}

	err = loadEtcdConfig()
	if err != nil {
		return err
	}

	return
}

func loadRedisConfig() (err error) {
	// proxy 2 layer
	{
		appConfig.Proxy2LayerRedis.RedisAddr = conf.String("redis::redis_proxy2layer_addr")
		if len(appConfig.Proxy2LayerRedis.RedisAddr) == 0 {
			return errors.New("init config redis::redis_proxy2layer_addr error")
		}

		appConfig.Proxy2LayerRedis.RedisMaxIdle, err = conf.Int("redis::redis_proxy2layer_idle")
		if err != nil {
			return errors.New("init config redis::redis_proxy2layer_idle error")
		}

		appConfig.Proxy2LayerRedis.RedisMaxActive, err = conf.Int("redis::redis_proxy2layer_active")
		if err != nil {
			return errors.New("init config redis::redis_proxy2layer_active error")
		}

		appConfig.Proxy2LayerRedis.RedisIdleTimeout, err = conf.Int("redis::redis_proxy2layer_idle_timeout")
		if err != nil {
			return errors.New("init config redis::redis_proxy2layer_idle_timeout error")
		}

		appConfig.Proxy2LayerRedis.RedisQueueName = conf.String("redis::redis_proxy2layer_queue_name")
		if len(appConfig.Proxy2LayerRedis.RedisAddr) == 0 {
			return errors.New("init config redis::redis_proxy2layer_queue_name error")
		}
	}
	// layer 2 proxy
	{
		appConfig.Layer2ProxyRedis.RedisAddr = conf.String("redis::redis_layer2proxy_addr")
		if len(appConfig.Proxy2LayerRedis.RedisAddr) == 0 {
			return errors.New("init config redis::redis_layer2proxy_addr error")
		}

		appConfig.Layer2ProxyRedis.RedisMaxIdle, err = conf.Int("redis::redis_layer2proxy_idle")
		if err != nil {
			return errors.New("init config redis::redis_layer2proxy_idle error")
		}

		appConfig.Layer2ProxyRedis.RedisMaxActive, err = conf.Int("redis::redis_layer2proxy_active")
		if err != nil {
			return errors.New("init config redis::redis_layer2proxy_active error")
		}

		appConfig.Layer2ProxyRedis.RedisIdleTimeout, err = conf.Int("redis::redis_layer2proxy_idle_timeout")
		if err != nil {
			return errors.New("init config redis::redis_layer2proxy_idle_timeout error")
		}

		appConfig.Layer2ProxyRedis.RedisQueueName = conf.String("redis::redis_layer2proxy_queue_name")
		if len(appConfig.Proxy2LayerRedis.RedisAddr) == 0 {
			return errors.New("init config redis::redis_layer2proxy_queue_name error")
		}
	}
	return nil
}

func loadServiceConfig() (err error) {
	appConfig.ReadProxy2LayerGoroutineNum, err = conf.Int("service::read_layer2proxy_goroutine_num")
	if err != nil {
		return errors.New("init config service::read_layer2proxy_goroutine_num error")
	}

	appConfig.WriteProxy2LayerGoroutineNum, err = conf.Int("service::write_layer2proxy_goroutine_num")
	if err != nil {
		return errors.New("init config service::write_layer2proxy_goroutine_num error")
	}

	appConfig.HandleUserGoroutineNum, err = conf.Int("service::handle_user_goroutine_num")
	if err != nil {
		return errors.New("init config service::handle_user_goroutine_num error")
	}

	appConfig.Read2HandleChanSize, err = conf.Int("service::read2handle_chan_size")
	if err != nil {
		return errors.New("init config service::read2handle_chan_size error")
	}

	appConfig.Handle2WriteChanSize, err = conf.Int("service::handle2write_chan_size")
	if err != nil {
		return errors.New("init config service::handle2write_chan_size error")
	}

	appConfig.MaxRequestWaitTimeout, err = conf.Int64("service::max_request_wait_timeout")
	if err != nil {
		return errors.New("init config service::max_request_wait_timeout error")
	}

	appConfig.SendToWriteChanTimeout, err = conf.Int("service::send_to_write_chan_timeout")
	if err != nil {
		return errors.New("init config service::send_to_write_chan_timeout error")
	}

	appConfig.SendToHandleChanTimeout, err = conf.Int("service::send_to_handle_chan_timeout")
	if err != nil {
		return errors.New("init config service::send_to_handle_chan_timeout error")
	}

	appConfig.TokenPasswd = conf.String("service::seckill_token_passwd")
	if len(appConfig.TokenPasswd) == 0 {
		return errors.New("init config service::seckill_token_passwd error")
	}
	return nil
}

func loadEtcdConfig() (err error) {
	appConfig.EtcdConfig.EtcdAddr = conf.String("etcd::server_addr")
	if len(appConfig.TokenPasswd) == 0 {
		return errors.New("init config etcd::server_addr error")
	}

	appConfig.EtcdConfig.EtcdSecProductKey = conf.String("etcd::etcd_product_key")
	if len(appConfig.TokenPasswd) == 0 {
		return errors.New("init config etcd::etcd_product_key error")
	}

	appConfig.EtcdConfig.EtcdSecKeyPrefix = conf.String("etcd::etcd_sec_key_prefix")
	if len(appConfig.TokenPasswd) == 0 {
		return errors.New("init config etcd::etcd_sec_key_prefix error")
	}

	appConfig.EtcdConfig.Timeout, err = conf.Int("etcd::timeout")
	if err != nil {
		return errors.New("init config etcd::timeout error")
	}
	return nil
}
