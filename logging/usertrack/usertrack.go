package usertrack

import (
	"bytes"
	"net"
	"net/url"
	"strconv"
	"strings"
)

type Slot string

const (
	SlotDeviceID          Slot = "device_id"
	SlotUserID                 = "user_id"
	SlotUsername               = "username"
	SlotMobilePhone            = "mobile_phone"
	SlotMacAddr                = "mac_addr"
	SlotIpAddr                 = "ip_addr"
	SlotSessionID              = "session_id"
	SlotAccessToken            = "access_token"
	SlotCanvasFingerprint      = "canvas_fingerprint"
)

type SlotType string

const (
	SlotTypeLoginFail  SlotType = "login_fail"
	SlotTypeSmsSend             = "sms_send"
	SlotTypeSmsOtpSent          = "sms_otp_sent"
)

type Options struct {
	IpDatabase string
}

type Feature struct {
	HttpRequestID string `json:"http_request_id,omitempty"` // HTTP 请求ID

	IpAddr             string       `json:"ip_addr,omitempty"`              // IP 地址
	RemoteAddr         string       `json:"remote_addr,omitempty"`          // 远端地址(IP:Port)
	WebrtcAddrs        []string     `json:"webrtc_addrs,omitempty"`         // WEBRTC 地址（未代理的真实地址）
	WebrtcPublicAddrs  []IpLocation `json:"webrtc_public_addrs,omitempty"`  // WEBRTC 公网地址（未代理的真实地址）
	WebrtcPrivateAddrs []string     `json:"webrtc_private_addrs,omitempty"` // WEBRTC 内网地址（未代理的真实地址）

	MacAddr string `json:"mac_addr,omitempty"` // MAC 地址

	WifiSSID  string `json:"wifi_ssid,omitempty"`  // Wi-Fi SSID 名
	WifiBSSID string `json:"wifi_bssid,omitempty"` // Wi-Fi MAC 地址

	AppVersion    string `json:"app_version,omitempty"`     // App 版本
	AppBundleID   string `json:"app_bundle_id,omitempty"`   // App bundle id
	AppBundleName string `json:"app_bundle_name,omitempty"` // App bundle name

	DeviceID            string `json:"device_id,omitempty"`              // 设备 ID
	DeviceName          string `json:"device_name,omitempty"`            // 设备名称
	DeviceModel         string `json:"device_model,omitempty"`           // 设备信号
	DeviceLocale        string `json:"device_locale,omitempty"`          // 设备语言
	DevicePrivateIpAddr string `json:"device_private_ip_addr,omitempty"` // 设备内网地址

	CanvasFingerprint string `json:"canvas_fingerprint,omitempty"` // Canvas 指纹

	SessionID   string `json:"session_id,omitempty"`   // 会话 ID
	AccessToken string `json:"access_token,omitempty"` // 访问令牌

	UserID   int64  `json:"user_id,omitempty"`  // 用户ID
	Username string `json:"username,omitempty"` // 用户名

	UserNodeID int64 `json:"user_node_id,omitempty"` // 用户ID

	UserAgent        string   `json:"useragent,omitempty"` // UserAgent
	UserAgentName    string   `json:"useragent_name,omitempty"`
	UserAgentVersion string   `json:"useragent_version,omitempty"`
	UserAgentPlugins []string `json:"useragent_plugins,omitempty"`
	ScreenResolution string   `json:"screen_resolution,omitempty"` // 屏幕分辨率

	PlatformName    string `json:"platform_name,omitempty"`    // 操作系统名称
	PlatformVersion string `json:"platform_version,omitempty"` // 操作系统版本

	IpAddrLocation *IpLocation `json:"ip_addr_location,omitempty"` // IP 地址位置
	DeviceLocation *IpLocation `json:"device_location,omitempty"`  // 设备位置
}

type IpLocation struct {
	IpAddr string `json:"ip_addr,omitempty"` // IP 地址

	Continent string `json:"continent,omitempty"` // 大洲
	Country   string `json:"country,omitempty"`   // 国家
	Province  string `json:"province,omitempty"`  // 省
	City      string `json:"city,omitempty"`      // 市
	District  string `json:"district,omitempty"`  // 县/区

	CountryEn   string `json:"country_en,omitempty"`   // 国家
	CountryCode string `json:"country_code,omitempty"` // 国家代码

	DmaCode string `json:"dma_code,omitempty"` // 行政规划代码
	ISP     string `json:"isp,omitempty"`      // 电信运营商

	Latitude  float64 `json:"latitude,omitempty"`  // 纬度
	Longitude float64 `json:"longitude,omitempty"` // 经度
	// GeoPoint  *elastic.GeoPoint `json:"geo_point,omitempty"`
}

func Parse(m map[string]interface{}) Feature {
	// TODO 完善请求日志记录

	webrtcAddrs := make([]string, 0)
	webrtcPublicAddrs := make([]IpLocation, 0)
	webrtcPrivateAddrs := make([]string, 0)

	userAgentPlugins, _ := url.PathUnescape(toString(m["X-User-Agent-Plugins"]))
	deviceName, _ := url.QueryUnescape(toString(m["X-Device-Name"]))
	bundleID := toString(m["X-App-Bundle-Id"])
	if bundleID == "" {
		bundleID = toString(m["X-Bundle-Id"])
	}
	f := Feature{
		HttpRequestID: toString(m["_request_id"]),

		IpAddr:  toString(m["_client_ip"]),
		MacAddr: toString(m["X-Client-Mac"]),

		AppVersion:    toString(m["X-App-Version"]),
		AppBundleID:   bundleID,
		AppBundleName: toString(m["X-App-Bundle-Name"]),

		WifiSSID:  toString(m["X-Device-Wifi-Ssid"]),
		WifiBSSID: toString(m["X-Device-Wifi-Bssid"]),

		DeviceID:            toString(m["X-Device-Id"]),
		DeviceName:          deviceName,
		DeviceModel:         toString(m["X-Device-Model"]),
		DeviceLocale:        toString(m["X-Device-Locale"]),
		DevicePrivateIpAddr: toString(m["X-Device-Private-Ip-Addr"]),

		PlatformName:    toString(m["X-Platform-Name"]),
		PlatformVersion: toString(m["X-Platform-Version"]),

		CanvasFingerprint: toString(m["X-Canvas-Fingerprint"]),
		ScreenResolution:  toString(m["X-Screen-Resolution"]),

		UserAgent:        toString(m["User-Agent"]),
		UserAgentPlugins: strings.Split(userAgentPlugins, ";"),

		WebrtcAddrs:        webrtcAddrs,
		WebrtcPublicAddrs:  webrtcPublicAddrs,
		WebrtcPrivateAddrs: webrtcPrivateAddrs,
	}

	// if !isIPv6(f.IpAddr) && ipParser != nil {
	if !isIPv6(f.IpAddr) {
	}
	f.IpAddrLocation = &IpLocation{
		Country: "未知",
	}

	if deviceLocation := toString(m["X-Device-Location"]); deviceLocation != "" {
		s := strings.Split(deviceLocation, ",")
		if len(s) == 2 {
			var dLocation IpLocation
			dLocation.Latitude, _ = strconv.ParseFloat(s[0], 64)
			dLocation.Longitude, _ = strconv.ParseFloat(s[1], 64)
			f.DeviceLocation = &dLocation
		}

	}

	f.RemoteAddr = f.IpAddr + ":" + toString(m["X-Remote-Port"])

	return f
}

func toString(s interface{}) string {
	v, ok := s.(string)
	if !ok {
		return ""
	}
	return v
}

func isIPv6(str string) bool {
	ip := net.ParseIP(str)
	return ip != nil && strings.Contains(str, ":")
}

// inRange - check to see if a given ip address is within a range given
func inRange(r ipRange, ipAddress net.IP) bool {
	// strcmp type byte comparison
	if bytes.Compare(ipAddress, r.start) >= 0 && bytes.Compare(ipAddress, r.end) < 0 {
		return true
	}
	return false
}

// ipRange - a structure that holds the start and end of a range of ip addresses
type ipRange struct {
	start net.IP
	end   net.IP
}
