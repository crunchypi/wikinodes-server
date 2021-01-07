package wapi

// # This file contains types which reflect
// # some json structures used in this wapi
// # as options in requests.

// options for:
// 	- db.DBManager.SearchNodes(...) ...
// 	- db.DBManager.SearchNodesBrief(...) ...
type jsonOptSearchNode struct {
	Title string `json:"title"`
	Brief bool   `json:"brief"`
}

// options for:
// 	- db.DBManager.SearchNodeNeighBrief(...) ...
type jsonOptSearchNeigh struct {
	Title   string   `json:"title"`
	Exclude []string `json:"exclude"`
	Limit   int      `json:"limit"`
}

// options for:
// 	- db.DBManager.RandomNodesBrief(...) ...
type jsonOptRandomNode struct {
	Amount int `json:"amount"`
}

type jsonOptCheckRels struct {
	Pairs [][2]string `json:"pairs"`
}
