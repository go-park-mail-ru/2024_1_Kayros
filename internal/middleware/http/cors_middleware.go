package http

import (
	"net/http"

	"go.uber.org/zap"
)

// Cors решает политику SOP с помощью CORS-заголовков
func Cors(handler http.Handler, logger *zap.Logger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", r.Header.Get("Origin"))
		w.Header().Set("Access-Control-Allow-Methods", "GET,POST,PUT,DELETE")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		handler.ServeHTTP(w, r)
	})
}

//var originAllowlist = []string{
//	"http://127.0.0.1:9999",
//	"http://cats.com",
//	"http://safe.frontend.net",
//}
//
//var methodAllowlist = []string{"GET", "POST", "DELETE", "OPTIONS"}
//
//func checkCORS(next http.Handler) http.Handler {
//	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
//		if isPreflight(r) {
//			origin := r.Header.Get("Origin")
//			method := r.Header.Get("Access-Control-Request-Method")
//			if slices.Contains(originAllowlist, origin) && slices.Contains(methodAllowlist, method) {
//				w.Header().Set("Access-Control-Allow-Origin", origin)
//				w.Header().Set("Access-Control-Allow-Methods", strings.Join(methodAllowlist, ", "))
//				w.Header().Add("Vary", "Origin")
//			}
//		} else {
//			// Not a preflight: regular request.
//			origin := r.Header.Get("Origin")
//			if slices.Contains(originAllowlist, origin) {
//				w.Header().Set("Access-Control-Allow-Origin", origin)
//				w.Header().Add("Vary", "Origin")
//			}
//		}
//		next.ServeHTTP(w, r)
//	})
//}
//
//func isPreflight(r *http.Request) bool {
//	return r.Method == "OPTIONS" &&
//		r.Header.Get("Origin") != "" &&
//		r.Header.Get("Access-Control-Request-Method") != ""
//}
