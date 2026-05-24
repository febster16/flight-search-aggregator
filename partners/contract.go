package partners

import "context"

type FlightProvider interface {
	Search(context.Context, *SearchRequest) (*SearchResponse, error)
}
