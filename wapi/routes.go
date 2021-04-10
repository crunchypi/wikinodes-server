package wapi

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

// setRoutes sets up routes for this API.
func (h *handler) setRoutes() {
	// # Serve static
	http.Handle("/", http.FileServer(http.Dir("./www/build/")))

	routes := map[string]func(w http.ResponseWriter, r *http.Request){
		"/data/search/articles/byid":      h.searchArticlesByID,
		"/data/search/articles/bytitle":   h.searchArticlesByTitle,
		"/data/search/articles/bycontent": h.searchArticlesByContent,
		"/data/search/articles/byneigh":   h.searchArticlesByNeighs,
		"/data/search/html/byid":          h.searchHMLByID,

		"/data/check/relsexist": h.checkRelsExist,
		"/data/random/articles": h.randomArticles,
	}
	for k, v := range routes {
		fmt.Sprintln("Setting route: " + k)
		http.Handle(k, h.midDOS(http.HandlerFunc(v)))
	}
}

// trySendWikiDataAny takes any <data>, then tries to marshal- and
// send it to a client. Here, this is meant to send any WikiData.
func (h *handler) trySendWikiData(
	w http.ResponseWriter, data interface{}, fetcherr error) {
	if fetcherr != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	b, err := json.Marshal(data)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
	// # Respond.
	w.WriteHeader(http.StatusOK)
	w.Write(b)
}

// tryUnpackRequestOptions will try to unmarshal the request body into
// <targetOpt>. If the task fails, then an automatic bad request
// response is sent to the requester and false is returned. Else,
// nothing is written to the requester and the return is true.
func (h *handler) tryUnpackRequestOptions(
	w http.ResponseWriter, r *http.Request, targetOpt interface{}) bool {
	// # Error is not necessary to check, if it's not nil then
	// # the body with JSON request isn't going to work anyway.
	body, _ := ioutil.ReadAll(r.Body)
	err := json.Unmarshal(body, targetOpt)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return false
	}
	return true
}

// searchArticlesByID endpoint accepts a JSON option {id:int}, where the
// id is used to search a database for an article with that id.
// Curl example:
// 	curl http://ip:port/data/search/articles/byid -d "{\"id\":4279}"
func (h *handler) searchArticlesByID(w http.ResponseWriter, r *http.Request) {
	// # Try get JSON options.
	options := struct {
		ID int64 `json:"id"`
	}{}
	if ok := h.tryUnpackRequestOptions(w, r, &options); !ok {
		return
	}
	// # Try db search.
	res, err := h.db.SearchArticlesByID(options.ID)
	// # Try response.
	h.trySendWikiData(w, res, err)
}

// searchArticlesByTitle endpoint accepts a JSON option {title:string}, where the
// title is used to search a database for an article with that title.
// Curl example:
// 	curl http://ip:port/data/search/articles/bytitle -d "{\"title\":\"Art\"}"
func (h *handler) searchArticlesByTitle(w http.ResponseWriter, r *http.Request) {
	// # Try get JSON options.
	options := struct {
		Title string `json:"title"`
	}{}
	if ok := h.tryUnpackRequestOptions(w, r, &options); !ok {
		return
	}
	// # Try db search.
	res, err := h.db.SearchArticlesByTitle(options.Title)
	// # Try response.
	h.trySendWikiData(w, res, err)
}

// searchArticlesByContent endpoint accepts a JOSN option {str:string, limit:int},
// where str is used to search a database for articles that have that string inside
// the bulk content (general search) -- the limit option limits the response.
// Curl example:
// 	curl http://ip:port/data/search/articles/bycontent -d "{\"str\":\"the\", \"limit\":1}"
func (h *handler) searchArticlesByContent(w http.ResponseWriter, r *http.Request) {
	// # Try get JSON options.
	options := struct {
		Str   string `json:"str"`
		Limit int    `json:"limit"`
	}{}
	if ok := h.tryUnpackRequestOptions(w, r, &options); !ok {
		return
	}
	// # Try db search.
	res, err := h.db.SearchArticlesByContent(options.Str, options.Limit)
	// # Try response.
	h.trySendWikiData(w, res, err)
}

// searchArticlesByNeighs endpoint accepts a JSON option {id:int, limit:int}, where
// id searches for a database for an article 'A' with that id, then returns all
// neighbours of 'A' (hyperlinked from 'A') -- the limit option limits the result.
// Curl example:
// 	curl http://ip:port/data/search/articles/byneigh -d "{\"id\":4394, \"limit\":1}"
func (h *handler) searchArticlesByNeighs(w http.ResponseWriter, r *http.Request) {
	// # Try get JSON options.
	options := struct {
		ID    int64 `json:"id"`
		Limit int   `json:"limit"`
	}{}
	if ok := h.tryUnpackRequestOptions(w, r, &options); !ok {
		return
	}
	// # Check what the user searched for last time and use that, if
	// # possible, to increment the relationship between the
	// # last wiki id -> current wiki id. Used for article recommendation.
	if ip, ok := extractIP(r); ok {
		// # Incr the rel if last id is found.
		if lastID, ok := h.cache.LastQueryID(ip); ok {
			h.db.IncrementRel(lastID, options.ID)
		}
		// # Update cache with new id.
		h.cache.SetLastQueryID(ip, options.ID)
	}
	// # Try db search.
	res, err := h.db.SearchArticlesNeighsByID(options.ID, options.Limit)
	// # Try response.
	h.trySendWikiData(w, res, err)
}

// searchHTMLByID endpoint accepts a JSON option {id:int}, where the id
// is used to search a database for an article. Then, the HTML content
// of that article is returned.
// Curl example:
// 	curl http://ip:port/data/search/html/byid -d "{\"id\":4394}"
func (h *handler) searchHMLByID(w http.ResponseWriter, r *http.Request) {
	// # Try get JSON option.
	options := struct {
		ID int64 `json:"id"`
	}{}
	if ok := h.tryUnpackRequestOptions(w, r, &options); !ok {
		return
	}
	// # Try db search.
	res, err := h.db.SearchArticlesHTMLByID(options.ID)
	// # Try response.
	h.trySendWikiData(w, res, err)
}

// checkRelsExist endpoint is used to check if article relationships exist and
// accepts a JSON with the form {rels:[][2]int} . The accepted data is a list
// of lists where index [0] represents a 'from' article id and index [0] represents
// a 'to' article id. In other words, if there are two articles in the database
// where the former has id 8, the latter has id 9, and the former has a link to
// the latter, then the query {rels:[[8,9]]} will return a list with a single true.
// Curl example:
// 	curl http://ip:port/data/check/relsexist -d "{\"rels\":[[4394, 4395]]}"
func (h *handler) checkRelsExist(w http.ResponseWriter, r *http.Request) {
	// # Try get JSON option.
	options := struct {
		Rels [][2]int64 `json:"rels"`
	}{}
	if ok := h.tryUnpackRequestOptions(w, r, &options); !ok {
		return
	}
	// # Try db search.
	res, err := h.db.CheckRelsExistByIDs(options.Rels)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	// # Try response.
	b, err := json.Marshal(res)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write(b)
}

// randomArticles endpoint accepts a JSON with form {limit:int}, where
// the limit specifies how many random articles to return.
// Curl example:
// 	curl http://ip:port/data/random/articles -d "{\"limit\":1}"
func (h *handler) randomArticles(w http.ResponseWriter, r *http.Request) {
	// # Try get JSON option.
	options := struct {
		Limit int `json:"limit"`
	}{}
	if ok := h.tryUnpackRequestOptions(w, r, &options); !ok {
		return
	}
	// # Try db search.
	res, err := h.db.RandomArticles(options.Limit)
	// # Try response.
	h.trySendWikiData(w, res, err)
}
