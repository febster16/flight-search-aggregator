package request

type Search struct {
	Origin        string `json:"origin" validate:"required"`
	Destination   string `json:"destination" validate:"required"`
	DepartureDate string `json:"departureDate" validate:"required"`
	ReturnDate    string `json:"returnDate"`
	Passengers    int    `json:"passengers" validate:"required,min=1"`
	CabinClass    string `json:"cabinClass" validate:"required"`
}
