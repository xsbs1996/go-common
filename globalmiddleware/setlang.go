package globalmiddleware

import "net/http"

// SetLang 设置语言
func SetLang(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		lang := r.Header.Get("lang")
		if len(lang) == 0 {
			r.Header.Set("lang", "en")
		}
		next(w, r)
	}

}
