package batik

import (
	"context"
	"encoding/json"
	"flight-search-aggregator/partners"
	"fmt"
	"math/rand"
	"os"
	"time"
)

func (*batikProvider) Name() string {
	return partners.ProviderBatik
}

func (*batikProvider) Search(
	ctx context.Context,
	req *partners.SearchRequest,
) (
	*partners.SearchResponse, error,
) {
	// Simulate 200-400ms delay
	delayMs := rand.Intn(400-200+1) + 200

	time.Sleep(time.Duration(delayMs) * time.Millisecond)

	file, err := os.Open("partners/batik/batik_air_search_response.json")
	if err != nil {
		return nil, fmt.Errorf("error opening batik air response file: %w", err)
	}

	defer file.Close()

	var res Response
	if err = json.NewDecoder(file).Decode(&res); err != nil {
		return nil, fmt.Errorf("error decode batik air response: %w", err)
	}

	searchResponse, err := res.mapToSearchResponse()
	if err != nil {
		return nil, fmt.Errorf("error mapping batik air response: %w", err)
	}

	return searchResponse, nil
}
