package wapi

import (
	"net/http"
)

func (h *handler) midDOS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip, ok := extractIP(r)
		if !ok {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		// # Check/Register for the purpose of identifying abuse.
		allow, err := h.cache.CheckRegDOSIP(ip)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		if !allow {
			w.WriteHeader(http.StatusTooManyRequests)
			return
		}

		next.ServeHTTP(w, r)
	})
}
