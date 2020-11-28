package neo4j

import (
	"fmt"
	"testing"
	"wikinodes-server/db"
)

var (
	// # DB credentials.
	uri = ""
	usr = ""
	pwd = ""
	// # Global db manager.
	manager db.DBManager = nil

	// # This test will search for a node,
	// # so it is necessary to specify a
	// # known title which is attached to
	// # a db node.
	knownNodeTitle = ""
)

func init() {
	// # Check global credentials.
	if uri == "" || usr == "" || pwd == "" {
		panic("uri and credentials not set")
	}

	// # Try set global db manager
	m, err := New(uri, usr, pwd)
	if err != nil {
		msg := fmt.Sprintf("db connection err: %v", err)
		panic(msg)
	}
	manager = m

	// #
	if knownNodeTitle == "" {
		panic("No known searchable node is set.")
	}

}

func TestNeightboursOfNodeBrief(t *testing.T) {

	res, err := manager.NeighboursOfNodeBrief(knownNodeTitle)
	if err != nil {
		t.Error("Fetch err")
	}

	if len(res) == 0 || res[0].Title == "" {
		t.Error("Empty result")
	}
}

func TestNeightboursOfNode(t *testing.T) {
	res, err := manager.NeighboursOfNode(knownNodeTitle)
	if err != nil {
		t.Error("Fetch err")
	}

	if len(res) == 0 || res[0].Title == "" {
		t.Error("Empty result")
	}
}
