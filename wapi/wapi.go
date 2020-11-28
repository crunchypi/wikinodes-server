package wapi

import (
	"net/http"
	"time"
	"wikinodes-server/db"
)

var (
	ip   = "localhost"
	port = "42040"
	// # Server IO time limitation.
	readTimeout  = time.Duration(time.Second * 5)
	writeTimeout = time.Duration(time.Second * 5)
)

// handler serves as a bridge between the app and
// other packages, mainly db.
type handler struct {
	db db.DBManager
}

// Start starts the app.
func Start(db db.DBManager) error {
	// # Enable interface to other ports of this api.
	handler := handler{db: db}
	handler.setRoutes()

	// # Server configs.
	server := http.Server{
		Addr:         ip + ":" + port,
		ReadTimeout:  readTimeout,
		WriteTimeout: writeTimeout,
	}
	return server.ListenAndServe()
}
