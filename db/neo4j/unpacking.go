package neo4j

import (
	"wikinodes-server/db"

	"github.com/neo4j/neo4j-go-driver/neo4j"
)

// # This file is used to unpack neo4j data (neo4j.Result)
// # into concrete values, including db.WikiData as well as
// # db.WikiDataBrief.

// # Unpack into int64, using the alias specified in CQL.
func (n *Neo4jManager) unpackInt64(
	r neo4j.Result, alias string) (int64, bool,
) {
	// # Guard value exists.
	v, ok := r.Record().Get(alias)
	if !ok {
		return 0, false
	}
	// # Guard expected val type.
	res, ok := v.(int64)
	return res, ok
}

// # Unpack into string, using the alias specified in CQL.
func (n *Neo4jManager) unpackString(
	r neo4j.Result, alias string) (string, bool,
) {
	// # Guard value exists.
	v, ok := r.Record().Get(alias)
	if !ok {
		return "", false
	}
	// # Guard expected val type.
	res, ok := v.(string)
	return res, ok
}

// # Unpack mimimum logical data unit for this pkg,
// # Specifically, int64 and string. Aliases are
// # the string aliases used in the CQL.
func (n *Neo4jManager) unpackMin(
	r neo4j.Result, aliasID, aliasTitle string) (
	int64, string, bool,
) {
	var id int64
	var title string
	var ok bool
	// # Try unpack & set id.
	id, ok = n.unpackInt64(r, aliasID)
	if !ok {
		return id, title, ok
	}

	// # Try unpack & set title.
	title, ok = n.unpackString(r, aliasTitle)
	if !ok {
		return id, title, ok
	}
	return id, title, ok
}

// # Unpack neo4j result into db.WikiDataBrief.
// # Aliases are the string aliases used in the CQL.
func (n *Neo4jManager) unpackWikiDataBrief(
	r neo4j.Result, aliasID, aliasTitle string) (
	*db.WikiDataBrief, bool,
) {

	res := db.WikiDataBrief{}
	id, title, ok := n.unpackMin(r, aliasID, aliasTitle)
	if !ok {
		return &res, ok
	}
	res.ID, res.Title = id, title

	return &res, true
}

// # Unpack neo4j result into db.WikiData. Aliases
// # are the string aliases used in the CQL.
func (n *Neo4jManager) unpackWikiData(
	r neo4j.Result, aliasID, aliasTitle, aliasHTML string) (
	*db.WikiData, bool,
) {

	res := db.WikiData{}
	id, title, ok := n.unpackMin(r, aliasID, aliasTitle)
	if !ok {
		return &res, ok
	}

	html, ok := n.unpackString(r, aliasHTML)
	if !ok {
		return &res, false
	}

	res.ID, res.Title, res.HTML = id, title, html

	return &res, true
}
