package neo4j

import (
	"fmt"
	"testing"
)

var (
	// DB credentials.
	uri = "neo4j://localhost:7687"
	usr = ""
	pwd = ""
	// Global db manager.
	n *Neo4jManager = nil
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
	n = m

}

// ------------ a few methods for creating data --------------- //
func (n *Neo4jManager) clear() error {
	return n.execute(executeParams{
		cypher: "MATCH (n:WikiData) DETACH DELETE n",
	})
}

func (n *Neo4jManager) createNode(title, content, html string) error {
	return n.execute(executeParams{
		cypher: "CREATE (:WikiData {title:$title, html:$html, content:$content})",
		bindings: map[string]interface{}{
			"title": title, "html": html, "content": content},
		callback: nil,
	})
}

func (n *Neo4jManager) createNodesAndRel(vTitle, wTitle string) error {
	cql := "MERGE (:WikiData{title:$vTitle})-[:HYPERLINKS]->(:WikiData{title:$wTitle})"
	return n.execute(executeParams{
		cypher:   cql,
		bindings: map[string]interface{}{"vTitle": vTitle, "wTitle": wTitle},
		callback: nil,
	})
}

// ------------------------------------------------------------ //

func TestSearchArticlesByTitle(t *testing.T) {
	n.clear()
	defer n.clear()
	title, content, html := "a", "", ""
	n.createNode(title, content, html)
	data, err := n.SearchArticlesByTitle(title)
	if err != nil {
		t.Fatal(err)
	}
	if len(data) == 0 {
		t.Fatal("empty result")
	}
	if data[0].Title != title {
		t.Fatal("unexpected result: ", data[0].Title)
	}
}

func TestSearchArticlesByContent(t *testing.T) {
	n.clear()
	defer n.clear()

	title, content, html := "a", "b", "c"
	n.createNode(title, content, html)

	res, err := n.SearchArticlesByContent(content, 1)
	if err != nil {
		t.Fatal(err)
	}
	if res[0].Title != title {
		t.Fatal("got incorrect result")
	}

}

func TestSearchArticlesByID(t *testing.T) {
	n.clear()
	defer n.clear()
	title, content, html := "a", "", ""
	n.createNode(title, content, html)

	// # Get and check unsafely for brevity. The
	// # previous test checks this properly.
	data, _ := n.SearchArticlesByTitle(title)
	if data[0].Title != title {
		t.Fatal("unexpected title result: ", data[0].Title)
	}

	res, _ := n.SearchArticlesByID(data[0].ID)
	if data[0].Title != res[0].Title {
		t.Fatal("unexpected id result")
	}
}

func TestSearchArticlesNeighsByID(t *testing.T) {
	n.clear()
	defer n.clear()
	vTitle, wTitle := "v", "w"
	n.createNodesAndRel(vTitle, wTitle)

	data, _ := n.SearchArticlesByTitle(vTitle)
	res, _ := n.SearchArticlesNeighsByID(data[0].ID, 1)
	t.Log("Got here")
	if res[0].Title != wTitle {
		t.Fatal("did not get neighbour")
	}
}

func TestSearchArticlesHTMLByID(t *testing.T) {
	n.clear()
	defer n.clear()

	title, content, html := "v", "", "some content"
	n.createNode(title, content, html)

	data, _ := n.SearchArticlesByTitle(title)
	res, _ := n.SearchArticlesHTMLByID(data[0].ID)
	if res != html {
		t.Fatal("expected html, got: " + res)
	}
}

func TestCheckRelsExistByIDs(t *testing.T) {
	n.clear()
	defer n.clear()

	vTitle, wTitle := "v", "w"
	n.createNode(vTitle, "", "")
	n.createNode(wTitle, "", "")

	// # This section should fail since there are no rels.
	vData, _ := n.SearchArticlesByTitle(vTitle)
	wData, _ := n.SearchArticlesByTitle(wTitle)

	r1, _ := n.CheckRelsExistByIDs([][2]int64{{vData[0].ID, wData[0].ID}})
	if r1[0] == true {
		t.Fatal("rel check should be false")
	}
	n.clear()

	// # This section should _not_ fail since there are rels.
	n.createNodesAndRel(vTitle, wTitle)
	vData, _ = n.SearchArticlesByTitle(vTitle)
	wData, _ = n.SearchArticlesByTitle(wTitle)

	r2, _ := n.CheckRelsExistByIDs([][2]int64{{vData[0].ID, wData[0].ID}})
	if r2[0] == false {
		t.Fatal("rel check should be true")
	}

}

func TestRandomArticles(t *testing.T) {
	n.clear()
	defer n.clear()

	// # Create a number of titles and use them for node creation
	titles := make([]string, 0, 10)
	for i := 0; i < 10; i++ {
		titles = append(titles, fmt.Sprintf("%v", i))
	}
	for _, title := range titles {
		n.createNode(title, "", "")
	}

	// # Check for thrice in a row, should be unlikely
	matches := 0
	for i := 0; i < 3; i++ {
		res, _ := n.RandomArticles(1)
		if res[0].Title == titles[5] { // # 5 is arbitrary.
			matches += 1
		}
	}
	if matches == 3 {
		t.Fatal("unlikely result (same 3 times in a row)")
	}
}

func TestIncrementRelAndSearchArticlesNeighsByID(t *testing.T) {
	n.clear()
	defer n.clear()
	// # rel: q -> a,b,c
	n.execute(executeParams{cypher: `
		CREATE (q:WikiData{title:'q'})-[:HYPERLINKS]
								->(:WikiData{title:'a'})
		CREATE (q)-[:HYPERLINKS]->(:WikiData{title:'b'})
		CREATE (q)-[:HYPERLINKS]->(:WikiData{title:'c'})
	`})
	// # Get ID for all new nodes.
	titles := []string{"q", "a", "b", "c"}
	ids := make([]int64, 4)
	for i, v := range titles {
		r, _ := n.SearchArticlesByTitle(v)
		ids[i] = (*r[0]).ID
	}
	// # Incr rels such that best-fit should be in order: a,b,c
	for i := 0; i < 200; i++ {
		n.IncrementRel(ids[0], ids[1])
	}
	for i := 0; i < 100; i++ {
		n.IncrementRel(ids[0], ids[2])
	}
	// # Get & check data.
	res, _ := n.SearchArticlesNeighsByID(ids[0], 3)
	resTitles := make([]string, 0, 4)
	for _, v := range res {
		resTitles = append(resTitles, (*v).Title)
	}

	if titles[1] != resTitles[0] || titles[2] != resTitles[1] {
		t.Errorf("wanted %v, got %v. might be a rand issue (retry)",
			titles[1:], resTitles)
	}
}
