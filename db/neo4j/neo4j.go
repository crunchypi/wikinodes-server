package neo4j

import (
	"sync"
	"wikinodes-server/db"

	"github.com/neo4j/neo4j-go-driver/neo4j"
)

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
//
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
func New(uri, usr, pwd string) (*Neo4jManager, error) {
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
