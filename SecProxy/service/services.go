package service

import (
	"crypto/md5"
	"errors"
	"fmt"
	"github.com/astaxie/beego/logs"
	"time"

	"github.com/astaxie/beego"
)

func SecKill(req *SecRequest) (data map[string]interface{}, code int, err error) {
	secKillServer.RWSecProductLock.Lock()
	defer secKillServer.RWSecProductLock.Unlock()

	// 检查请求身份是否合法
	err = userCheck(req)
	if err != nil {
		code = ErrUserCheckAuthFailed
		beego.Warn(fmt.Sprintf("userID:%d identification invalid, Check failed , req[%v]", req.UserId, req))
		return
	}

	// 检查请求频率是否合法  过滤机器人的高并发请求
	err = antiSpam(req)
	if err != nil {
		code = ErrUserServiceBusy
		beego.Warn(fmt.Sprintf("UserId[%d] request's frequency invalid, req[%v]", req.UserId, req))
		err = errors.New("server is busy")
		return
	}

	data, code, err = SecInfoById(req.ID)
	if err != nil {
		beego.Warn(fmt.Sprintf("UserId[%d] SecInfo By Id fail , req[%v]", req.UserId, req))
		return
	}

	if code != 0 {
		beego.Warn(fmt.Sprintf("UserId[%d] SecInfo By Id fail ,code[%d] req[%v]", req.UserId, code, req))
		return
	}

	userKey := fmt.Sprintf("%d_%d", req.UserId, req.ID)
	logs.Debug(userKey)
	secKillServer.UserConnMap[userKey] = req.ResultChan

	secKillServer.SecReqChan <- req

	ticker := time.NewTicker(time.Second * 5)

	defer func() {
		ticker.Stop()
		secKillServer.UserConnMapLock.Lock()
		delete(secKillServer.UserConnMap, userKey)
		secKillServer.UserConnMapLock.Unlock()
	}()

	select {
	case <-ticker.C:
		code = ErrProcessTimeout
		err = fmt.Errorf("request timeout")

		return
	case <-req.CloseNotify:
		code = ErrClientClosed
		err = fmt.Errorf("client already closed")
		return
	case result := <-req.ResultChan:
		code = result.Code
		data["id"] = result.ID
		data["product_id"] = result.ProductId
		data["token"] = result.Token
		data["user_id"] = result.UserId

		return
	}
	return
}

func SecInfoList() (data []map[string]interface{}, code int, err error) {

	secKillServer.RWSecProductLock.RLock()
	defer secKillServer.RWSecProductLock.RUnlock()

	for _, v := range secKillServer.SecProductInfoConfMap {

		item, _, err := SecInfoById(v.ID)
		if err != nil {
			beego.Debug(fmt.Sprintf("get product_id[%d] failed, err:%v", v.ProductId, err))
			continue
		}
		data = append(data, item)
	}

	return
}

func SecInfo(Id int64) (data []map[string]interface{}, code int, err error) {

	secKillServer.RWSecProductLock.RLock()
	defer secKillServer.RWSecProductLock.RUnlock()

	item, code, err := SecInfoById(Id)
	if err != nil {
		return
	}

	data = append(data, item)
	return
}

func SecInfoById(Id int64) (data map[string]interface{}, code int, err error) {

	v, ok := secKillServer.SecProductInfoConfMap[Id]
	if !ok {
		code = ErrNotFoundProductId
		err = fmt.Errorf("not found product_id:%d", Id)
		return
	}

	start := false
	end := false
	status := "success"

	now := time.Now().Unix()
	if now-v.StartTime < 0 {
		start = false
		end = false
		status = "sec kill is not start"
		code = ErrActiveNotStart
	}

	if now-v.StartTime > 0 {
		start = true
	}

	if now-v.EndTime > 0 {
		start = false
		end = true
		status = "sec kill is already end"
		code = ErrActiveAlreadyEnd
	}

	if v.Status == ProductStatusForceSaleOut || v.Status == ProductStatusSaleOut {
		start = false
		end = true
		status = "product is sale out"
		code = ErrActiveSaleOut
	}

	data = make(map[string]interface{})

	data["id"] = v.ID
	data["product_id"] = v.ProductId
	data["start"] = start
	data["end"] = end
	data["status"] = status

	return
}

func NewSecRequest() (secRequest *SecRequest) {
	secRequest = &SecRequest{
		ResultChan: make(chan *SecResult, 1),
	}

	return
}

func userCheck(req *SecRequest) (err error) {
	//found := false
	//// 检查用户的请求refer 是否在白名单中 如果没在会拒绝请求
	//logs.Debug(fmt.Sprintf("client refence is %v", req.ClientRefence))
	//for _, refer := range secKillServer.ReferWhiteList {
	//	if refer == req.ClientRefence {
	//		found = true
	//		break
	//	}
	//}
	//
	//if !found {
	//	err = fmt.Errorf("invalid request ")
	//	beego.Warn(fmt.Sprintf("user [%d] is reject by refer , request :%v", req.UserId, req))
	//	return
	//}

	// 检车用户身份是否合法
	authData := fmt.Sprintf("%d:%s", req.UserId, secKillServer.CookieSecretKey)
	authSign := fmt.Sprintf("%x", md5.Sum([]byte(authData)))

	if authSign != req.UserAuthSign {
		beego.Debug(fmt.Sprintf("authSign is %s, req.UserAuthSign is %s", authSign, req.UserAuthSign))
		err = fmt.Errorf("invalid user cookie auth")
		return err
	}
	return nil
}
