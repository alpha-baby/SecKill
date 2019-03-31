package activity

import (
	"SecKill/SecAdmin/model"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"net/http"
	"strings"
)

type ActivityController struct {
	beego.Controller
}

func (p *ActivityController) CreateActivity() {
	p.TplName = "activity/create.html"
	p.Layout = "layout/layout.html"
	return
}

func (p *ActivityController) ListActivity() {

	activityModel := model.NewActivityModel()
	activityList, err := activityModel.GetActivityList()
	if err != nil {
		logs.Warn("get activity list failed, err :%v", err)
		return
	}

	p.Data["activity_list"] = activityList
	p.TplName = "activity/list.html"
	p.Layout = "layout/layout.html"
	return
}

func (p *ActivityController) SubmitActivity() {

	activityModel := model.NewActivityModel()
	var activity model.Activity
	var err error
	var errMsg string

	defer func() {
		p.Layout = "layout/layout.html"

		if err != nil {
			p.TplName = "activity/error.html"
			p.Data["Error"] = errMsg
		}
	}()

	name := strings.TrimSpace(p.GetString("activity_name"))
	if len(name) == 0 {
		errMsg = "活动名称不能为空"
		logs.Warn("submit activity the activity name err")
		err = fmt.Errorf("submit activity the product name is err")
		return
	}
	productId, err := p.GetInt("product_id")
	if err != nil {
		errMsg = "商品的Id错误"
		logs.Warn("submit activity the product id err")
		err = fmt.Errorf("submit activity the product id err")
		return
	}
	startTime, err := p.GetInt64("start_time")
	if (err != nil) {
		err = fmt.Errorf("开始时间 非法, err:%v", err)
		errMsg = "开始时间 不合要求"
		return
	}

	endTime, err := p.GetInt64("end_time")
	if (err != nil) {
		err = fmt.Errorf("结束时间 非法, err:%v", err)
		errMsg = "结束时间 不合要求"
		return
	}

	total, err := p.GetInt("total")
	if (err != nil) {
		err = fmt.Errorf("商品数量 非法, err:%v", err)
		errMsg = "商品数量错误"
		return
	}

	speed, err := p.GetInt("speed")
	if (err != nil) || speed < 1 {
		err = fmt.Errorf("商品抢购限制速度 非法, err:%v", err)
		errMsg = "商品抢购限制速度 错误"
		return
	}

	limit, err := p.GetInt("buy_limit")
	if (err != nil) {
		err = fmt.Errorf("商品个人限购数量 非法, err:%v", err)
		errMsg = "商品个人限购数量 错误"
		return
	}

	buyRate, err := p.GetFloat("buy_rate")
	if (err != nil) || (buyRate <=0 || buyRate > 1){
		err = fmt.Errorf("商品个人抢购概率 非法, err:%v", err)
		errMsg = "商品个人抢购概率 错误 请重新输入"
		return
	}
	activity.ActivityName = name
	activity.ProductId = productId
	activity.StartTime = startTime
	activity.EndTime = endTime
	activity.Total = total
	activity.Speed = speed
	activity.BuyLimit = limit
	activity.BuyRate = buyRate

	err = activityModel.CreateActivity(&activity)
	if err != nil {
		logs.Warn("submit activity failed , err is %v", err)
		err = fmt.Errorf("创建活动失败， %s", err.Error())
		errMsg = err.Error()
		return
	}

	p.Redirect("/activity/list", http.StatusMovedPermanently)

	return
}