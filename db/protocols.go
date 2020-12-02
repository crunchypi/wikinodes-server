package db

// DBManager specifies interface for necessary work, which is
// defined by sub-packages (concrete db impl), and used by
// the server/api to forward data.
type DBManager interface {
	SearchNode(title string) ([]*WikiData, error)
	SearchNodeBrief(title string) ([]*WikiDataBrief, error)
	// ---
	// # Only Brief allowed, as 'full' is very lkely to be
	// # too much for the system to handle at this moment (201202)
	SearchNodeNeighBrief(title string, exclude []string, limit int) (
		[]*WikiDataBrief, error)
	// ---
	// # Only Brief allowd, pretty much the same reason as above.
	RandomNodesBrief(amount int) ([]*WikiDataBrief, error)
}

// WikiData should represent a database obj/entry for a
// Wikipedia article.
type WikiData struct {
	ID    int64
	Title string
	HTML  string
}

// WikiDataBrief should represent a database obj/entry for a
// Wikipedia article, but only the ID and title.
type WikiDataBrief struct {
	ID    int64
	Title string
}
