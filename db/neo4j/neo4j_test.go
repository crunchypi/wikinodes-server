package neo4j

import (
	"fmt"
	"testing"
	"wikinodes-server/db"
)

var (
	// # DB credentials.
	uri = "neo4j://localhost:7687"
	usr = ""
	pwd = ""
	// # Global db manager.
	manager db.DBManager = nil

	// # This test will search for a node,
	// # so it is necessary to specify a
	// # known title which is attached to
	// # a db node.
	knownNodeTitle = "Art"
	// # Known neighbours (titles) of above.
	knownNodeNeigh = []string{
		"Abstract animation",
		"Art history",
		"Art manifesto",
		"Art movement",
		"Avant-garde",
	}
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

func TestSearchNodeBrief(t *testing.T) {

	res, err := manager.SearchNodeBrief(knownNodeTitle)
	if err != nil {
		t.Error("Fetch err")
	}

	if len(res) == 0 || res[0].Title != knownNodeTitle {
		t.Error("Empty or incorrect result")
	}
}

func TestSearchNode(t *testing.T) {
	res, err := manager.SearchNode(knownNodeTitle)
	if err != nil {
		t.Error("Fetch err")
	}

	if len(res) == 0 || res[0].Title != knownNodeTitle {
		t.Error("Empty or incorrect result")
	}
}

func TestSearchNodeNeighBrief(t *testing.T) {
	res, err := manager.SearchNodeNeighBrief(
		knownNodeTitle,
		[]string{},
		len(knownNodeNeigh),
	)
	if err != nil {
		t.Error("Fetch err")
	}
	resTitles := make([]string, 0, len(res))
	for _, w := range res {
		resTitles = append(resTitles, w.Title)
	}
	for _, knownNeigh := range knownNodeNeigh {
		if !contains(knownNeigh, resTitles) {
			t.Error(fmt.Sprintf("Did not get %s",
				knownNeigh))
		}
	}
}

func TestRandomNodesBrief(t *testing.T) {
	res1, err1 := manager.RandomNodesBrief(1)
	res2, err2 := manager.RandomNodesBrief(1)
	if err1 != nil || err2 != nil {
		s := fmt.Sprintf("node fetch failed: no1: %s, no2: %s",
			err1, err2)
		t.Error(s)
	}
	if res1[0].Title == res2[0].Title {
		t.Error("rand test: both titles are equal. Try again?")
	}

}

func TestCheckRel(t *testing.T) {
	// # This order should be correct, but naturally depends on
	// # 'knownNodeNeigh' & 'knownNodeTitle'
	q := [][2]string{[2]string{knownNodeNeigh[0], knownNodeTitle}}
	res, err := manager.CheckRel(q)
	if err != nil {
		t.Error("unexpected error")
	}
	if res[0] != true {
		t.Errorf("unexpected result: not neighbours. %v", res)
	}

	// # Test for false positive.

	q = [][2]string{[2]string{"LAHDfL", "DPSBPP"}}
	res, err = manager.CheckRel(q)
	if err != nil {
		t.Error("unexpected error")
	}
	if res[0] != false {
		t.Errorf("unexpected result: not neighbours. %v", res)
	}

}
