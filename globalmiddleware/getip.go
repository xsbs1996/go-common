package globalmiddleware

import (
	"net"
	"net/http"
	"strings"
)

// GetIP returns request real ip.
func GetIP(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ipString := getIpString(r)
		r.Header.Set("ip", ipString)
		next(w, r)
	}

}

func getIpString(r *http.Request) string {
	ip := r.Header.Get("X-Forwarded-For")
	for _, i := range strings.Split(ip, ",") {
		if net.ParseIP(i) != nil {
			return i
		}
	}

	ip = r.Header.Get("X-Real-IP")
	if net.ParseIP(ip) != nil {
		return ip
	}

	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return ""
	}

	if net.ParseIP(ip) != nil {
		return ip
	}

	return ""
}
