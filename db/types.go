package db

// WikiData represents a packet of 'normal' data
// retrieved from the database for normal front-
// end operations. This does not include the html
// because that must be fetched separately, as
// that is relatively rare.
type WikiData struct {
	ID    int64  `json:"id"`
	Title string `json:"title"`
}
