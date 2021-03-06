package access

import (
	"encoding/base64"
	"encoding/json"
	"net/http"
	"regexp"
	"strings"
	"time"
)

var forwardedForRegex = regexp.MustCompile(`for=([^;]*)(?:[; ]|$)`)
var ipRegex = regexp.MustCompile(`^(?:(\d+\.\d+\.\d+\.\d+)|\[(.*)\])(?::\d+)?$`)

const HeaderSessionExtraInfo = "x-authgear-extra-info"

type Info struct {
	InitialAccess Event `json:"initial_access"`
	LastAccess    Event `json:"last_access"`
}

type Event struct {
	Timestamp time.Time      `json:"time"`
	RemoteIP  string         `json:"ip,omitempty"`
	UserAgent string         `json:"user_agent,omitempty"`
	Extra     EventExtraInfo `json:"extra,omitempty"`
}

func NewEvent(timestamp time.Time, req *http.Request, trustProxy bool) Event {
	remote := EventConnInfo{
		RemoteAddr:    req.RemoteAddr,
		XForwardedFor: req.Header.Get("X-Forwarded-For"),
		XRealIP:       req.Header.Get("X-Real-IP"),
		Forwarded:     req.Header.Get("Forwarded"),
	}

	extra := EventExtraInfo{}
	extraData, err := base64.StdEncoding.DecodeString(req.Header.Get(HeaderSessionExtraInfo))
	const extraDataSizeLimit = 1024
	if err == nil && len(extraData) <= extraDataSizeLimit {
		_ = json.Unmarshal(extraData, &extra)
	}

	return Event{
		Timestamp: timestamp,
		RemoteIP:  remote.IP(trustProxy),
		UserAgent: req.UserAgent(),
		Extra:     extra,
	}
}

type EventConnInfo struct {
	RemoteAddr    string `json:"remote_addr,omitempty"`
	XForwardedFor string `json:"x_forwarded_for,omitempty"`
	XRealIP       string `json:"x_real_ip,omitempty"`
	Forwarded     string `json:"forwarded,omitempty"`
}

func (i EventConnInfo) IP(trustProxy bool) (ip string) {
	defer func() {
		ip = strings.TrimSpace(ip)
		// remove ports from IP
		if matches := ipRegex.FindStringSubmatch(ip); len(matches) > 0 {
			ip = matches[1]
			if len(matches[2]) > 0 {
				ip = matches[2]
			}
		}
	}()

	if trustProxy && i.Forwarded != "" {
		if matches := forwardedForRegex.FindStringSubmatch(i.Forwarded); len(matches) > 0 {
			ip = matches[1]
			return
		}
	}
	if trustProxy && i.XForwardedFor != "" {
		parts := strings.SplitN(i.XForwardedFor, ",", 2)
		ip = parts[0]
		return
	}
	if trustProxy && i.XRealIP != "" {
		ip = i.XRealIP
		return
	}
	ip = i.RemoteAddr
	return
}

type EventExtraInfo map[string]interface{}

func (i EventExtraInfo) DeviceName() string {
	deviceName, _ := i["device_name"].(string)
	return deviceName
}
