package response

type Search struct {
	Flights []Flight `json:"flights"`
}

type Flight struct {
	ID             string    `json:"id"`
	Provider       string    `json:"provider"`
	Airline        Airline   `json:"airline"`
	FlightNumber   string    `json:"flight_number"`
	Departure      Departure `json:"departure"`
	Arrival        Arrival   `json:"arrival"`
	Duration       Duration  `json:"duration"`
	Stops          int       `json:"stops"`
	Price          Price     `json:"price"`
	AvailableSeats int       `json:"available_seats"`
	CabinClass     string    `json:"cabin_class"`
	Aircraft       *string   `json:"aircraft,omitempty"`
	Amenities      []string  `json:"amenities,omitempty"`
	Baggage        Baggage   `json:"baggage"`
}

type Airline struct {
	Name string `json:"name"`
	Code string `json:"code"`
}

type Departure struct {
	Airport   string `json:"airport"`
	City      string `json:"city"`
	Datetime  string `json:"datetime"`
	Timestamp int64  `json:"timestamp"`
}

type Arrival struct {
	Airport   string `json:"airport"`
	City      string `json:"city"`
	Datetime  string `json:"datetime"`
	Timestamp int64  `json:"timestamp"`
}

type Duration struct {
	TotalMinutes int    `json:"total_minutes"`
	Formatted    string `json:"formatted"`
}

type Price struct {
	Amount   int    `json:"amount"`
	Currency string `json:"currency"`
}

type Baggage struct {
	CarryOn string `json:"carry_on"`
	Checked string `json:"checked"`
}
