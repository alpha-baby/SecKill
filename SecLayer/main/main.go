package main

import (
	"SecKill/SecLayer/service"
	"fmt"
	"github.com/astaxie/beego/logs"
)

func main() {
	// 加载配置文件
	confType := "ini"
	confFilePath := "/Users/alphababy/go/src/SecKill/SecLayer/conf/secLayer.conf"
	err := initConfig(confType, confFilePath)
	if err != nil {
		panic(fmt.Sprintf("init config error , err: %v", err))
		return
	}

	// 初始化日志库
	err = initLogger()
	if err != nil {
		panic("Init log error")
		return
	}
	logs.Debug("Init log success")

	// 初始化秒杀逻辑
	err = service.InitService(appConfig)
	if err != nil {
		logs.Error(fmt.Sprintf("init seckill error, err is %v", err))
		panic(fmt.Sprintf("init seckill error, err is %v", err))
		return
	}
	logs.Debug("Init service.InitService success")

	// 启动服务
	logs.Debug("service SecLayer Run")
	service.Run()

	// 初始化etcd
	//ctx, cancelFunc := context.WithCancel(context.Background())
	//err = InitEtcd(appConfig.EtcdAddr, ctx)
	//if err != nil {
	//	log.Println("Init Etcd error,", err)
	//	panic("loadConf error")
	//}
	//defer cancelFunc()
	//logs.Debug("Init Etcd success ")
}
