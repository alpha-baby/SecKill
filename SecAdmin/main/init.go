package main

import (
	"fmt"
	"github.com/astaxie/beego/logs"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

var Db *sqlx.DB

func Init() (err error) {
	//  初始化读取配置文件
	err = initConfig()
	if err != nil {
		return errors.New(fmt.Sprintf("%v", err))
	}

	// 初始化model
	err = initDb()
	if err != nil {
		return errors.New(fmt.Sprintf("%v", err))
	}
	return nil
}

func initDb() (err error) {
	mysqlConfig := AppConf.MysqlConfig
	Db, err = sqlx.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s",
		mysqlConfig.UserName,
		mysqlConfig.Passwd,
		mysqlConfig.Host,
		mysqlConfig.Port,
		mysqlConfig.Database))

	if err != nil {
		logs.Warn("open mysql failed,", err)
		return err
	}
	return nil
}