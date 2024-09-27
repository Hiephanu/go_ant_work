package middleware

// import (
// 	"fmt"
// 	"net/http"
// 	"time"
// )

// func LoggingMiddleware(next http.Handler) http.Handler {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		start := time.Now()
// 		next.ServeHTTP(w, r)
// 		fmt.Printf("%s %s %s %v\n", r.Method, r.RequestURI, r.RemoteAddr, time.Since(start))
// 	})
// }
