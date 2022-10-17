package util

import (
	"encoding/json"
	"net"
	"net/http"
	"strings"
	"unsafe"
)

// ToJson 序列化为json
func ToJson(data interface{}) string {
	b, _ := json.Marshal(data)
	return *(*string)(unsafe.Pointer(&b))
}

func ClientIP(r *http.Request) string {
	// 容器环境下获取真实IP
	clientIP := getFirstIpAddr(r.Header.Get("X-Original-Forwarded-For"))
	if len(clientIP) > 0 {
		return clientIP
	}

	clientIP = getFirstIpAddr(r.Header.Get("X-Forwarded-For"))
	if len(clientIP) > 0 {
		return clientIP
	}

	clientIP = getFirstIpAddr(r.Header.Get("X-Client-Ip"))
	if len(clientIP) > 0 {
		return clientIP
	}

	clientIP = getFirstIpAddr(r.Header.Get("Cdn-Src-Ip"))
	if len(clientIP) > 0 {
		return clientIP
	}

	clientIP = getFirstIpAddr(r.Header.Get("X-Real-Ip"))
	if len(clientIP) > 0 {
		return clientIP
	}

	if ip, _, err := net.SplitHostPort(strings.TrimSpace(r.RemoteAddr)); err == nil {
		return ip
	}
	return ""
}

func getFirstIpAddr(clientIP string) string {
	if index := strings.IndexByte(clientIP, ','); index >= 0 {
		clientIP = clientIP[0:index]
	}
	clientIP = strings.TrimSpace(clientIP)
	if len(clientIP) > 0 {
		return clientIP
	}
	return ""
}
