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

// # ---------------- utils ------------------ //

// simple 'contains' func, because yay go!
func contains(s string, others []string) bool {
	for i := 0; i < len(others); i++ {
		if s == others[i] {
			return true
		}
	}
	return false
}

// # ----------------------------------------- //

// SearchNodeBrief uses a title to search for a node in the db.
// The result is a slice of WikiDataBrief pointers.
func (n *Neo4jManager) SearchNodeBrief(
	title string) (
	[]*db.WikiDataBrief, error,
) {
	res := make([]*db.WikiDataBrief, 0, 10) // # 10 is arbitrary.

	alias := "n"
	aID, aTitle, csv := aliasedPropsBrief(alias)
	node := cqlNode(alias, true)
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

// SearchNode uses a title to search for a node in the db.
// The result is a slice of WikiData pointers.
func (n *Neo4jManager) SearchNode(title string) ([]*db.WikiData, error) {
	res := make([]*db.WikiData, 0, 10) // # 10 is arbitrary.

	/// # Helpers for cql construction and result unwrapping
	alias := "n"
	aID, aTitle, aHTML, csv := aliasedProps(alias)
	node := cqlNode(alias, true)
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

// SearchNodeNeighBrief uses a <title> to search any direct bi-directional
// neighbours of a node with that title. Neihbours with titles in <exclude>
// are ignored and the search is limited to <limit>.
func (n *Neo4jManager) SearchNodeNeighBrief(
	title string, exclude []string, limit int) (
	[]*db.WikiDataBrief, error,
) {
	res := make([]*db.WikiDataBrief, 0, limit)

	// # Helpers for cql construction and result unwrapping
	vAlias, wAlias := "v", "w"
	wID, wTitle, csv := aliasedPropsBrief(wAlias)
	v, w := cqlNode(vAlias, true), cqlNode(wAlias, false)

	cql := fmt.Sprintf(`
		MATCH %s, %s
		WHERE (%s)-[]->(%s)
		   OR (%s)-[]->(%s)
	   RETURN %s LIMIT %d
	`, v, w, vAlias, wAlias, wAlias, vAlias, csv, limit)

	// # Execute cql with bindings, where..
	err := n.execute(executeParams{
		cypher:   cql,
		bindings: map[string]interface{}{nodeSpec.title: title},
		callback: func(r neo4j.Result) {
			// # Try extract data.
			newNode, ok := n.unpackWikiDataBrief(r, wID, wTitle)
			// # Dont add junk.
			if !ok {
				return
			}
			// # Dont add nodes that are not wanted..
			if contains(newNode.Title, exclude) {
				return
			}
			res = append(res, newNode)

		},
	})

	return res, err
}

// RandomNodesBrief returns a specified <amount> of random nodes from the db.
func (n *Neo4jManager) RandomNodesBrief(amount int) ([]*db.WikiDataBrief, error) {
	res := make([]*db.WikiDataBrief, 0, amount)

	// # Helpers for cql construction and result unwrapping.
	alias := "n"
	node := cqlNode(alias, false)
	id, title, csv := aliasedPropsBrief(alias)

	cql := fmt.Sprintf("MATCH %s WHERE rand() < 0.1 RETURN %s LIMIT %d",
		node, csv, amount)
	// # Execute cql with bindings, where..
	err := n.execute(executeParams{
		cypher: cql,
		callback: func(r neo4j.Result) {
			// # Try extract data.
			newNode, ok := n.unpackWikiDataBrief(r, id, title)
			// # Dont add junk.
			if !ok {
				return
			}
			res = append(res, newNode)
		},
	})

	return res, err
}

// CheckRel Checks if relationships exists between tuples
// of titles. Example: If a relationship exists in the db
// such that A-[]->B, and 'rels' = [[A, B]], then the
// return is [true]. Such relationships can be checked
// in bulk, hence slice of slices. Note; this method isn't
// optimised, so it'll make a db request for each tuple.
func (n *Neo4jManager) CheckRel(rels [][2]string) ([]bool, error) {
	res := make([]bool, len(rels))

	// # Construct node for [0] in tuple.
	vPropAlias := "vp" // # For binding.
	v := fmt.Sprintf("(v:%s {%s:$%s})",
		nodeSpec.label, nodeSpec.title, vPropAlias)

	// # Construct node for [1] in tuple.
	wPropAlias := "wp" // # For binding.
	w := fmt.Sprintf("(w:%s {%s:$%s})",
		nodeSpec.label, nodeSpec.title, wPropAlias)

	// # Only interested if match is successful, so 'RETURN'
	// # doesn't have to be much more complex than '1'
	cql := fmt.Sprintf("MATCH %s-[]->%s RETURN 1", v, w)

	// # Check each tuple.
	for i := 0; i < len(rels); i++ {
		// # Execute cql with bindings, where..
		err := n.execute(executeParams{
			cypher: cql,
			bindings: map[string]interface{}{
				vPropAlias: rels[i][0],
				wPropAlias: rels[i][1],
			},
			callback: func(r neo4j.Result) {
				// # Simply flip 'res' bool at 'i' on success.
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
