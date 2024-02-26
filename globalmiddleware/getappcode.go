package globalmiddleware

import (
	"net/http"
)

// GetAppCode 获取app_code
func GetAppCode(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		appCode := r.Header.Get("X-APP-CODE")
		r.Header.Set("app_code", appCode)
		next(w, r)
	}
}
