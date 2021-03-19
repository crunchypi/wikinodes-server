package wapi

import (
	"net/http"
	"wikinodes-server/wapi/dosguard"
)

func midDOS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip, ok := extractIP(r)
		if !ok {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		// # Control is a package var in doscheck which handles ip registration.
		// # Guard too many access attempts.
		if ok := dosguard.Control.RegisterCheck(ip); !ok {
			w.WriteHeader(http.StatusTooManyRequests)
			return
		}
		next.ServeHTTP(w, r)
	})
}
