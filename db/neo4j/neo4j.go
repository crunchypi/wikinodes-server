package neo4j

import (
	"fmt"
	"sync"
	"wikinodes-server/db"

	"github.com/neo4j/neo4j-go-driver/neo4j"
)

// ---------------- helpers -------------------------- //

// # Describes properties of relevant nodes in neo4j.
var (
	nodeLabel     = "WikiData"
	nodePropTitle = "title"
	nodePropHTML  = "html"
)

// # Simply constructs a neo4j property bindings string.
// # Example:
// # 	let propIDs = []string{"a", "b", "c"}
// # 	result = "{a:$a, b:$b, c:$c}"
func nodePropBindStr(propIDs []string) string {

	cql := "{" // Open.

	l := len(propIDs) // Len once.
	for i := 0; i < l; i++ {
		cql += fmt.Sprintf("%s:$%s", propIDs[i], propIDs[i])

		// # Add commas after each prop, except for last
		if i < l-1 {
			cql += ", "
		}
	}

	cql += "}" // Close.
	return cql
}

// ---------------- /helpers -------------------------- //

// Exclusively used for Neo4jManager.execute(). Defined
// as a struct mainly for briefer method signatures.
type executeParams struct {
	cypher   string
	bindings map[string]interface{}
	callback func(neo4j.Result)
}

// Neo4jManager implements db.DBManager interface.
// But seriously, why can't this be done like the
// Rust folks do it, ya know, like sane people?!?
//								   __		   _
// 			i m p l i c i t 		\(︶︹ ︺')/
//
var _ db.DBManager = &Neo4jManager{}

// Neo4jManager -- manages neo4j connection and friends.
type Neo4jManager struct {
	mx sync.Mutex
	db neo4j.Driver
}

// New attempts to return Neo4jManager with an active
// Neo4j driver.
func New(uri, usr, pwd string) (db.DBManager, error) {
	new := Neo4jManager{}

	driver, err := neo4j.NewDriver(
		uri,
		neo4j.BasicAuth(usr, pwd, ""),
		func(c *neo4j.Config) {
			c.Encrypted = false
		},
	)
	if err != nil {
		return &new, err
	}

	new.db = driver
	return &new, nil
}

// General async-safe executor, expects T executeParams
// as arg, see type def in this pkg.
func (n *Neo4jManager) execute(x executeParams) error {
	// # Standard syncing.
	n.mx.Lock()
	defer n.mx.Unlock()

	// # Open.
	session, err := n.db.Session(neo4j.AccessModeWrite)
	if err != nil {
		return err
	}
	defer session.Close()

	// # Execute.
	res, err := session.Run(x.cypher, x.bindings)
	if err != nil {
		return err
	}

	// # Optional callback.
	if x.callback != nil {
		for res.Next() {
			x.callback(res)
		}
	}
	return nil
}

// ----------------------------------------------------------
// !!
// !!	Dear reader: I apologise for the code below -- i was
// !! 	in a rush and will fix this at a later date, promise!
// !!
// !!	In the meantime, rest assured that the methods work.
// ----------------------------------------------------------

// NeighboursOfNodeBrief accepts a <title> which is expected to be
// a property of a node 'V' in Neo4j. Label and other relevant props
// of this 'V' is defined at the top of this file.
// All neighbours of 'V' will be returned as []db.WikiDataBrief and
// appropriate props.
func (n *Neo4jManager) NeighboursOfNodeBrief(title string) ([]db.WikiDataBrief, error) {
	res := make([]db.WikiDataBrief, 0, 10) // # 10 is arbitrary.
	// # Node aliasing & property construction.
	alias := "n"
	bindStr := nodePropBindStr([]string{nodePropTitle})

	cql := fmt.Sprintf("MATCH (%s:%s %s) RETURN id(%s), %s.%s",
		alias, nodeLabel, bindStr, alias, alias, nodePropTitle)

	// # Execute cql with bindings, where..
	err := n.execute(executeParams{
		cypher:   cql,
		bindings: map[string]interface{}{nodePropTitle: title},
		callback: func(r neo4j.Result) {

			// # .. expected return is ...
			newNode := db.WikiDataBrief{}

			// # .. with property ...
			search := fmt.Sprintf("%s.%s", alias, nodePropTitle)
			if v, ok := r.Record().Get(search); ok {
				// # Because YOLO.
				newNode.Title = v.(string)
			}
			// # .. and ...
			if v, ok := r.Record().Get("id(n)"); ok {
				// # Because YOLO -- the return.
				newNode.ID = v.(int64)
			}

			// # Only add to result if complete:
			if newNode.Title != "" && newNode.ID != 0 {
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
func (n *Neo4jManager) NeighboursOfNode(title string) ([]db.WikiData, error) {

	res := make([]db.WikiData, 0, 10) // # 10 is arbitrary.
	// # Node aliasing & property construction.
	alias := "n"
	bindStr := nodePropBindStr([]string{nodePropTitle})

	cql := fmt.Sprintf("MATCH (%s:%s %s) RETURN %s.%s, id(%s), %s.%s",
		alias, nodeLabel, bindStr, alias,
		nodePropTitle, alias, alias, nodePropHTML)

	// # Execute cql with bindings, where..
	err := n.execute(executeParams{
		cypher:   cql,
		bindings: map[string]interface{}{nodePropTitle: title},
		callback: func(r neo4j.Result) {

			// # .. expected return is ...
			newNode := db.WikiData{}

			// # .. with property ...
			if v, ok := r.Record().Get("id(n)"); ok {
				// # Because YOLO -- the return.
				newNode.ID = v.(int64)
			}
			// # .. and ...
			search := fmt.Sprintf("%s.%s", alias, nodePropTitle)
			if v, ok := r.Record().Get(search); ok {
				// # Because YOLO.
				newNode.Title = v.(string)
			}
			// # .. and ..
			search = fmt.Sprintf("%s.%s", alias, nodePropHTML)
			if v, ok := r.Record().Get(search); ok {
				// # Because YOLO.
				newNode.HTML = v.(string)
			}

			// # Only add to result if complete:
			if newNode.Title != "" && newNode.ID != 0 {
				res = append(res, newNode)
			}
		},
	})

	return res, err
}
