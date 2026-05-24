package application

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"flight-search-aggregator/flight/domain"
	"flight-search-aggregator/flight/request"
	"flight-search-aggregator/flight/response"
	"flight-search-aggregator/partners"
	"flight-search-aggregator/utilities"
)

type flight struct {
	// TODO: Consider using map[string]partners.FlightProvider for named provider access.
	// This would allow per-provider retry policies, circuit breakers, and runtime toggling
	// without restarting the service.
	providers []partners.FlightProvider
}

func (f *flight) Search(
	ctx context.Context,
	req *request.Search,
) (*response.Search, error) {
	searchReq := &partners.SearchRequest{
		Origin:        req.Origin,
		Destination:   req.Destination,
		DepartureDate: req.DepartureDate,
		ReturnDate:    req.ReturnDate,
		Passengers:    req.Passengers,
		CabinClass:    req.CabinClass,
	}

	var (
		mu         sync.Mutex
		wg         sync.WaitGroup
		allFlights []domain.Flight
	)

	// TODO: introduce a bounded worker pool for concurrency control
	// under high traffic.
	for _, provider := range f.providers {
		wg.Add(1)
		go searchFlights(ctx, searchReq, provider, &wg, &mu, &allFlights)
	}

	wg.Wait()

	if len(allFlights) == 0 {
		return &response.Search{Flights: []response.Flight{}}, nil
	}

	return mapToSearchResponse(allFlights), nil
}

func searchFlights(
	ctx context.Context,
	req *partners.SearchRequest,
	provider partners.FlightProvider,
	wg *sync.WaitGroup,
	mu *sync.Mutex,
	allFlights *[]domain.Flight,
) {
	defer wg.Done()

	var res *partners.SearchResponse

	if err := utilities.Retry(
		ctx, 3, 15*time.Millisecond, func(ctx context.Context) error {
			var err error

			if res, err = provider.Search(ctx, req); err != nil {
				return fmt.Errorf("error provider search: %w", err)
			}

			return nil
		},
	); err != nil {
		// Log and continue — don't fail the entire search if one provider errors.
		// TODO: probably would be better to send to Slack or structured logging.
		log.Printf("provider search for {%s} failed: %v", provider.Name(), err)

		return
	}

	mu.Lock()

	*allFlights = append(*allFlights, res.Flights...)

	mu.Unlock()
}
