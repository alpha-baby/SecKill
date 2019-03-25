package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/astaxie/beego/logs"
	"go.etcd.io/etcd/clientv3"
	"strings"
	"time"
)

func initEtcd(conf *SecLayerConf) (err error) {
	seclayerContext.EtcdClient, err = clientv3.New(clientv3.Config{
		Endpoints:   strings.Split(conf.EtcdConfig.EtcdAddr, ","),
		DialTimeout: time.Duration(conf.EtcdConfig.Timeout) * time.Second,
	})
	if err != nil {
		return errors.New(fmt.Sprintf("init service etcd error, err is %v", err))
	}

	ctx, cancelFunc := context.WithTimeout(context.Background(), time.Second*3)
	defer cancelFunc()
	res, err := seclayerContext.EtcdClient.Get(ctx, "version")
	if err != nil {
		return errors.New(fmt.Sprintf("init service etcd error, err is %v", err))
	}
	logs.Debug(fmt.Sprintf("test init etcd , version is %v", res.Kvs))
	return nil
}
