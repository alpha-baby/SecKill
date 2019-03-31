package main

import (
	"SecKill/SecAdmin/model"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/astaxie/beego/logs"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	MyLogs "SecKill/SecAdmin/logs"
	"go.etcd.io/etcd/clientv3"
	"strings"
	"time"
)

func Init() (err error) {
	//  初始化读取配置文件
	err = initConfig()
	if err != nil {
		return errors.New(fmt.Sprintf(" 初始化读取配置文件 err is %v", err))
	}

	// 初始化日志
	err = MyLogs.InitLog(AppConf.LogConfig)
	if err != nil {
		return errors.New(fmt.Sprintf("初始化日志 err is %v", err))
	}
	logs.Debug("logger init success")

	// 初始化model
	DbMysql, err := initDb()
	if err != nil {
		return errors.New(fmt.Sprintf("初始化model err is %v", err))
	}
	model.InitDB(DbMysql)
	logs.Debug("db mysql init success")

	// 初始化 etcd
	etcdClient, err := initEtcd()
	if err != nil {
		return errors.New(fmt.Sprintf("初始化 etcd err is %v", err))
	}
	model.InitEtcd(etcdClient, AppConf.EtcdConf.EtcdKeyPrefix, AppConf.EtcdConf.EtcdKey)
	logs.Debug("etcd init success")

	return nil
}

func initDb() (Db *sqlx.DB, err error){
	mysqlConfig := AppConf.MysqlConfig
	Db, err = sqlx.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s",
		mysqlConfig.UserName,
		mysqlConfig.Passwd,
		mysqlConfig.Host,
		mysqlConfig.Port,
		mysqlConfig.Database),
	)

	if err != nil {
		logs.Warn("open mysql failed,", err)
		return nil, err
	}

	err = Db.Ping()

	if err != nil {
		logs.Warn("ping mysql failed,", err)
		return nil, err
	}

	return Db, nil
}

func initEtcd() (c *clientv3.Client, err error) {
	c, err = clientv3.New(clientv3.Config{
		Endpoints: strings.Split(AppConf.EtcdConf.Addr, ","),
		DialTimeout: time.Duration(AppConf.EtcdConf.Timeout) * time.Second,
	})

	if err != nil {
		logs.Error("init etcd err, Addr is %v, err is %v", AppConf.EtcdConf.Addr, err)
		return nil, err
	}

	return c,nil
}