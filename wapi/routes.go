package wapi

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

// setRoutes sets up routes for this API.
func (h *handler) setRoutes() {
	http.Handle("/data/search/node", midDOS(http.HandlerFunc(h.searchNode)))
	http.Handle("/data/search/neigh", midDOS(http.HandlerFunc(h.searchNeigh)))
	http.Handle("/data/search/rand", midDOS(http.HandlerFunc(h.searchRand)))
}

// trySendWikiDataAny takes any <data>, then tries to marshal- and
// send it to a client. Here, this is meant to send any WikiData.
func trySendWikiDataAny(w http.ResponseWriter, data interface{}) {
	b, err := json.Marshal(data)
	if err != nil {
		msg := fmt.Sprintf("db err on url: %v", err)
		log.Println(msg) // @ TODO: Logfile.
		w.WriteHeader(http.StatusInternalServerError)
	}

	// # Respond.
	w.WriteHeader(http.StatusOK)
	w.Write(b)
}

// searchNode forwards request from client->db to find
// a node with a certain entry. Uses:
// 	- db.DBManager.SearchNode(...) ...
//  - db.DBManager.SearchNodeBrief(...) ...
//
// Search options are specified in JSON found at:
// 	- wapi/json.go (jsonOptSearchNode)
//
// testjson/ -d option in curl:
// 		"{\"title\":\"Art\",\"brief\":true}"
func (h *handler) searchNode(w http.ResponseWriter, r *http.Request) {
	body, _ := ioutil.ReadAll(r.Body) // @ TODO: err.

	// # Extract options/config.
	opt := jsonOptSearchNode{}
	if err := json.Unmarshal(body, &opt); err != nil {
		w.WriteHeader(http.StatusNotAcceptable)
		fmt.Println(err) // @ TODO: Logfile.
		return
	}

	// # Result & err, independent on options.
	var resp interface{}
	var err error

	// # Handle options & get db data.
	switch opt.Brief {
	case true:
		resp, err = h.db.SearchNodeBrief(opt.Title)
	case false:
		resp, err = h.db.SearchNode(opt.Title)
	}

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	// # Handle err & try respond.
	trySendWikiDataAny(w, resp)
}

// searchNeigh forwards request from client->db to find
// a neighbours of a db node. Uses:
//  - db.DBManager.SearchNodeNeighBrief(...) ...
//
// Search options are specified in JSON found at:
// 	- wapi/json.go (jsonOptSearchNeigh)
//
// testjson/ -d option in curl:
// 	"{\"title\":\"Art\",\"exclude\":[\"Art history\"], \"limit\":2}"
func (h *handler) searchNeigh(w http.ResponseWriter, r *http.Request) {
	body, _ := ioutil.ReadAll(r.Body) // @ TODO: err.

	// # Extract options/config.
	opt := jsonOptSearchNeigh{}
	if err := json.Unmarshal(body, &opt); err != nil {
		w.WriteHeader(http.StatusNotAcceptable)
		fmt.Println(err) // @ TODO: Logfile.
		return
	}

	// # Get db data.
	resp, err := h.db.SearchNodeNeighBrief(
		opt.Title, opt.Exclude, opt.Limit)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	// # Handle err & try respond.
	trySendWikiDataAny(w, resp)
}

// searchRand forwards request from client->db to find
// random db node(s). Uses:
//  - db.DBManager.RandomNodesBrief(...) ...
//
// Search options are specified in JSON found at:
// 	- wapi/json.go (jsonOptRandomNode)
//
// testjson/ -d option in curl:
//	"{\"amount\":2}"
func (h *handler) searchRand(w http.ResponseWriter, r *http.Request) {
	body, _ := ioutil.ReadAll(r.Body) // @ TODO: err.

	// # Extract options/config.
	opt := jsonOptRandomNode{}
	if err := json.Unmarshal(body, &opt); err != nil {
		w.WriteHeader(http.StatusNotAcceptable)
		fmt.Println(err) // @ TODO: Logfile.
		return
	}
	resp, err := h.db.RandomNodesBrief(opt.Amount)

	// # Handle err & try respond.
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	trySendWikiDataAny(w, resp)
}
