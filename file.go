package logs

import (
	"encoding/json"
	"fmt"

	BeegoLogs "github.com/astaxie/beego/logs"
)

type MultiFileLogWriter struct {
	Maxlines int `json:"maxlines"`
	// Rotate at size
	Maxsize int `json:"maxsize"`
	// Rotate daily
	Daily   bool  `json:"daily"`
	Maxdays int64 `json:"maxdays"`

	Rotate    bool   `json:"rotate"`
	Level     int    `json:"level"`
	LevelName string `json:"levelname"`

	LevelFiles []*struct {
		LevelNames []string `json:"levelnames"`
		Levels     []int
		FileName   string `json:"filename"`
	} `json:"levelfiles"`

	levelLoggerMap map[int]BeegoLogs.LoggerInterface
}

func NewMultiFileLogWriter() BeegoLogs.LoggerInterface {
	return &MultiFileLogWriter{
		Maxlines: 1000000,
		Maxsize:  1 << 28, //256 MB
		Daily:    true,
		Maxdays:  7,
		Rotate:   true,
		Level:    BeegoLogs.LevelTrace,

		levelLoggerMap: make(map[int]BeegoLogs.LoggerInterface),
	}
}

func (w *MultiFileLogWriter) Init(jsonconfig string) error {
	err := json.Unmarshal([]byte(jsonconfig), w)
	if err != nil {
		return err
	}
	if len(w.LevelName) > 0 {
		tmp, err := logLevelName2Int(w.LevelName)
		if err == nil {
			w.Level = tmp
		}
	}
	if len(w.LevelFiles) == 0 {
		return fmt.Errorf("levelfiles is empty")
	}
	for _, l := range w.LevelFiles {
		if len(l.FileName) == 0 {
			return fmt.Errorf("filename in levelfiles is empty")
		}
		if len(l.LevelNames) == 0 {
			return fmt.Errorf("levelnames for file[%s] is empty", l.FileName)
		}
		l.Levels = make([]int, 0, len(l.LevelNames))
		for _, levelName := range l.LevelNames {
			level, err := logLevelName2Int(levelName)
			if err != nil {
				return err
			}
			l.Levels = append(l.Levels, level)
		}
	}

	w.initInnerLoggers()

	return nil
}

func (w *MultiFileLogWriter) initInnerLoggers() error {
	for _, l := range w.LevelFiles {
		config := fmt.Sprintf(`
		{
			"filename": "%s",
			"maxlines": %d,
			"maxsize": %d,
			"daily": %v,
			"maxdays": %d,
			"rotate": %v,
			"level": %d
		}
		`, l.FileName, w.Maxlines, w.Maxsize, w.Daily, w.Maxdays, w.Rotate, BeegoLogs.LevelDebug /*use highest log level*/)
		innerLogger := BeegoLogs.NewFileWriter()
		err := innerLogger.Init(config)
		if err != nil {
			return err
		}
		for _, level := range l.Levels {
			w.levelLoggerMap[level] = innerLogger
		}
	}

	return nil
}

func (w *MultiFileLogWriter) WriteMsg(msg string, level int) error {
	if level > w.Level {
		return nil
	}
	logger, exist := w.levelLoggerMap[level]
	if !exist {
		return nil
	}
	return logger.WriteMsg(msg, level)
}

func (w *MultiFileLogWriter) Destroy() {
	for _, logger := range w.levelLoggerMap {
		logger.Destroy()
	}
}

func (w *MultiFileLogWriter) Flush() {
	for _, logger := range w.levelLoggerMap {
		logger.Flush()
	}
}

func init() {
	BeegoLogs.Register("multi_file", NewMultiFileLogWriter)
}
