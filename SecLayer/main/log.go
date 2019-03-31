package main

import (
	"encoding/json"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"log"
	"strings"
)

// init loger
func initLogger() (err error) {
	if strings.ToLower(appConfig.LogLevel) == "debug" {
		beego.SetLogFuncCall(true)
	}else {
		config := make(map[string]interface{})
		config["filename"] = appConfig.LogPath
		config["level"] = convertLogLevel(appConfig.LogLevel)

		configStr, err := json.Marshal(config)
		if err != nil {
			log.Println("init loger fialed, json marshal error ,", err)
			return err
		}

		//adapter := logs.AdapterFile
		//if config["level"] == logs.LevelDebug {
		//	adapter = logs.AdapterConsole
		//}
		err = beego.SetLogger(logs.AdapterFile, string(configStr))
		if err != nil {
			log.Println("init loger fialed, SetLogger error,", err)
			return err
		}
		beego.SetLogFuncCall(true)
	}
	return nil
}

// loglevel string convert to int
func convertLogLevel(level string) int {
	level = strings.ToLower(level)
	switch level {
	case "debug":
		return logs.LevelDebug
	case "warn":
		return logs.LevelWarn
	case "info":
		return logs.LevelInfo
	case "trace":
		return logs.LevelTrace
	default:
		return logs.LevelDebug
	}
}
