package logs

import (
	"fmt"
	"strings"

	BeegoLogs "github.com/astaxie/beego/logs"
)

func logLevelName2Int(levelName string) (logLevel int, err error) {

	switch strings.ToLower(levelName) {
	case "emergency":
		logLevel = BeegoLogs.LevelEmergency
	case "alert":
		logLevel = BeegoLogs.LevelAlert
	case "critical":
		logLevel = BeegoLogs.LevelCritical
	case "error":
		logLevel = BeegoLogs.LevelError
	case "warnning":
		logLevel = BeegoLogs.LevelWarning
	case "notice":
		logLevel = BeegoLogs.LevelNotice
	case "info", "informational":
		logLevel = BeegoLogs.LevelInformational
	case "debug":
		logLevel = BeegoLogs.LevelDebug
	case "trace":
		logLevel = BeegoLogs.LevelTrace
	default:
		err = fmt.Errorf("level name[%s] is invalid", levelName)
	}

	return logLevel, err
}
