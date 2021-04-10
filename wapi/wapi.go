package wapi

import (
	"net/http"
	"strings"
	"wikinodes-server/config"
	"wikinodes-server/db"
)

var (
	ip   = config.WAPIIP
	port = config.WAPIPort
	// # Server IO time limitation.
	readTimeout  = config.ReadTimeout
	writeTimeout = config.WriteTimeout

	pathToReactApp = config.PathToReactApp
)

// handler serves as a bridge between the app and
// other packages, mainly db.
type handler struct {
	db    db.StoredWikiManager
	cache db.CacheManager
}

// Start starts the app.
func Start(db db.StoredWikiManager, cache db.CacheManager) error {
	// # Enable interface to other ports of this api.
	handler := handler{db: db, cache: cache}
	handler.setRoutes()

	// # Server configs.
	server := http.Server{
		Addr:         ip + ":" + port,
		ReadTimeout:  readTimeout,
		WriteTimeout: writeTimeout,
	}
	return server.ListenAndServe()
}

func extractIP(r *http.Request) (string, bool) {
	// # Only identify xx.xx.xx.xx:xxxx format (local can be [::1]:xxxx)
	ipport := strings.Split(r.RemoteAddr, ":")
	// # Guard indexing crash.
	if len(ipport) < 1 {
		return "", false
	}
	return ipport[0], true
}
