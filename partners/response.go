package partners

import "flight-search-aggregator/flight/domain"

// SearchResponse is the unified output returned by each FlightProvider.
type SearchResponse struct {
	Flights []domain.Flight
}
