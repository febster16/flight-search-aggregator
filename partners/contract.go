package partners

import "context"

type FlightProvider interface {
	Name() string
	Search(context.Context, *SearchRequest) (*SearchResponse, error)
}
