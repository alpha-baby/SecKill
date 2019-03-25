package service

import (
	"fmt"
	"github.com/astaxie/beego"
	"sync"
)

type SecLimitMgr struct {
	UserLimitMap map[int]*SecLimit
	IpLimitMap   map[string]*SecLimit
	lock         sync.Mutex
}

func antiSpam(req *SecRequest) (err error) {

	// 用户频率检测
	secKillServer.secLimitMgr.lock.Lock()
	secLimit, ok := secKillServer.secLimitMgr.UserLimitMap[req.UserId]
	if !ok {
		secLimit = &SecLimit{}
		secKillServer.secLimitMgr.UserLimitMap[req.UserId] = secLimit
	}
	count := secLimit.Count(req.AccessTime.Unix())
	secKillServer.secLimitMgr.lock.Unlock()
	beego.Warn(fmt.Sprintf("userId[%d] count is [%d]", req.UserId, count))
	if count > secKillServer.UserSecAccessLimit {
		beego.Warn(fmt.Sprintf("count is %d, secKillServer.UserSecAccessLimit is %d", count, secKillServer.UserSecAccessLimit))
		err = fmt.Errorf("UserId[%d] invalid request's frequency", req.UserId)
		return
	}

	// ip 频率检测
	secKillServer.secLimitMgr.lock.Lock()
	ipLimit, ok := secKillServer.secLimitMgr.IpLimitMap[req.ClientAddr]
	if !ok {
		ipLimit = &SecLimit{count: 1}
		secKillServer.secLimitMgr.IpLimitMap[req.ClientAddr] = ipLimit
	}

	count = ipLimit.Count(req.AccessTime.Unix())
	secKillServer.secLimitMgr.lock.Unlock()

	if count > secKillServer.UserSecAccessLimit {
		err = fmt.Errorf("IP:[%s] invalid request's frequency", req.ClientAddr)
		return
	}
	return nil
}
