package batik

import (
	"fmt"
	"regexp"

	"flight-search-aggregator/flight/domain"
	"flight-search-aggregator/partners"
)

type Response struct {
	Code    int      `json:"code"`
	Message string   `json:"message"`
	Results []Flight `json:"results"`
}

type Flight struct {
	FlightNumber      string       `json:"flightNumber"`
	AirlineName       string       `json:"airlineName"`
	AirlineIATA       string       `json:"airlineIATA"`
	Origin            string       `json:"origin"`
	Destination       string       `json:"destination"`
	DepartureDateTime string       `json:"departureDateTime"`
	ArrivalDateTime   string       `json:"arrivalDateTime"`
	TravelTime        string       `json:"travelTime"`
	NumberOfStops     int          `json:"numberOfStops"`
	Connections       []Connection `json:"connections,omitempty"`
	Fare              Fare         `json:"fare"`
	SeatsAvailable    int          `json:"seatsAvailable"`
	AircraftModel     string       `json:"aircraftModel"`
	BaggageInfo       string       `json:"baggageInfo"`
	OnboardServices   []string     `json:"onboardServices,omitempty"`
}

type Connection struct {
	StopAirport  string `json:"stopAirport"`
	StopDuration string `json:"stopDuration"`
}

type Fare struct {
	BasePrice    int    `json:"basePrice"`
	Taxes        int    `json:"taxes"`
	TotalPrice   int    `json:"totalPrice"`
	CurrencyCode string `json:"currencyCode"`
	Class        string `json:"class"`
}

func (r Response) mapToSearchResponse() (*partners.SearchResponse, error) {
	var flights []domain.Flight

	for _, f := range r.Results {
		depTime, err := partners.ParseTime(f.DepartureDateTime, partners.TimeLayoutISO8601)
		if err != nil {
			return nil, fmt.Errorf("error parsing departure time for flight %s: %w", f.FlightNumber, err)
		}

		arrTime, err := partners.ParseTime(f.ArrivalDateTime, partners.TimeLayoutISO8601)
		if err != nil {
			return nil, fmt.Errorf("error parsing arrival time for flight %s: %w", f.FlightNumber, err)
		}

		durationMinutes := partners.ParseDurationString(f.TravelTime)
		aircraftModel := f.AircraftModel

		flight := domain.NewFlight(
			f.FlightNumber,
			partners.ProviderBatik,
			domain.NewAirline(f.AirlineName, f.AirlineIATA),
			f.FlightNumber,
			domain.NewEndpoint(f.Origin, partners.AirportCity(f.Origin), depTime),
			domain.NewEndpoint(f.Destination, partners.AirportCity(f.Destination), arrTime),
			domain.NewDuration(durationMinutes),
			f.NumberOfStops,
			domain.NewPrice(f.Fare.TotalPrice, f.Fare.CurrencyCode),
			f.SeatsAvailable,
			partners.NormalizeCabinClass(f.Fare.Class),
			&aircraftModel,
			f.OnboardServices,
			domain.NewBaggage(
				extractBaggagePart(f.BaggageInfo, "cabin"),
				extractBaggagePart(f.BaggageInfo, "checked"),
			),
		)

		flights = append(flights, flight)
	}

	return &partners.SearchResponse{
		Flights: flights,
	}, nil
}

func extractBaggagePart(info string, part string) string {
	re := regexp.MustCompile(`(\d+kg)\s+` + part)

	if matches := re.FindStringSubmatch(info); len(matches) == 2 {
		return matches[1]
	}

	return ""
}
