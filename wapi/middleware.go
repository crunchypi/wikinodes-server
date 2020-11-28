package wapi

import (
	"net/http"
	"strings"
	"wikinodes-server/wapi/dosguard"
)

func midDOS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// # Only identify xx.xx.xx.xx:xxxx format (local can be [::1]:xxxx)
		ipport := strings.Split(r.RemoteAddr, ":")
		// # Guard indexing crash.
		if len(ipport) < 1 {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		ip := ipport[0]
		// # Control is a package var in doscheck which handles ip registration.
		// # Guard too many access attempts.
		if ok := dosguard.Control.RegisterCheck(ip); !ok {
			w.WriteHeader(http.StatusTooManyRequests)
			return
		}
		next.ServeHTTP(w, r)
	})
}
