package garuda

import (
	"context"
	"encoding/json"
	"flight-search-aggregator/partners"
	"fmt"
	"math/rand"
	"os"
	"time"
)

func (*garudaProvider) Search(
	ctx context.Context,
	req *partners.SearchRequest,
) (
	*partners.SearchResponse, error,
) {
	// Simulate 50-100ms delay
	delayMs := rand.Intn(100-50+1) + 50

	time.Sleep(time.Duration(delayMs) * time.Millisecond)

	file, err := os.Open("partners/garuda/garuda_indonesia_search_response.json")
	if err != nil {
		return nil, fmt.Errorf("error opening garuda indonesia response file: %w", err)
	}

	defer file.Close()

	var res Response
	if err = json.NewDecoder(file).Decode(&res); err != nil {
		return nil, fmt.Errorf("error decode garuda indonesia response: %w", err)
	}

	searchResponse, err := res.mapToSearchResponse()
	if err != nil {
		return nil, fmt.Errorf("error mapping garuda indonesia response: %w", err)
	}

	return searchResponse, nil
}
