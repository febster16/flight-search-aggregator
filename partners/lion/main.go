package lion

import (
	"context"
	"encoding/json"
	"flight-search-aggregator/partners"
	"fmt"
	"math/rand"
	"os"
	"time"
)

func (*lionProvider) Name() string {
	return partners.ProviderLion
}

func (*lionProvider) Search(
	ctx context.Context,
	req *partners.SearchRequest,
) (
	*partners.SearchResponse, error,
) {
	// Simulate 100-200ms delay
	delayMs := rand.Intn(200-100+1) + 100

	time.Sleep(time.Duration(delayMs) * time.Millisecond)

	file, err := os.Open("partners/lion/lion_air_search_response.json")
	if err != nil {
		return nil, fmt.Errorf("error opening lion air response file: %w", err)
	}

	defer file.Close()

	var res Response
	if err = json.NewDecoder(file).Decode(&res); err != nil {
		return nil, fmt.Errorf("error decode lion air response: %w", err)
	}

	searchResponse, err := res.mapToSearchResponse()
	if err != nil {
		return nil, fmt.Errorf("error mapping lion air response: %w", err)
	}

	return searchResponse, nil
}
