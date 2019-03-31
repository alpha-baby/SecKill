package main

import (
	"SecKill/SecProxy/service"
	"context"
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego"
	"go.etcd.io/etcd/mvcc/mvccpb"
	"strings"
	"time"

	"github.com/pkg/errors"
	"go.etcd.io/etcd/clientv3"
)

func initEtcd() (err error) {
	EtcdClient, err = clientv3.New(clientv3.Config{
		Endpoints:   strings.Split(SecKillConfig.EtcdConfig.EtcdAddr, ","),
		DialTimeout: time.Duration(SecKillConfig.EtcdConfig.EtcdTimeout) * time.Second,
	})
	if err != nil {
		return err
	}

	// 检测etcd是否连接成功
	ctx, cancelFunc := context.WithTimeout(context.Background(), 2*time.Second)
	_, err = EtcdClient.Get(ctx, "/logagent/confs/")
	defer cancelFunc()
	if err != nil {
		return errors.New(fmt.Sprintf("connect etcd server error, ", err))
	}

	return nil
}

func loadSecConf() (err error) {
	resp, err1 := EtcdClient.Get(context.Background(), SecKillConfig.EtcdConfig.EtcdProductKey)
	if err != nil {
		return err1
	}
	var secProductInfo []service.SecProductInfoConf
	for k, v := range resp.Kvs {
		beego.Debug(fmt.Sprintf("key:[%v] value:[%v] ", k, v))
		err = json.Unmarshal(v.Value, &secProductInfo)
		if err != nil {
			return errors.New("Unmarshal secProductInfo err")
		}
		beego.Debug(fmt.Sprintf("load sec conf is %v", secProductInfo))
	}

	// 把从etcd中读取到的秒杀商品配置信息存放到系统全局变量中去
	go updateSecProductInfo(secProductInfo)
	// 开启go协程监听etcd是否有改动或者增加
	go watchSecProductKey(SecKillConfig.EtcdConfig.EtcdProductKey)
	return nil
}

func updateSecProductInfo(secProductInfo []service.SecProductInfoConf) {
	var tmp map[int64]*service.SecProductInfoConf = make(map[int64]*service.SecProductInfoConf, 1024)
	for _, v := range secProductInfo {
		produtInfo := v
		tmp[v.ID] = &produtInfo
	}

	SecKillConfig.RWSecProductLock.Lock()
	SecKillConfig.SecProductInfoConfMap = tmp
	SecKillConfig.RWSecProductLock.Unlock()

	beego.Debug("add Sec Product info Conf :", tmp)

}

func watchSecProductKey(key string) {

	beego.Debug(fmt.Sprintf("begin watch key:%s", key))
	for {
		rch := EtcdClient.Watch(context.Background(), key)
		var secProductInfo []service.SecProductInfoConf

		for wresp := range rch {
			for _, ev := range wresp.Events {
				if ev.Type == mvccpb.DELETE {
					beego.Warn(fmt.Sprintf("key[%s] 's config deleted", key))
					secProductInfo = []service.SecProductInfoConf{}
					updateSecProductInfo(secProductInfo)
					continue
				}

				if ev.Type == mvccpb.PUT && string(ev.Kv.Key) == key {
					err := json.Unmarshal(ev.Kv.Value, &secProductInfo)
					if err != nil {
						beego.Error(fmt.Sprintf("key [%s], Unmarshal[%s], err:%v ", err))
						continue
					}

					beego.Debug(fmt.Sprintf("get new config from etcd: %v", secProductInfo))
					updateSecProductInfo(secProductInfo)
				}
			}
		}
	}
}
