package klog

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
	"time"

	"github.com/go-kratos/kratos/v2/log"
	rotatelogs "github.com/lestrrat/go-file-rotatelogs"
	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"
)

var (
	isLogFile  bool                   // 是否写入文件
	logDir     string                 // 目录
	logFile    string                 // 文件名
	ignorePath map[int]*regexp.Regexp // 忽略路径前缀
)

// Level 日志级别
type Level = string

const (
	LevelDebug Level = "debug"
	LevelInfo  Level = "info"
	LevelFatal Level = "fatal"
	LevelError Level = "error"
	LevelWarn  Level = "warn"
	LevelPanic Level = "panic"
)

var (
	allLevels = []string{LevelDebug, LevelInfo, LevelFatal, LevelError, LevelWarn, LevelPanic}
)

type Logger struct {
	log *logrus.Entry
}

var _ log.Logger = (*Logger)(nil)

// Log Implementation of logger interface.
func (l *Logger) Log(level log.Level, keyVals ...interface{}) error {
	// l.log.Info(level, keyVals)
	logLevel, _ := logrus.ParseLevel(level.String())
	l.log.Log(logLevel, keyVals...)
	return nil
}

//func NewLogger(prefix string, level logrus.Level) *Logger {
//	l := newLogger()
//	l.SetReportCaller(true)
//	return &Logger{
//		log: l.WithField("prefix", prefix),
//	}
//}

func NewLogger(prefix string, level Level) *Logger {
	l := newLogger()
	l.SetReportCaller(true)
	if lvl, err := logrus.ParseLevel(level); err != nil {
		l.Panicf("Please set log-level one of %v instead of %s.\n\n",
			allLevels, lvl)
	} else {
		l.Level = lvl
	}
	return &Logger{
		log: l.WithField("prefix", prefix),
	}
}

func getLogData(data logrus.Fields) string {
	if data == nil || len(data) <= 0 {
		return " "
	}
	payload := strings.Builder{}
	idx := 0
	for k, v := range data {
		if idx > 0 {
			payload.WriteString(" ")
		}
		payload.WriteString(k + "=")
		payload.WriteString(fmt.Sprintf("%v", v))
		idx++
	}
	return payload.String()
}

// SetIgnorePath 设置忽略路径
// _ignorePath 目录 例：/Users/ha666/gopath/src/git.ztosys.com/ZTO_CS/go-contrib/
func SetIgnorePath(_ignorePath []string) {
	if _ignorePath != nil && len(_ignorePath) > 0 {
		ignorePath = make(map[int]*regexp.Regexp)
		for i, v := range _ignorePath {
			ignorePath[i] = regexp.MustCompile(v)
		}
	}
}

// SetFileLogger 设置日志写入文件
// _logDir 目录 例：/data/logs/tenant-sso
// _logFile 文件名 例：sso.log
func SetFileLogger(_logDir, _logFile string) {
	if _logDir == "" || _logFile == "" {
		return
	}
	isLogFile = true
	logDir = _logDir
	logFile = _logFile
}

func newLogger() *logrus.Logger {
	l := logrus.New()
	l.AddHook(NewLHook(2))
	if isLogFile {
		src, err := os.OpenFile(os.DevNull, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
		if err != nil {
			panic("打开文件出错:" + err.Error())
		}
		l.Out = src
		filePatten := fmt.Sprintf("%s/%s", logDir, logFile)
		filePatten = strings.ReplaceAll(filePatten, ".log", ".%Y-%m-%d.log")
		logWriter, _ := rotatelogs.New(
			filePatten,
			rotatelogs.WithMaxAge(2*24*time.Hour),     // 文件最大保存时间
			rotatelogs.WithRotationTime(time.Hour*24), // 日志切割时间间隔
			rotatelogs.WithClock(rotatelogs.Local),
			rotatelogs.WithLocation(time.Local),
		)
		writeMap := lfshook.WriterMap{
			logrus.PanicLevel: logWriter,
			logrus.FatalLevel: logWriter,
			logrus.ErrorLevel: logWriter,
			logrus.WarnLevel:  logWriter,
			logrus.InfoLevel:  logWriter,
			logrus.DebugLevel: logWriter,
			logrus.TraceLevel: logWriter,
		}
		lfHook := lfshook.NewHook(writeMap, &MyFormatter{})
		l.AddHook(lfHook)
	} else {
		l.Formatter = &logrus.TextFormatter{
			TimestampFormat: "2006-01-02 15:04:05.999",
			ForceColors:     true,
			FullTimestamp:   true,
			CallerPrettyfier: func(frame *runtime.Frame) (function string, file string) {
				fName := frame.File
				if ignorePath == nil || len(ignorePath) <= 0 {
					fName = filepath.Base(frame.File)
				} else {
					for _, v := range ignorePath {
						fName = v.ReplaceAllString(fName, "")
					}
				}
				return "", fmt.Sprintf("%s:%d", fName, frame.Line)
			},
		}
		l.SetOutput(os.Stdout)
	}
	return l
}
