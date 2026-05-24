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
	start := time.Now()

	// TODO: Implement cache or single flight based on request filter.
	searchReq := &partners.SearchRequest{
		Origin:        req.Origin,
		Destination:   req.Destination,
		DepartureDate: req.DepartureDate,
		ReturnDate:    req.ReturnDate,
		Passengers:    req.Passengers,
		CabinClass:    req.CabinClass,
	}

	var (
		mu      sync.Mutex
		wg      sync.WaitGroup
		results []providerResult
	)

	// TODO: introduce a bounded worker pool for concurrency control
	// under high traffic.
	for _, provider := range f.providers {
		wg.Add(1)
		go func(p partners.FlightProvider) {
			defer wg.Done()

			result := searchProvider(ctx, searchReq, p)

			mu.Lock()
			results = append(results, result)
			mu.Unlock()
		}(provider)
	}

	wg.Wait()

	var (
		allFlights         []domain.Flight
		providersSucceeded int
		providersFailed    int
	)

	for _, r := range results {
		if r.success {
			providersSucceeded++
			allFlights = append(allFlights, r.flights...)
		} else {
			providersFailed++
		}
	}

	filtered := filterFlights(allFlights, searchReq)

	flights := mapFlightsToResponse(filtered)

	return &response.Search{
		SearchCriteria: mapSearchCriteriaResponse(req),
		Metadata: response.Metadata{
			TotalResults:       len(flights),
			ProvidersQueried:   len(f.providers),
			ProvidersSucceeded: providersSucceeded,
			ProvidersFailed:    providersFailed,
			SearchTimeMs:       time.Since(start).Milliseconds(),
			CacheHit:           false, // TODO: Implement cache.
		},
		Flights: flights,
	}, nil
}

type providerResult struct {
	flights []domain.Flight
	success bool
}

func searchProvider(
	ctx context.Context,
	req *partners.SearchRequest,
	provider partners.FlightProvider,
) providerResult {
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

		return providerResult{success: false}
	}

	return providerResult{
		flights: res.Flights,
		success: true,
	}
}
