package application

import (
	"context"
	"flight-search-aggregator/flight/request"
	"flight-search-aggregator/flight/response"
	"flight-search-aggregator/partners"
	"flight-search-aggregator/partners/airasia"
	"flight-search-aggregator/partners/batik"
	"flight-search-aggregator/partners/garuda"
	"flight-search-aggregator/partners/lion"
)

type Flight interface {
	Search(context.Context, *request.Search) (*response.Search, error)
}

func NewFlightApplication() Flight {
	return &flight{
		providers: initProviders(),
	}
}

func initProviders() []partners.FlightProvider {
	// TODO: Consider using a more dynamic injection for providers,
	// either loading from config (requires re-deployment), or
	// runtime loading without code changes e.g. db.
	return []partners.FlightProvider{
		garuda.NewFlightProvider(),
		lion.NewFlightProvider(),
		batik.NewFlightProvider(),
		airasia.NewFlightProvider(),
	}
}
