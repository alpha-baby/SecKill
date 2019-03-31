package product

import (
	"SecKill/SecAdmin/model"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
)

type ProductController struct {
	beego.Controller
}

func (p *ProductController) ListProduct() {
	productModel := model.NewProductModel()
	productList, err := productModel.GetProductList()
	if err != nil {
		logs.Warn("get product list failed, err:%v", err)
		return
	}

	p.Data["product_list"] = productList
	p.TplName = "product/list.html"
	p.Layout = "layout/layout.html"
}


func (p *ProductController) CreateProduct() {
	p.TplName = "product/create.html"
	p.Layout = "layout/layout.html"
}

func (p *ProductController) SubmitProduct() {
	errMsg := ""
	productName := p.GetString("product_name")
	if len(productName) == 0 {
		logs.Warn("invalid submit product get string productName")
		errMsg = "商品名不能为空"
		return
	}
	productTotal, err := p.GetInt("product_total")
	if err != nil {
		logs.Warn("invalid submit product get int product_total err: %v", err)
		errMsg =  "商品数量填写错误"
		return
	}
	productStatus, err := p.GetInt("product_status")
	if err != nil {
		logs.Warn("invalid submit product get int product_status err: %v", err)
		errMsg = "商品状态填写错误"
		return
	}

	defer func() {
		if err != nil {
			p.Data["Error"] = errMsg
			p.TplName = "product/error.html"
			p.Layout = "layout/layout.html"
		}else {
			p.TplName = "product/create.html"
			p.Layout = "layout/layout.html"
		}
	}()

	productModel := model.NewProductModel()
	product := model.Product{
		ProductName: productName,
		Total: productTotal,
		Status: productStatus,
	}

	err = productModel.CreateProduct(&product)
	if err != nil {
		logs.Warn("create product failed, err:%v", err)
		errMsg = "创建商品错误"
		return
	}

	logs.Debug("product name[%s], product total[%d], product status[%d]", productName, productTotal, productStatus)
}