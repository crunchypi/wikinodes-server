package wapi

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

// setRoutes sets up routes for this API.
func (h *handler) setRoutes() {
	// http.Handle("/stuff", midDOS(http.HandlerFunc(h.stuff)))
	http.Handle("/data/read", http.HandlerFunc(h.readData))
	http.Handle("/data/read-brief", http.HandlerFunc(h.readDataBrief))

}

// readData handles requests of wikipedia data from the db..
func (h *handler) readData(w http.ResponseWriter, r *http.Request) {
	body, _ := ioutil.ReadAll(r.Body) // # middleware err handled.

	// # Try db retrieve.
	res, err := h.db.NeighboursOfNode(string(body))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	// # Struct -> JSON.
	b, err := json.Marshal(res)
	if err != nil {
		msg := fmt.Sprintf("db err on url: %v", r.URL)
		panic(msg)
	}

	// # Respond.
	w.WriteHeader(http.StatusOK)
	w.Write(b)
}

// readData handles requests of _brief_ wikipedia data from the db.
func (h *handler) readDataBrief(w http.ResponseWriter, r *http.Request) {
	body, _ := ioutil.ReadAll(r.Body) // # middleware err handled.

	// # Try db retrieve.
	res, err := h.db.NeighboursOfNodeBrief(string(body))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	// # Struct -> JSON.
	b, err := json.Marshal(res)
	if err != nil {
		msg := fmt.Sprintf("db err on url: %v", r.URL)
		panic(msg)
	}

	// # Respond.
	w.WriteHeader(http.StatusOK)
	w.Write(b)
}
