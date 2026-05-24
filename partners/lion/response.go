package lion

import (
	"fmt"

	"flight-search-aggregator/flight/domain"
	"flight-search-aggregator/partners"
)

type Response struct {
	Success bool `json:"success"`
	Data    Data `json:"data"`
}

type Data struct {
	AvailableFlights []Flight `json:"available_flights"`
}

type Flight struct {
	ID         string    `json:"id"`
	Carrier    Carrier   `json:"carrier"`
	Route      Route     `json:"route"`
	Schedule   Schedule  `json:"schedule"`
	FlightTime int       `json:"flight_time"`
	IsDirect   bool      `json:"is_direct"`
	StopCount  int       `json:"stop_count,omitempty"`
	Layovers   []Layover `json:"layovers,omitempty"`
	Pricing    Pricing   `json:"pricing"`
	SeatsLeft  int       `json:"seats_left"`
	PlaneType  string    `json:"plane_type"`
	Services   Services  `json:"services"`
}

type Carrier struct {
	Name string `json:"name"`
	IATA string `json:"iata"`
}

type Route struct {
	From Airport `json:"from"`
	To   Airport `json:"to"`
}

type Airport struct {
	Code string `json:"code"`
	Name string `json:"name"`
	City string `json:"city"`
}

type Schedule struct {
	Departure         string `json:"departure"`
	DepartureTimezone string `json:"departure_timezone"`
	Arrival           string `json:"arrival"`
	ArrivalTimezone   string `json:"arrival_timezone"`
}

type Layover struct {
	Airport         string `json:"airport"`
	DurationMinutes int    `json:"duration_minutes"`
}

type Pricing struct {
	Total    int    `json:"total"`
	Currency string `json:"currency"`
	FareType string `json:"fare_type"`
}

type Services struct {
	WifiAvailable    bool             `json:"wifi_available"`
	MealsIncluded    bool             `json:"meals_included"`
	BaggageAllowance BaggageAllowance `json:"baggage_allowance"`
}

type BaggageAllowance struct {
	Cabin string `json:"cabin"`
	Hold  string `json:"hold"`
}

func (r Response) mapToSearchResponse() (*partners.SearchResponse, error) {
	var flights []domain.Flight

	for _, f := range r.Data.AvailableFlights {
		depTime, err := partners.ParseTimeInLocation(f.Schedule.Departure, partners.TimeLayoutNoTimezone, f.Schedule.DepartureTimezone)
		if err != nil {
			return nil, fmt.Errorf("error parsing departure time for flight %s: %w", f.ID, err)
		}

		arrTime, err := partners.ParseTimeInLocation(f.Schedule.Arrival, partners.TimeLayoutNoTimezone, f.Schedule.ArrivalTimezone)
		if err != nil {
			return nil, fmt.Errorf("error parsing arrival time for flight %s: %w", f.ID, err)
		}

		stops := f.StopCount
		if f.IsDirect {
			stops = 0
		}

		planeType := f.PlaneType

		flight := domain.NewFlight(
			f.ID,
			partners.ProviderLion,
			domain.NewAirline(f.Carrier.Name, f.Carrier.IATA),
			f.ID,
			domain.NewEndpoint(f.Route.From.Code, f.Route.From.City, depTime),
			domain.NewEndpoint(f.Route.To.Code, f.Route.To.City, arrTime),
			domain.NewDuration(f.FlightTime),
			stops,
			domain.NewPrice(f.Pricing.Total, f.Pricing.Currency),
			f.SeatsLeft,
			f.Pricing.FareType,
			&planeType,
			buildAmenities(f.Services),
			domain.NewBaggage(
				f.Services.BaggageAllowance.Cabin,
				f.Services.BaggageAllowance.Hold,
			),
		)

		flights = append(flights, flight)
	}

	return &partners.SearchResponse{
		Flights: flights,
	}, nil
}

func buildAmenities(services Services) []string {
	var amenities []string
	if services.WifiAvailable {
		amenities = append(amenities, "wifi")
	}
	if services.MealsIncluded {
		amenities = append(amenities, "meal")
	}
	return amenities
}
