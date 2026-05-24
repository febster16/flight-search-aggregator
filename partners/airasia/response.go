package airasia

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
	FlightCode    string  `json:"flight_code"`
	Airline       string  `json:"airline"`
	FromAirport   string  `json:"from_airport"`
	ToAirport     string  `json:"to_airport"`
	DepartTime    string  `json:"depart_time"`
	ArriveTime    string  `json:"arrive_time"`
	DurationHours float64 `json:"duration_hours"`
	DirectFlight  bool    `json:"direct_flight"`
	Stops         []Stop  `json:"stops,omitempty"`
	PriceIDR      int     `json:"price_idr"`
	Seats         int     `json:"seats"`
	CabinClass    string  `json:"cabin_class"`
	BaggageNote   string  `json:"baggage_note"`
}

type Stop struct {
	Airport         string `json:"airport"`
	WaitTimeMinutes int    `json:"wait_time_minutes"`
}

func (r Response) mapToSearchResponse() (*partners.SearchResponse, error) {
	var flights []domain.Flight

	for _, f := range r.Flights {
		depTime, err := partners.ParseTime(f.DepartTime, partners.TimeLayoutRFC3339)
		if err != nil {
			return nil, fmt.Errorf("error parsing departure time for flight %s: %w", f.FlightCode, err)
		}

		arrTime, err := partners.ParseTime(f.ArriveTime, partners.TimeLayoutRFC3339)
		if err != nil {
			return nil, fmt.Errorf("error parsing arrival time for flight %s: %w", f.FlightCode, err)
		}

		durationMinutes := partners.HoursToMinutes(f.DurationHours)

		stops := len(f.Stops)
		if f.DirectFlight {
			stops = 0
		}

		flight := domain.NewFlight(
			f.FlightCode,
			partners.ProviderAirAsia,
			domain.NewAirline(f.Airline, partners.ExtractAirlineCode(f.FlightCode)),
			f.FlightCode,
			domain.NewEndpoint(f.FromAirport, partners.AirportCity(f.FromAirport), depTime),
			domain.NewEndpoint(f.ToAirport, partners.AirportCity(f.ToAirport), arrTime),
			domain.NewDuration(durationMinutes),
			stops,
			domain.NewPrice(f.PriceIDR, partners.CurrencyIDR),
			f.Seats,
			partners.NormalizeCabinClass(f.CabinClass),
			nil,
			nil,
			// TODO: The baggage_note field contains free-text
			// e.g. "Cabin baggage only, checked bags additional fee"
			// which doesn't map cleanly to carry_on/checked weight values.
			// Implement a more detailed parser or normalize
			// to a structured format when requirements are clarified.
			domain.NewBaggage(f.BaggageNote, ""),
		)

		flights = append(flights, flight)
	}

	return &partners.SearchResponse{
		Flights: flights,
	}, nil
}
