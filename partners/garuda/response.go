package garuda

import (
	"fmt"

	"flight-search-aggregator/flight/domain"
	"flight-search-aggregator/partners"
)

type Response struct {
	Status  string   `json:"status"`
	Flights []Flight `json:"flights"`
}

type Flight struct {
	FlightID        string    `json:"flight_id"`
	Airline         string    `json:"airline"`
	AirlineCode     string    `json:"airline_code"`
	Departure       Endpoint  `json:"departure"`
	Arrival         Endpoint  `json:"arrival"`
	DurationMinutes int       `json:"duration_minutes"`
	Stops           int       `json:"stops"`
	Aircraft        string    `json:"aircraft"`
	Price           Price     `json:"price"`
	AvailableSeats  int       `json:"available_seats"`
	FareClass       string    `json:"fare_class"`
	Baggage         Baggage   `json:"baggage"`
	Amenities       []string  `json:"amenities,omitempty"`
	Segments        []Segment `json:"segments,omitempty"`
}

type Endpoint struct {
	Airport  string `json:"airport"`
	City     string `json:"city"`
	Time     string `json:"time"`
	Terminal string `json:"terminal"`
}

type Price struct {
	Amount   int    `json:"amount"`
	Currency string `json:"currency"`
}

type Baggage struct {
	CarryOn int `json:"carry_on"`
	Checked int `json:"checked"`
}

type Segment struct {
	FlightNumber    string          `json:"flight_number"`
	Departure       SegmentEndpoint `json:"departure"`
	Arrival         SegmentEndpoint `json:"arrival"`
	DurationMinutes int             `json:"duration_minutes"`
	LayoverMinutes  int             `json:"layover_minutes,omitempty"`
}

type SegmentEndpoint struct {
	Airport string `json:"airport"`
	Time    string `json:"time"`
}

func (r Response) mapToSearchResponse() (*partners.SearchResponse, error) {
	var flights []domain.Flight

	for _, f := range r.Flights {
		depTime, err := partners.ParseTime(f.Departure.Time, partners.TimeLayoutRFC3339)
		if err != nil {
			return nil, fmt.Errorf("error parsing departure time for flight %s: %w", f.FlightID, err)
		}

		arrTime, err := partners.ParseTime(f.Arrival.Time, partners.TimeLayoutRFC3339)
		if err != nil {
			return nil, fmt.Errorf("error parsing arrival time for flight %s: %w", f.FlightID, err)
		}

		aircraft := f.Aircraft

		flight := domain.NewFlight(
			f.FlightID,
			partners.ProviderGaruda,
			domain.NewAirline(f.Airline, f.AirlineCode),
			f.FlightID,
			domain.NewEndpoint(f.Departure.Airport, f.Departure.City, depTime),
			domain.NewEndpoint(f.Arrival.Airport, f.Arrival.City, arrTime),
			domain.NewDuration(f.DurationMinutes),
			f.Stops,
			domain.NewPrice(f.Price.Amount, f.Price.Currency),
			f.AvailableSeats,
			f.FareClass,
			&aircraft,
			f.Amenities,
			domain.NewBaggage(
				fmt.Sprint(f.Baggage.CarryOn),
				fmt.Sprint(f.Baggage.Checked),
			),
		)

		flights = append(flights, flight)
	}

	return &partners.SearchResponse{
		Flights: flights,
	}, nil
}
