package service

import (
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego"
)

func writeHandle() {

	for {
		req := <-secKillServer.SecReqChan

		conn := secKillServer.Proxy2LayerRedisPool.Get()
		defer conn.Close()

		data, err := json.Marshal(req)
		if err != nil {
			beego.Error(fmt.Sprintf("Proxy2Layer writeHandle json marshal error, error is %v req is %v", err, req))
			continue
		}
		_, err = conn.Do("LPUSH", data)
		if err != nil {
			beego.Error(fmt.Sprintf("Proxy2Layer writeHandle LPUSH data error, erros is %v data is %s", err, data))
			continue
		}
	}
}

func readHandle() {

}
