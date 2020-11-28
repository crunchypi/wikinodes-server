package db

// DBManager specifies interface for necessary work, which is
// defined by sub-packages (concrete db impl), and used by
// the server/api to forward data.
type DBManager interface {
	NeighboursOfNodeBrief(title string) ([]WikiDataBrief, error)
	NeighboursOfNode(title string) ([]WikiData, error)
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
