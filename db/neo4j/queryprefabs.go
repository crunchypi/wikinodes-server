package neo4j

import (
	"github.com/neo4j/neo4j-go-driver/neo4j"
	"wikinodes-server/db"
)

// This file contains exported funcs (the API)
// of this pkg -- they communicate with the db.
// The purpose is to satisfy the db.DBManager
// behaviour in db/protocols.go.
//
//
//

// ---------------- utils ------------------ //

// simple 'contains' func, because yay go!
func contains(s string, others []string) bool {
	for i := 0; i < len(others); i++ {
		if s == others[i] {
			return true
		}
	}
	return false
}

// ----------------------------------------- //

// ------------ prefabs below -------------- //

// SearchArticlesByID will search through articles
// by their IDs and return all matches.
func (n *Neo4jManager) SearchArticlesByID(id int64) ([]*db.WikiData, error,
) {
	res := make([]*db.WikiData, 0, 1) // # 1 is logically expected.
	cql := `
		MATCH (v:WikiNode)
		WHERE id(v) = $id
		RETURN id(v) as i, v.title as t 
	`
	err := n.execute(executeParams{
		cypher:   cql,
		bindings: map[string]interface{}{"id": id},
		callback: func(r neo4j.Result) {
			v, ok := n.unpackWikiData(r, "i", "t")
			if ok {
				res = append(res, v)
			}
		},
	})
	return res, err
}

// SearchArticlesByTitle will search through articles
// by their title and return all matches.
func (n *Neo4jManager) SearchArticlesByTitle(title string) ([]*db.WikiData, error,
) {
	res := make([]*db.WikiData, 0, 5) // # 5 is arbitrary.
	cql := `
		MATCH (v:WikiNode {title:$title}) RETURN id(v) as i, v.title as t 
	`
	err := n.execute(executeParams{
		cypher:   cql,
		bindings: map[string]interface{}{"title": title},
		callback: func(r neo4j.Result) {
			v, ok := n.unpackWikiData(r, "i", "t")
			if ok {
				res = append(res, v)
			}
		},
	})
	return res, err
}

// SearchArticlesByContent will do a full-text search through
// the database for content that contains the specified string.
// This will be a search on an index named 'ArticleContantIndex'
// so thah must be enabled with the following indexing:
// 	CALL db.index.fulltext.createNodeIndex(
//		"ArticleContentIndex",["WikiNode"],["content"])
func (n *Neo4jManager) SearchArticlesByContent(
	str string, limit int) ([]*db.WikiData, error,
) {
	res := make([]*db.WikiData, 0, limit)
	cql := `
		CALL db.index.fulltext.queryNodes(
			"ArticleContentIndex", $str
		) YIELD node
		RETURN id(node) as i, node.title as t LIMIT $limit 
	`
	err := n.execute(executeParams{
		cypher:   cql,
		bindings: map[string]interface{}{"str": str, "limit": limit},
		callback: func(r neo4j.Result) {
			v, ok := n.unpackWikiData(r, "i", "t")
			if ok {
				res = append(res, v)
			}
		},
	})
	return res, err
}

// SearchArticlesNeightsByIDs will search for article 'A'
// by its ID and return articles that were linked from 'A'.
func (n *Neo4jManager) SearchArticlesNeighsByID(
	id int64, limit int) ([]*db.WikiData, error,
) {
	res := make([]*db.WikiData, 0, 10) // # 10 is arbitrary
	cql := `
		MATCH (v:WikiNode)-[:links]->(w:WikiNode)
		WHERE id(v) = $id
		RETURN id(w) as i, w.title as t
		LIMIT $limit
	`
	err := n.execute(executeParams{
		cypher:   cql,
		bindings: map[string]interface{}{"id": id, "limit": limit},
		callback: func(r neo4j.Result) {
			v, ok := n.unpackWikiData(r, "i", "t")
			if ok {
				res = append(res, v)
			}
		},
	})
	return res, err
}

// SearchArticlesHTMLByID will get the HTML from an article
// with the specified ID.
func (n *Neo4jManager) SearchArticlesHTMLByID(id int64) (string, error) {
	res := ""
	cql := `
		MATCH (v:WikiNode)
		WHERE id(v) = $id
		RETURN v.html as html 
	`
	err := n.execute(executeParams{
		cypher:   cql,
		bindings: map[string]interface{}{"id": id},
		callback: func(r neo4j.Result) {
			v, ok := n.unpackString(r, "html")
			if ok {
				res = v
			}
		},
	})
	return res, err

}

// CheckRelsExistsByIDs will check if there is a relationship
// between articles, i.e if one links another. The Expected
// argument should be an slice containing another two-element
// array where index 0 should have a 'from' id and index 0
// should have a 'to' id. For example if there are two
// articles, where the first one has id 1 and links to
// another article with id 2, then the query [[1,2]] will
// field [true]. If there was no relationship, then the
// result is [false]. This applies to all nested slices.
func (n *Neo4jManager) CheckRelsExistsByIDs(relIDs [][2]int64) ([]bool, error,
) {
	res := make([]bool, len(relIDs))
	cql := `
		MATCH (v:WikiNode)-[:links]->(w:WikiNode)
		WHERE id(v) = $vID
		  AND id(w) = $wID
	   RETURN v.title, w.title 
	`
	for i := 0; i < len(relIDs); i++ {
		err := n.execute(executeParams{
			cypher: cql,
			bindings: map[string]interface{}{
				"vID": relIDs[i][0], "wID": relIDs[i][1]},
			callback: func(r neo4j.Result) {
				if len(r.Record().Values()) > 0 {
					res[i] = true
				}
			},
		})
		if err != nil {
			return res, err
		}
	}
	return res, nil
}

// RandomArticles will return a specified amount of
// randomly picked articles.
func (n *Neo4jManager) RandomArticles(amount int) ([]*db.WikiData, error,
) {
	res := make([]*db.WikiData, 0, amount)
	cql := `
		MATCH (v:WikiNode)
		 WITH COUNT(v) as c
		MATCH (v:WikiNode)
		WHERE rand() < 1/c+0.1 // Decreasing chance per node + a bias.
	   RETURN id(v) as i, v.title as t LIMIT $amount  
	`
	err := n.execute(executeParams{
		cypher:   cql,
		bindings: map[string]interface{}{"amount": amount},
		callback: func(r neo4j.Result) {
			v, ok := n.unpackWikiData(r, "i", "t")
			if ok {
				res = append(res, v)
			}
		},
	})
	return res, err
}
