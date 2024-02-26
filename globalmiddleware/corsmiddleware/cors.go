package corsmiddleware

import (
	"fmt"
	"net/http"
)

// CorsMiddleware 跨域请求处理中间件
type CorsMiddleware struct{}

const DefaultHandler = "Content-Length,X-CSRF-Token,Token,session,X_Requested_With,Accept,Origin,Host,Connection,Upgrade,Accept-Encoding,Accept-Language,DNT,X-CustomHeader,Keep-Alive,User-Agent,X-Requested-With,If-Modified-Since,Cache-Control,Content-Type,Pragma,lang,Front"
const EjxcHandler = "Authorization,X-Is-Organic-Install,X-Is-Vpn,X-DEVICE-INFO,X-APP-CODE,X-Appsflyer-ID,X-Install-Referrer"

// NewCorsMiddleware 新建跨域请求处理中间件
func NewCorsMiddleware() *CorsMiddleware {
	return &CorsMiddleware{}
}

// Handle 跨域请求处理
func (m *CorsMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		setHeader(w, fmt.Sprintf("%s,%s", DefaultHandler, EjxcHandler))
		// 放行所有 OPTIONS 方法
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		// 处理请求
		next(w, r)
	}
}

// Handler 跨域请求处理器
func (m *CorsMiddleware) Handler(headers string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		allowHeaders := fmt.Sprintf("%s,%s", DefaultHandler, EjxcHandler)
		if len(headers) > 0 {
			allowHeaders = fmt.Sprintf("%s,%s", allowHeaders, headers)
		}
		setHeader(w, allowHeaders)
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
		} else {
			w.WriteHeader(http.StatusNotFound)
		}
	})
}

// setHeader 设置响应头
func setHeader(w http.ResponseWriter, allowHeaders string) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST,GET,OPTIONS,PUT,DELETE,UPDATE")
	w.Header().Set("Access-Control-Allow-Headers", allowHeaders)
	w.Header().Set("Access-Control-Expose-Headers", "Content-Length,Access-Control-Allow-Origin,Access-Control-Allow-Headers,Cache-Control,Content-Language,Content-Type,Expires,Last-Modified,Pragma,FooBar")
	w.Header().Set("Access-Control-Max-Age", "172800")
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	w.Header().Set("Content-type", "application/json")
}
