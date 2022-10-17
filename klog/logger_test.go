package klog

import (
	"testing"

	"github.com/go-kratos/kratos/v2/log"
)

func TestLog(t *testing.T) {
	SetIgnorePath([]string{`/\S+kratos-middleware[0-9\-]*/`})
	//SetFileLogger(
	//	"/Users/ha666/v/abc",
	//	"abc.log")
	l := NewLogger("test", "debug")
	l.Log(log.LevelInfo, "初始化应用出错:", "abc", "jsdfl")
}
