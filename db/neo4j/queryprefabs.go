package neo4j

import (
	"fmt"
	"wikinodes-server/db"

	"github.com/neo4j/neo4j-go-driver/neo4j"
)

// # This file contains exported funcs (the API)
// # of this pkg -- they communicate with the db.
// # The purpose is to satisfy the db.DBManager
// # behaviour in db/protocols.go.
// #
// #
// #

// NeighboursOfNodeBrief accepts a <title> which is expected to be
// a property of a node 'V' in Neo4j. Label and other relevant props
// of this 'V' is defined at the top of this file.
// All neighbours of 'V' will be returned as []db.WikiDataBrief and
// appropriate props.
func (n *Neo4jManager) NeighboursOfNodeBrief(
	title string) (
	[]*db.WikiDataBrief, error,
) {
	res := make([]*db.WikiDataBrief, 0, 10) // # 10 is arbitrary.

	alias := "n"
	aID, aTitle, csv := aliasedPropsBrief(alias)
	node := cqlNode(alias)
	cql := fmt.Sprintf("MATCH %s RETURN %s", node, csv)

	// # Execute cql with bindings, where..
	err := n.execute(executeParams{
		cypher:   cql,
		bindings: map[string]interface{}{nodeSpec.title: title},
		callback: func(r neo4j.Result) {
			// # Try extract data.
			newNode, ok := n.unpackWikiDataBrief(r, aID, aTitle)
			if ok {
				res = append(res, newNode)
			}
		},
	})

	return res, err
}

// NeighboursOfNode accepts a <title> which is expected to be
// a property of a node 'V' in Neo4j. Label and other relevant props
// of this 'V' is defined at the top of this file.
// All neighbours of 'V' will be returned as []db.WikiData and
// appropriate props.
func (n *Neo4jManager) NeighboursOfNode(title string) ([]*db.WikiData, error) {

	res := make([]*db.WikiData, 0, 10) // # 10 is arbitrary.
	// # Node aliasing. These Aliases will be used to construct a cypher str
	// # with fmt:
	// #   			 MATCH (alias:aliasLabel {...})
	// #  			RETURN aID, aTitle, aHTML
	alias := "n"
	aID, aTitle, aHTML, csv := aliasedProps(alias)
	node := cqlNode(alias)
	cql := fmt.Sprintf("MATCH %s RETURN %s", node, csv)

	// # Execute cql with bindings, where callback is ..
	err := n.execute(executeParams{
		cypher:   cql,
		bindings: map[string]interface{}{nodeSpec.title: title},
		callback: func(r neo4j.Result) {
			// # .. trying to extract data.
			newNode, ok := n.unpackWikiData(r, aID, aTitle, aHTML)
			if ok {
				res = append(res, newNode)
			}

		},
	})

	return res, err
}
