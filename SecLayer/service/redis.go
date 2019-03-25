package service

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego/logs"
	"github.com/garyburd/redigo/redis"
	"github.com/pkg/errors"
	"math/rand"
	"time"
)

func initRedis(conf *SecLayerConf) (err error) {
	seclayerContext.Proxy2LayerRedisPool, err = initRedisPool(conf.Proxy2LayerRedis)
	if err != nil {
		return errors.New(fmt.Sprintf("init proxy 2 layer redis pool error, err is %v", err))
	}

	seclayerContext.Layer2ProxyRedisPool, err = initRedisPool(conf.Layer2ProxyRedis)
	if err != nil {
		return errors.New(fmt.Sprintf("init layer 2 proxy redis pool error, err is %v", err))
	}
	return nil
}

func initRedisPool(redisConf RedisConfig) (pool *redis.Pool, err error) {
	pool = &redis.Pool{
		MaxIdle:     redisConf.RedisMaxIdle,
		MaxActive:   redisConf.RedisMaxActive,
		IdleTimeout: time.Duration(redisConf.RedisIdleTimeout) * time.Second,
		Dial: func() (conn redis.Conn, e error) {
			return redis.Dial("tcp", redisConf.RedisAddr)
		},
	}

	conn := pool.Get()
	defer conn.Close()
	_, err = conn.Do("ping")
	if err != nil {
		return nil, err
	}
	return pool, nil
}

func RunProcess() {
	for i := 0; i < seclayerContext.SecLayerConf.ReadProxy2LayerGoroutineNum; i++ {
		seclayerContext.waitGroup.Add(1)
		go HandleReader()
	}

	for i := 0; i < seclayerContext.SecLayerConf.WriteProxy2LayerGoroutineNum; i++ {
		seclayerContext.waitGroup.Add(1)
		go HandleWriter()
	}

	for i := 0; i < seclayerContext.SecLayerConf.HandleUserGoroutineNum; i++ {
		seclayerContext.waitGroup.Add(1)
		go HandleUser()
	}

	logs.Debug("all process goroutine started")
	seclayerContext.waitGroup.Wait()
	logs.Debug("wait all goroutine exited")
	return
}

func HandleReader() {
	logs.Debug("read goroutine running")
	for {
		conn := seclayerContext.Proxy2LayerRedisPool.Get()
		for {
			data, err := redis.String(conn.Do("blpop", seclayerContext.SecLayerConf.Proxy2LayerRedis.RedisQueueName, 0))
			if err != nil {
				conn.Close()
				break
			}
			logs.Debug("pop from queue, data:%s", data)
			var req SecRequest
			err = json.Unmarshal([]byte(data), &req)
			if err != nil {
				logs.Debug("handle reader josn unmarshal request error:%v", err)
				continue
			}
			now := time.Now().Unix()
			if now-req.AccessTime.Unix() >= seclayerContext.SecLayerConf.MaxRequestWaitTimeout {
				logs.Warn("req[%v] is expire", req)
				continue
			}
			timer := time.NewTicker(time.Millisecond * time.Duration(seclayerContext.SecLayerConf.SendToHandleChanTimeout))
			select {
			case seclayerContext.Read2HandleChan <- &req:

			case <-timer.C:
				logs.Warn("send to handle chan timeout, res:%v", req)
			}
		}
		conn.Close()
	}
}

func HandleWriter() {
	logs.Debug("handle write running")
	for res := range seclayerContext.Handle2WriteChan {
		err := sendToRedis(res)
		if err != nil {
			logs.Error("handle write send to redis err:%v, res:%v", err, res)
			continue
		}
	}
}
func sendToRedis(res *SecResponse) (err error) {
	data, err := json.Marshal(res)
	if err != nil {
		return err
	}
	conn := seclayerContext.Layer2ProxyRedisPool.Get()
	_, err = conn.Do("rpush", seclayerContext.SecLayerConf.Layer2ProxyRedis.RedisQueueName, string(data))
	if err != nil {
		return err
	}
	return nil
}

func HandleUser() {
	logs.Debug("handle user running")
	for req := range seclayerContext.Read2HandleChan {
		logs.Debug("begin process request:%v", req)
		res, err := HandleSecKill(req)
		if err != nil {
			logs.Warn("process request %v failed, err:%v", req, err)
			res = &SecResponse{
				Code: ErrServiceBusy,
			}
		}

		timer := time.NewTicker(time.Millisecond * time.Duration(seclayerContext.SecLayerConf.SendToWriteChanTimeout))
		select {
		case seclayerContext.Handle2WriteChan <- res:

		case <-timer.C:
			logs.Warn("send to response chan timeout, res:%v", res)
		}
	}
}

func HandleSecKill(req *SecRequest) (res *SecResponse, err error) {
	seclayerContext.RWSecProductLock.Lock()
	defer seclayerContext.RWSecProductLock.Unlock()

	res = &SecResponse{}
	product, ok := seclayerContext.SecLayerConf.SecProductInfoMap[req.ProductId]
	if !ok {
		logs.Error("not found product :%v", req.ProductId)
		res.Code = ErrNotFoundProduct
		return res, errors.New(fmt.Sprintf("not found product :%v", req.ProductId))
	}

	if product.Status == ProductStatusSoldout {
		res.Code = ErrSoldout
		return res, errors.New(fmt.Sprintf("product soldout :%v", req.ProductId))
	}

	now := time.Now().Unix()
	alreadySoldCount := product.secLimit.Check(now)
	if alreadySoldCount > product.SoldMaxLimit {
		res.Code = ErrRetry
		return
	}

	seclayerContext.HistoryMapLock.Lock()
	userHistory, ok := seclayerContext.HistoryMap[req.ProductId]
	if !ok {
		userHistory = &UserBuyHistory{
			history: make(map[int]int, 16),
		}

		seclayerContext.HistoryMap[req.ProductId] = userHistory
	}

	historyCount := userHistory.GetProductByCount(req.ProductId)
	seclayerContext.HistoryMapLock.Unlock()
	if historyCount > product.OnePersonBuyLimit {
		res.Code = ErrAlreadyBuy
		return
	}

	curSoldCount := seclayerContext.ProductCountMgr.Count(req.ProductId)
	if curSoldCount >= product.Total {
		res.Code = ErrSoldout
		product.Status = ProductStatusSoldout
		return
	}

	curRate := rand.Float64()
	if curRate > curRate {
		res.Code = ErrRetry
		return
	}

	userHistory.Add(req.ProductId, 1)
	seclayerContext.ProductCountMgr.Add(req.ProductId, 1)

	// //用户id&商品id&当前时间&密钥
	res.Code = ErrSecKillSucc
	tokenData := fmt.Sprintf("userId=%d&productId=%d&timestamp=%d&security=%s",
		req.UserId, req.ProductId, now, seclayerContext.SecLayerConf.TokenPasswd)

	res.Token = fmt.Sprintf("%x", md5.Sum([]byte(tokenData)))
	res.TokenTime = now
	return res, err
}
