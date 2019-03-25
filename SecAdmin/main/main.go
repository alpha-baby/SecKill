package main

import (
	_ "SecKill/SecAdmin/router"
	"fmt"
	"github.com/astaxie/beego"
)

func main() {
	err := Init()
	if err != nil {
		panic(fmt.Sprintf("init err, err is %v", err))
		return
	}
	beego.Run()
}
