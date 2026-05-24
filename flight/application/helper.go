package application

import (
	"flight-search-aggregator/flight/domain"
	"flight-search-aggregator/flight/response"
	"fmt"
	"time"
)

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
