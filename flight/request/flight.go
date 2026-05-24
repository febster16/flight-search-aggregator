package request

type Search struct {
	// Core search criteria
	Origin        string `json:"origin" validate:"required"`
	Destination   string `json:"destination" validate:"required"`
	DepartureDate string `json:"departure_date" validate:"required"`
	ReturnDate    string `json:"return_date"`
	Passengers    int    `json:"passengers" validate:"required,min=1"`
	CabinClass    string `json:"cabin_class" validate:"required"`

	// Additional filters
	MinPrice    *int     `json:"min_price"`
	MaxPrice    *int     `json:"max_price"`
	MaxStops    *int     `json:"max_stops"`
	Airlines    []string `json:"airlines"`
	MaxDuration *int     `json:"max_duration"` // in minutes

	// Possible values: "price_asc", "price_desc", "duration_asc", "duration_desc",
	// "departure_asc", "departure_desc", "arrival_asc", "arrival_desc"
	SortBy string `json:"sort_by"`
}
