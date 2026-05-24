package airasia

import (
	"context"
	"encoding/json"
	"flight-search-aggregator/partners"
	"fmt"
	"math/rand"
	"os"
	"time"
)

func (*airasiaProvider) Name() string {
	return partners.ProviderAirAsia
}

func (*airasiaProvider) Search(
	ctx context.Context,
	req *partners.SearchRequest,
) (
	*partners.SearchResponse, error,
) {
	// Simulate 90% success rate
	if rand.Intn(100) < 10 {
		return nil, fmt.Errorf("airasia temporarily unavailable")
	}

	// Simulate 50-150ms delay
	delayMs := rand.Intn(150-50+1) + 50
	time.Sleep(time.Duration(delayMs) * time.Millisecond)

	file, err := os.Open("partners/airasia/airasia_search_response.json")
	if err != nil {
		return nil, fmt.Errorf("error opening airasia response file: %w", err)
	}

	defer file.Close()

	var res Response
	if err = json.NewDecoder(file).Decode(&res); err != nil {
		return nil, fmt.Errorf("error decode airasia response: %w", err)
	}

	searchResponse, err := res.mapToSearchResponse()
	if err != nil {
		return nil, fmt.Errorf("error mapping airasia response: %w", err)
	}

	return searchResponse, nil
}
