package application

import (
	"context"
	"fmt"
	"log"
	"time"

	"flight-search-aggregator/flight/domain"
	"flight-search-aggregator/flight/request"
	"flight-search-aggregator/flight/response"
	"flight-search-aggregator/partners"
)

type flight struct {
	// TODO: Consider using map[string]partners.FlightProvider for named provider access.
	// This would allow per-provider retry policies, circuit breakers, and runtime toggling
	// without restarting the service.
	providers []partners.FlightProvider
}

func (f *flight) Search(ctx context.Context, req *request.Search) (*response.Search, error) {
	searchReq := &partners.SearchRequest{
		Origin:        req.Origin,
		Destination:   req.Destination,
		DepartureDate: req.DepartureDate,
		ReturnDate:    req.ReturnDate,
		Passengers:    req.Passengers,
		CabinClass:    req.CabinClass,
	}

	var allFlights []domain.Flight

	for _, provider := range f.providers {
		res, err := provider.Search(ctx, searchReq)
		if err != nil {
			// Log and continue — don't fail the entire search if one provider errors.
			log.Printf("[%s] provider search failed: %v", time.Now().Format(time.RFC3339), err)
			continue
		}

		allFlights = append(allFlights, res.Flights...)
	}

	if len(allFlights) == 0 {
		return &response.Search{Flights: []response.Flight{}}, nil
	}

	return mapToSearchResponse(allFlights), nil
}

func mapToSearchResponse(flights []domain.Flight) *response.Search {
	result := &response.Search{
		Flights: make([]response.Flight, 0, len(flights)),
	}

	for _, f := range flights {
		result.Flights = append(result.Flights, response.Flight{
			ID:           f.ID(),
			Provider:     f.Provider(),
			Airline:      response.Airline{Name: f.Airline().Name(), Code: f.Airline().Code()},
			FlightNumber: f.FlightNumber(),
			Departure: response.Endpoint{
				Airport:   f.Departure().Airport(),
				City:      f.Departure().City(),
				Datetime:  formatDatetime(f.Departure().Datetime()),
				Timestamp: f.Departure().Timestamp(),
			},
			Arrival: response.Endpoint{
				Airport:   f.Arrival().Airport(),
				City:      f.Arrival().City(),
				Datetime:  formatDatetime(f.Arrival().Datetime()),
				Timestamp: f.Arrival().Timestamp(),
			},
			Duration: response.Duration{
				TotalMinutes: f.Duration().TotalMinutes(),
				Formatted:    f.Duration().Formatted(),
			},
			Stops: f.Stops(),
			Price: response.Price{
				Amount:   f.Price().Amount(),
				Currency: f.Price().Currency(),
			},
			AvailableSeats: f.AvailableSeats(),
			CabinClass:     f.CabinClass(),
			Aircraft:       f.Aircraft(),
			Amenities:      f.Amenities(),
			Baggage: response.Baggage{
				CarryOn: f.Baggage().CarryOn(),
				Checked: f.Baggage().Checked(),
			},
		})
	}

	return result
}

func formatDatetime(t time.Time) string {
	if t.IsZero() {
		return ""
	}

	return fmt.Sprintf("%sT%s", t.Format("2006-01-02"), t.Format("15:04:05Z07:00"))
}
