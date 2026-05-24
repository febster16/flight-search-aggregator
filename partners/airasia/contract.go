package airasia

import "flight-search-aggregator/partners"

type airasiaProvider struct{}

func NewFlightProvider() partners.FlightProvider {
	return &airasiaProvider{}
}
