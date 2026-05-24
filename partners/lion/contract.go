package lion

import "flight-search-aggregator/partners"

type lionProvider struct{}

func NewFlightProvider() partners.FlightProvider {
	return &lionProvider{}
}
