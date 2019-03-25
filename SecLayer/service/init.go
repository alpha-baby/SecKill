package service

import "github.com/astaxie/beego/logs"

func InitService(conf *SecLayerConf) (err error) {
	secLayerConf = conf
	seclayerContext.SecLayerConf = conf
	seclayerContext.Read2HandleChan = make(chan *SecRequest, conf.Read2HandleChanSize)
	seclayerContext.Handle2WriteChan = make(chan *SecResponse, conf.Handle2WriteChanSize)
	seclayerContext.HistoryMap = make(map[int]*UserBuyHistory, 100000)
	seclayerContext.ProductCountMgr = NewProductCountMgr()
	// 初始化redis
	err = initRedis(conf)
	if err != nil {
		return err
	}
	logs.Debug("init redis success")
	// 初始化 etcd
	err = initEtcd(conf)
	if err != nil {
		return err
	}

	logs.Debug("init etcd success")
	err = LoadProductFromEtcd(conf)
	if err != nil {
		return err
	}
	logs.Debug("Load Product From Etcd success")
	return nil
}
