package request

// TODO: Filter by: price range, number of stops, departure/arrival time, airlines, duration
// Sort by: price (lowest/highest), duration (shortest/longest), departure time, arrival time
type Search struct {
	Origin        string `json:"origin" validate:"required"`
	Destination   string `json:"destination" validate:"required"`
	DepartureDate string `json:"departure_date" validate:"required"`
	ReturnDate    string `json:"return_date"`
	Passengers    int    `json:"passengers" validate:"required,min=1"`
	CabinClass    string `json:"cabin_class" validate:"required"`
}
