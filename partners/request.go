package partners

// SearchRequest represents the search parameters passed to each flight provider.
type SearchRequest struct {
	Origin        string
	Destination   string
	DepartureDate string
	ReturnDate    string
	Passengers    int
	CabinClass    string
}
