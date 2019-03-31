package model

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"time"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gpmgo/gopm/modules/log"
)

const (
	ActivityStatusNormal  = 0
	ActivityStatusDisable = 1
	ActivityStatusExpire  = 2
)

type Activity struct {
	ActivityId   int    `db:"id"`
	ActivityName string `db:"name"`
	ProductId    int    `db:"product_id"`
	StartTime    int64  `db:"start_time"`
	EndTime      int64  `db:"end_time"`
	Total        int    `db:"total"`
	Status       int    `db:"status"`

	StartTimeStr string
	EndTimeStr   string
	StatusStr    string
	Speed        int `db:"sec_speed"`
	BuyLimit     int `db:"buy_limit"`
	BuyRate     float64 `db:"buy_rate"`
}

type SecProductInfoConf struct {
	ID int64
	ProductId         int
	StartTime         int64
	EndTime           int64
	Status            int
	Total             int
	Left              int
	OnePersonBuyLimit int
	BuyRate           float64
	//每秒最多能卖多少个
	SoldMaxLimit int
}

type ActivityModel struct {
}

func NewActivityModel() *ActivityModel {
	return &ActivityModel{}
}

func (p *ActivityModel) GetActivityList() (activityList []*Activity, err error) {
	sql := "select id, name, product_id, start_time, end_time, total, status from activity order by id desc"
	err = Db.Select(&activityList, sql)
	if err != nil {
		logs.Error("select activity from database failed, err:%v", err)
		return
	}

	for _, v := range activityList {
		t := time.Unix(v.StartTime, 0)
		v.StartTimeStr = t.Format("2006-01-02 15:04:05")

		t = time.Unix(v.EndTime, 0)
		v.EndTimeStr = t.Format("2006-01-02 15:04:05")

		now := time.Now().Unix()

		if now > v.EndTime {
			v.StatusStr = "已结束"
			continue
		}

		if v.Status == ActivityStatusNormal {
			v.StatusStr = "正常"
		} else if v.Status == ActivityStatusDisable {
			v.StatusStr = "已禁用"
		}
	}
	logs.Debug("get activity succ,  activity list is[%v]", activityList)
	return
}

func (p *ActivityModel) ProductValid(productId int, total int) (exist bool, totalValid bool, err error) {
	sql := "select id, name, total, status from product where id = ?"
	var productList []*Product

	err = Db.Select(&productList, sql, productId)
	if err != nil {
		logs.Debug("select product failed , err:%v", err)
		return false, false, err
	}
	logs.Debug("product is %v", productList[0])
	if len(productList) == 0 {
		return false, true, nil
	} else {
		if productList[0].Total < total {
			return true, false, nil
		}
		return true, true, nil
	}
}

func (p *ActivityModel) CreateActivity(activity *Activity) (err error) {
	// 判断对应要创建活动的商品是否存在
	exists, totalValid, err := p.ProductValid(activity.ProductId, activity.Total)
	if err != nil {
		logs.Error("Product Valid mysql select err is %v", err)
		err = fmt.Errorf("商品id[%d]不存在，或这设置的数量%d过多", activity.ProductId, activity.Total)
		return err
	}
	if !exists {
		err = fmt.Errorf("商品id[%d]不存在", activity.ProductId)
		logs.Warn("product id[%d] is not exists", activity.ProductId)
		return err
	}
	if !totalValid {
		err = fmt.Errorf("设置的商品数量%d大于库存", activity.Total)
		logs.Warn("product id[%d] total more than stock ", activity.ProductId)
		return err
	}

	if activity.StartTime <= 0 || activity.EndTime <= 0 || activity.EndTime <= activity.StartTime || activity.StartTime <= time.Now().Unix() {
		err = fmt.Errorf("开始时间[%d] 或者结束时间[%d] 错误", activity.StartTime, activity.EndTime)
		logs.Warn("invalid start[%d] end[%d] time", activity.StartTime, activity.EndTime)
		return err
	}

	sql := "insert into activity(name, product_id, start_time, end_time, total, buy_limit, sec_speed, buy_rate)values(?,?,?,?,?,?,?,?)"
	_, err = Db.Exec(sql, activity.ActivityName, activity.ProductId,
		activity.StartTime, activity.EndTime, activity.Total, activity.BuyLimit, activity.Speed, activity.BuyRate)
	if err != nil {
		logs.Warn("select from mysql failed, err:%v sql:%v", err, sql)
		err = fmt.Errorf(" 请检查数据是否填写正确！")
		return
	}

	logs.Debug("insert into database succ")
	err = p.SyncToEtcd(activity)
	if err != nil {
		logs.Warn("sync to etcd failed, err:%v data:%v", err, activity)
		return nil
	}
	logs.Debug("product:{%v} sync to etcd key{%s} success",activity, EtcdProductKey)
	return nil
}

func (p *ActivityModel) SyncToEtcd(activity *Activity) (err error){
	secProductInfoList, err := LoadProductFromEtcd(EtcdProductKey)
	if err != nil {
		log.Warn("load product from etcd err : %v", err)
		return err
	}

	var secProductInfo SecProductInfoConf
	secProductInfo.ID = getRandomSecKillID()
	secProductInfo.EndTime =  activity.EndTime
	secProductInfo.OnePersonBuyLimit = activity.BuyLimit
	secProductInfo.ProductId = activity.ProductId
	secProductInfo.SoldMaxLimit = activity.Speed
	secProductInfo.StartTime = activity.StartTime
	secProductInfo.Status = activity.Status
	secProductInfo.Total = activity.Total
	secProductInfo.BuyRate = activity.BuyRate

	secProductInfoList = append(secProductInfoList, secProductInfo)

	data, err := json.Marshal(secProductInfoList)
	if err != nil {
		logs.Error("json marshal failed, err:%v", err)
		return
	}

	_, err = EtcdClient.Put(context.Background(), EtcdProductKey, string(data))
	if err != nil {
		logs.Error("put to etcd failed, err:%v, data[%v]", err, string(data))
		return
	}

	return nil
}

func LoadProductFromEtcd(key string) (secProductInfo []SecProductInfoConf, err error) {
	ctx, cancelFunc := context.WithTimeout(context.Background(), time.Second*1)
	resp, err1 := EtcdClient.Get(ctx, key)
	if err1 != nil {
		log.Warn("load product from etcd error : %v", err1)
		return nil, err1
	}
	defer cancelFunc()

	for k, v := range resp.Kvs {
		beego.Debug(fmt.Sprintf("key:[%v] value:[%v] ", k, v))
		err = json.Unmarshal(v.Value, &secProductInfo)
		if err != nil {
			return nil, errors.New(fmt.Sprintf("Unmarshal secProductInfo err: %v", err))
		}
		beego.Debug(fmt.Sprintf("load sec conf is %v", secProductInfo))
	}

	return secProductInfo, nil
}

func getRandomSecKillID() int64 {
	unixTimeStamp := time.Now().Unix()
	randomInt := rand.Int63n(1000000000)
	randomIDStr := 1000000000 * unixTimeStamp + randomInt
	return randomIDStr
}
