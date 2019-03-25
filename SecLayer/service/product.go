package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/gpmgo/gopm/modules/log"
	"go.etcd.io/etcd/mvcc/mvccpb"
	"time"
)

func LoadProductFromEtcd(conf *SecLayerConf) (err error) {
	ctx, cancelFunc := context.WithTimeout(context.Background(), time.Second*1)
	resp, err1 := seclayerContext.EtcdClient.Get(ctx, conf.EtcdConfig.EtcdSecProductKey)
	if err != nil {
		log.Debug("load product from etcd error")
		return err1
	}
	defer cancelFunc()
	var secProductInfo []SecProductInfoConf
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
	// 开启
	go watchSecProductKey(conf.EtcdConfig.EtcdSecProductKey)
	return nil
}

func updateSecProductInfo(secProductInfo []SecProductInfoConf) {
	if len(secProductInfo) == 0 {
		return
	}

	var tmp map[int]*SecProductInfoConf = make(map[int]*SecProductInfoConf, 1024)
	seclayerContext.RWSecProductLock.RLock()
	if len(secLayerConf.SecProductInfoMap) != 0 {
		tmp = secLayerConf.SecProductInfoMap
	}
	seclayerContext.RWSecProductLock.RUnlock()
	for _, v := range secProductInfo {
		//product := service.SecProductInfoConf{}
		product := v
		product.secLimit = &SecLimit{}
		tmp[v.ProductId] = &product
	}

	seclayerContext.RWSecProductLock.Lock()
	secLayerConf.SecProductInfoMap = tmp
	seclayerContext.RWSecProductLock.Unlock()

	beego.Debug("add Sec Product info Conf :", tmp)

}

func watchSecProductKey(key string) {

	beego.Debug(fmt.Sprintf("begin watch key:%s", key))
	for {
		rch := seclayerContext.EtcdClient.Watch(context.Background(), key)
		var secProductInfo []SecProductInfoConf
		var getConfSucc = true

		for wresp := range rch {
			for _, ev := range wresp.Events {
				if ev.Type == mvccpb.DELETE {
					beego.Warn(fmt.Sprintf("key[%s] 's config deleted", key))
					continue
				}

				if ev.Type == mvccpb.PUT && string(ev.Kv.Key) == key {
					err := json.Unmarshal(ev.Kv.Value, &secProductInfo)
					if err != nil {
						beego.Error(fmt.Sprintf("key [%s], Unmarshal[%s], err:%v ", err))
						getConfSucc = false
						continue
					}
				}
			}

			if getConfSucc {
				beego.Debug(fmt.Sprintf("get new config from etcd: %v", secProductInfo))
				updateSecProductInfo(secProductInfo)
			}
		}
	}
}