package db

// StoredWikiManager specifies interface for interacting with
// a DB which keeps wikipedia articles.
type StoredWikiManager interface {
	// SearchArticlesByID will search through articles by
	// their IDs and return all matches.
	SearchArticlesByID(id int64) ([]*WikiData, error)
	// SearchArticlesByTitle will search through articles
	// by their title and return all matches.
	SearchArticlesByTitle(title string) ([]*WikiData, error)

	// SearchArticlesByContent will do a full-text search through
	// the database for content that contains the specified string.
	// This will be a search on an index named 'ArticleContantIndex'
	// so thah must be enabled with the following indexing:
	// 	CALL db.index.fulltext.createNodeIndex(
	//		"ArticleContentIndex",["WikiNode"],["content"])
	SearchArticlesByContent(str string, limit int) ([]*WikiData, error)
	// SearchArticlesNeightsByIDs will search for article 'A'
	// by its ID and return articles that were linked from 'A'.
	SearchArticlesNeighsByID(id int64, limit int) ([]*WikiData, error)
	// SearchArticlesHTMLByID will get the HTML from an article
	// with the specified ID.
	SearchArticlesHTMLByID(id int64) (string, error)

	// CheckRelsExistsByIDs will check if there is a relationship
	// between articles, i.e if one links another. The Expected
	// argument should be an slice containing another two-element
	// array where index 0 should have a 'from' id and index 0
	// should have a 'to' id. For example if there are two
	// articles, where the first one has id 1 and links to
	// another article with id 2, then the query [[1,2]] will
	// field [true]. If there was no relationship, then the
	// result is [false]. This applies to all nested slices.
	CheckRelsExistByIDs(relIDs [][2]int64) ([]bool, error)

	// RandomArticles will return a specified amount of
	// randomly picked articles.
	RandomArticles(amount int) ([]*WikiData, error)

	// IncrementRel increments the relationship between two nodes with
	// the given IDs. The incremented relationship is of type HYPERLINKS,
	// where property is 'lookups'. This method is intended to be used
	// for increments such for the purpose of treating the graph as a
	// markov-chain (for article recommendation). Note, the 'lookups'
	// property does not need to exist before using this method.
	IncrementRel(vID, wID int64) error
}

// CacheManager specifies interface for using a cache
// for service improvements.
type CacheManager interface {
	// SetLastQueryID tries to set a query id for an ip. Intenden
	// to be used for keeping track of which Wikipedia Articles a
	// front-end client searches for, for the purpose of article
	// recommendation.
	SetLastQueryID(ip string, id int64) bool
	// LastQueryID is the counterpart of SetLastQueryID, it simply
	// tries to retrieve a Wikipedia Article for a given IP.
	LastQueryID(ip string) (int64, bool)

	// Used to prevent service spam. Calling this method will
	// increment the counter for an IP and check if it has
	// exceeded an allowance over a time period (see pkg vars
	// dosguardAllowance & dosguardExpiration). If True is
	// returned, then the IP is good for more requests.
	CheckRegDOSIP(ip string) (bool, error)
}
