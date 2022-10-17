package klog

import (
	"bytes"
	"fmt"
	"path/filepath"

	"github.com/sirupsen/logrus"
)

type MyFormatter struct{}

func (m *MyFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	var b *bytes.Buffer
	if entry.Buffer != nil {
		b = entry.Buffer
	} else {
		b = &bytes.Buffer{}
	}

	timestamp := entry.Time.Format("2006-01-02 15:04:05.999")
	var newLog string

	//HasCaller()为true才会有调用信息
	if entry.HasCaller() {
		fName := entry.Caller.File
		if ignorePath == nil || len(ignorePath) <= 0 {
			fName = filepath.Base(entry.Caller.File)
		} else {
			for _, v := range ignorePath {
				fName = v.ReplaceAllString(fName, "")
			}
		}
		newLog = fmt.Sprintf("[%s] [%s] [%s:%d] [%s] %s\n",
			timestamp, entry.Level, fName, entry.Caller.Line, getLogData(entry.Data), entry.Message)
	} else {
		newLog = fmt.Sprintf("[%s] [%s] %s\n", timestamp, entry.Level, entry.Message)
	}

	b.WriteString(newLog)
	return b.Bytes(), nil
}
