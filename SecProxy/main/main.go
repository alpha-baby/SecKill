package main

import (
	_ "SecKill/SecProxy/router"
	"fmt"
	"github.com/astaxie/beego"
)

func main() {
	// 初始化配置信息
	err := initConfig()
	if err != nil {
		panic(fmt.Sprintf("init config error:%v", err))
	}

	err = initSec()
	if err != nil {
		panic(fmt.Sprintf("init sec error:%v", err))
	}

	beego.Run()
}
