package klog

import (
	"fmt"
	"runtime"
	"strings"

	"github.com/sirupsen/logrus"
)

type LHook struct {
	Field     string
	Skip      int
	Jumped    int
	levels    []logrus.Level
	Formatter func(file, function string, line int) string
}

func (l *LHook) Levels() []logrus.Level {
	return l.levels
}

func (l *LHook) Fire(entry *logrus.Entry) error {
	entry.Data[l.Field] = l.Formatter(findCaller(l.Skip, l.Jumped))
	return nil
}

func findCaller(skip, jumped int) (string, string, int) {
	var (
		pc       uintptr
		file     string
		function string
		line     int
	)
	for i := 0; i < 15; i++ {
		pc, file, line = getCaller(skip + i)
		if !strings.HasPrefix(file, "log") && !strings.HasPrefix(file, "logrus") {
			pc, file, line = getCaller(skip + i + jumped)
			break
		}
	}
	if pc != 0 {
		frames := runtime.CallersFrames([]uintptr{pc})
		frame, _ := frames.Next()
		function = frame.Function
	}

	return file, function, line
}

func getCaller(skip int) (uintptr, string, int) {
	pc, file, line, ok := runtime.Caller(skip)
	if !ok {
		return 0, "", 0
	}

	n := 0
	for i := len(file) - 1; i > 0; i-- {
		if file[i] == '/' {
			n += 1
			if n >= 2 {
				file = file[i+1:]
				break
			}
		}
	}

	return pc, file, line
}

func NewLHook(jumped int, levels ...logrus.Level) logrus.Hook {
	hook := LHook{
		Field:  "caller",
		Skip:   3,
		Jumped: jumped,
		levels: levels,
		Formatter: func(file, function string, line int) string {
			return fmt.Sprintf("%s:%d", file, line)
		},
	}
	if len(hook.levels) == 0 {
		hook.levels = logrus.AllLevels
	}
	return &hook
}
