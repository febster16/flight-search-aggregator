package garuda

import "flight-search-aggregator/partners"

type garudaProvider struct{}

func NewFlightProvider() partners.FlightProvider {
	return &garudaProvider{}
}
