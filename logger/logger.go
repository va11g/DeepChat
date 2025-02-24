package logger

import (
	"fmt"
	"io"
	"log"
	"logger/config"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"sync"
)

type LogLevel int
type LogConfig struct {
	Path  string
	Name  string
	Level LogLevel
}

const (
	REPLY LogLevel = iota
	INFO
	WARNING
	ERROR
	PANIC
)

var (
	logFile            *os.File
	logger             *log.Logger
	logMu              sync.Mutex
	levelLabels        = []string{"reply", "info", "warning", "error", "panic"}
	logCfg             *LogConfig
	defaultCallerDepth = 2
	logPrefix          = ""
)

func SetUp(cfg *config.Config) error {
	var err error
	logCfg = &LogConfig{
		Path:  cfg.LogDir,
		Name:  "redis.log",
		Level: INFO,
	}
	for i, v := range levelLabels {
		if v == cfg.LogLevel {
			logCfg.Level = LogLevel(i)
			break
		}
	}

	if _, err = os.Stat(logCfg.Path); err != nil {
		mkErr := os.Mkdir(logCfg.Path, 0755)
		if mkErr != nil {
			return mkErr
		}
	}

	logfile := path.Join(logCfg.Path, logCfg.Name)
	logFile, err = os.OpenFile(logfile, os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		return err
	}
	writer := io.MultiWriter(os.Stdout, logFile)
	logger = log.New(writer, "", log.LstdFlags)
	return nil
}

// 前缀
func setPrefix(level LogLevel) {
	_, file, line, ok := runtime.Caller(defaultCallerDepth)
	if ok {
		logPrefix = fmt.Sprintf("[%s][%s:%d] ", levelLabels[level], filepath.Base(file), line)
	} else {
		logPrefix = fmt.Sprintf("[%s] ", levelLabels[level])
	}
	logger.SetPrefix(logPrefix)
}

func Reply(v ...any) {
	if logCfg.Level > REPLY {
		return
	}
	logMu.Lock()
	defer logMu.Unlock()
	setPrefix(REPLY)
	logger.Println(v)
}

func Info(v ...any) {
	if logCfg.Level > INFO {
		return
	}
	logMu.Lock()
	defer logMu.Unlock()
	setPrefix(INFO)
	logger.Println(v)
}

func Warning(v ...any) {
	if logCfg.Level > WARNING {
		return
	}
	logMu.Lock()
	defer logMu.Unlock()
	setPrefix(WARNING)
	logger.Println(v)
}

func Error(v ...any) {
	if logCfg.Level > ERROR {
		return
	}
	logMu.Lock()
	defer logMu.Unlock()
	setPrefix(ERROR)
	logger.Println(v)
}

func Panic(v ...any) {
	if logCfg.Level > PANIC {
		return
	}
	logMu.Lock()
	defer logMu.Unlock()
	setPrefix(PANIC)
	logger.Println(v)
}
