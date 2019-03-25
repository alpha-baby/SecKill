package model

import (
	"github.com/astaxie/beego/logs"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

type ProductModel struct {
	Db *sqlx.DB
}

type Product struct {
	ProductId   int    `db:"id"`
	ProdcutName string `db:"name"`
	Total       int    `db:"total"`
	Status      int    `db:"status"`
}

func NewProductModel(db *sqlx.DB) *ProductModel {
	productModel := &ProductModel{
		Db: db,
	}
	return productModel
}

func (p *ProductModel) GetProductList() (list []*Product, err error){
	sql := "select id, name, total, status from product"
	err = p.Db.Select(list,sql)
	if err != nil {
		logs.Warn("select from mysql failed ,err is %v sql is %v", err, sql)
		return
	}
}
