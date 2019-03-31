package service

import (
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/garyburd/redigo/redis"
	"time"
)

func writeHandle() {

	for {
		req := <-secKillServer.SecReqChan

		conn := secKillServer.Proxy2LayerRedisPool.Get()

		data, err := json.Marshal(req)
		if err != nil {
			beego.Error(fmt.Sprintf("Proxy2Layer writeHandle json marshal error, error is %v req is %v", err, req))
			conn.Close()
			continue
		}
		_, err = conn.Do("LPUSH", "sec_queue", string(data))
		if err != nil {
			beego.Error(fmt.Sprintf("Proxy2Layer writeHandle LPUSH data error, erros is %v data is %s", err, data))
			conn.Close()
			continue
		}

		conn.Close()
	}

}

func readHandle() {
	for {

		conn := secKillServer.Proxy2LayerRedisPool.Get()

		reply, err := conn.Do("RPOP", "recv_queue")
		data, err := redis.String(reply, err)

		if err == redis.ErrNil {
			time.Sleep(time.Second)
			conn.Close()
			continue
		}
		logs.Debug("rpop from reids success, data is %v", data)
		if err != nil {
			logs.Error("rpop failed, err:%v", err)
			conn.Close()
			continue
		}

		var result SecResult
		err = json.Unmarshal([]byte(data), &result)
		if err != nil {
			logs.Error("json.Unmarshal failed, err:%v", err)
			conn.Close()
			continue
		}

		userKey := fmt.Sprintf("%d_%d", result.UserId, result.ID)

		secKillServer.UserConnMapLock.Lock()
		resultChan, ok := secKillServer.UserConnMap[userKey]
		secKillServer.UserConnMapLock.Unlock()
		if !ok {
			conn.Close()
			logs.Warn("user not found:%v", userKey)
			continue
		}

		resultChan <- &result
		conn.Close()
	}
}
