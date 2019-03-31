package main

import (
	"fmt"
	"github.com/astaxie/beego"
	"github.com/pkg/errors"
	MyLogs "SecKill/SecAdmin/logs"
	"strings"
)

type MysqlConfig struct {
	UserName string
	Passwd   string
	Port     string
	Database string
	Host     string
}

type EtcdConf struct {
	Addr string
	EtcdKeyPrefix string
	ProductKey	string
	Timeout int64
	EtcdKey string
}

type Config struct {
	MysqlConfig MysqlConfig
	LogConfig   MyLogs.LogConfig
	EtcdConf EtcdConf
}

var AppConf Config

func initConfig() (err error) {
	// mysql
	err = loadMysqlConfig()
	if err != nil {
		return err
	}

	// log config
	err = loadLogConfig()
	if err != nil {
		return err
	}

	// etcd
	err = loadEtcdConfig()
	if err != nil {
		return err
	}
	return nil
}

func loadMysqlConfig() (err error) {
	AppConf.MysqlConfig.Host = beego.AppConfig.String("mysql_host")
	if len(AppConf.MysqlConfig.Host) == 0 {
		return errors.New("init config mysql mysql_host err")
	}

	AppConf.MysqlConfig.Database = beego.AppConfig.String("mysql_database")
	if len(AppConf.MysqlConfig.Database) == 0 {
		return errors.New("init config mysql mysql_database err")
	}

	AppConf.MysqlConfig.Passwd = beego.AppConfig.String("mysql_passwd")
	if len(AppConf.MysqlConfig.Passwd) == 0 {
		return errors.New("init config mysql mysql_passwd err")
	}

	AppConf.MysqlConfig.Database = beego.AppConfig.String("mysql_database")
	if len(AppConf.MysqlConfig.Database) == 0 {
		return errors.New("init config mysql mysql_database err")
	}

	AppConf.MysqlConfig.Port = beego.AppConfig.String("mysql_port")
	if len(AppConf.MysqlConfig.Port) == 0 {
		return errors.New("init config mysql mysql_port err")
	}

	AppConf.MysqlConfig.UserName = beego.AppConfig.String("mysql_user_name")
	if len(AppConf.MysqlConfig.Port) == 0 {
		return errors.New("init config mysql mysql_user_name err")
	}
	return nil
}

func loadLogConfig() (err error) {
	AppConf.LogConfig.Level = beego.AppConfig.String("log_level")
	if len(AppConf.MysqlConfig.Port) == 0 {
		return errors.New("init config mysql log_level err")
	}

	AppConf.LogConfig.Path = beego.AppConfig.String("log_path")
	if len(AppConf.MysqlConfig.Port) == 0 {
		return errors.New("init config mysql log_path err")
	}
	return nil
}

func loadEtcdConfig() (err error) {
	AppConf.EtcdConf.Addr = beego.AppConfig.String("etcd_addr")
	if len(AppConf.EtcdConf.Addr) == 0 {
		return errors.New("init config etcd etcd_addr err")
	}

	AppConf.EtcdConf.EtcdKeyPrefix = beego.AppConfig.String("etcd_sec_key_prefix")
	if len(AppConf.EtcdConf.EtcdKeyPrefix) == 0 {
		return errors.New("init config etcd etcd_sec_key_prefix err")
	}

	AppConf.EtcdConf.ProductKey = beego.AppConfig.String("etcd_product_key")
	if len(AppConf.EtcdConf.ProductKey) == 0 {
		return errors.New("init config etcd etcd_product_key err")
	}

	AppConf.EtcdConf.Timeout, err = beego.AppConfig.Int64("etcd_timeout")
	if err != nil || AppConf.EtcdConf.Timeout <= 0 {
		return errors.New("init config etcd etcd_timeout err")
	}

	if strings.HasSuffix(AppConf.EtcdConf.EtcdKeyPrefix, "/") == false {
		AppConf.EtcdConf.EtcdKeyPrefix = AppConf.EtcdConf.EtcdKeyPrefix + "/"
	}
	AppConf.EtcdConf.EtcdKey = fmt.Sprintf("%s%s", AppConf.EtcdConf.EtcdKeyPrefix, AppConf.EtcdConf.ProductKey)
	return nil
}