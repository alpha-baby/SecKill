package main

import (
	"github.com/astaxie/beego"
	"github.com/pkg/errors"
)

type MysqlConfig struct {
	UserName string
	Passwd   string
	Port     string
	Database string
	Host     string
}

type Config struct {
	MysqlConfig MysqlConfig
}

var AppConf Config

func initConfig() (err error) {
	// mysql
	err = loadMysql()
	if err != nil {
		return err
	}
	return nil
}

func loadMysql() (err error) {
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
