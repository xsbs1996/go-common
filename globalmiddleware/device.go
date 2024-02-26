package globalmiddleware

import (
	"net/http"
	"strings"
)

// GetDevice 获取设备信息
func GetDevice(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		deviceInfo := strings.Split(r.Header.Get("X-DEVICE-INFO"), "_")
		var device string
		if len(deviceInfo) >= 1 {
			device = deviceInfo[0]
		}
		r.Header.Set("device", device)
		next(w, r)
	}

}
