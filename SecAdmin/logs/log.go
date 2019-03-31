package logs

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
)

type LogConfig struct {
	Level string
	Path  string
}

func InitLog(conf LogConfig) (err error) {
	if strings.ToLower(conf.Level) == "debug" {

	} else {
		err = beego.BeeLogger.DelLogger("console")
		if err != nil {
			return errors.New(fmt.Sprintf("beego BeeLogger Delelte Logger console error:%v", err))
		}
		config := make(map[string]interface{})
		config["filename"] = conf.Path
		config["level"] = convertLogLevel(conf.Level)

		configStr, err1 := json.Marshal(config)
		if err1 != nil {
			return errors.New(fmt.Sprintf("beego set new file logger error,json Marshal the config of logger :%v", err))
		}
		err = beego.SetLogger(logs.AdapterFile, string(configStr))
		if err != nil {
			return errors.New(fmt.Sprintf("beego set new file logger error:%v", err))
		}
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
