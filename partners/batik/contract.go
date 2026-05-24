package batik

import "flight-search-aggregator/partners"

type batikProvider struct{}

func NewFlightProvider() partners.FlightProvider {
	return &batikProvider{}
}
