package db

// DBManager specifies interface for necessary work, which is
// defined by sub-packages (concrete db impl), and used by
// the server/api to forward data.
type DBManager interface {
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
	CheckRelsExistsByIDs(relIDs [][2]int64) ([]bool, error)

	// RandomArticles will return a specified amount of
	// randomly picked articles.
	RandomArticles(amount int) ([]*WikiData, error)
}
