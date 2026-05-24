package application

import (
	"flight-search-aggregator/flight/domain"
	"flight-search-aggregator/flight/request"
	"flight-search-aggregator/flight/response"
	"fmt"
	"time"
)

func mapSearchCriteriaResponse(req *request.Search) response.SearchCriteria {
	return response.SearchCriteria{
		Origin:        req.Origin,
		Destination:   req.Destination,
		DepartureDate: req.DepartureDate,
		Passengers:    req.Passengers,
		CabinClass:    req.CabinClass,
	}
}

func mapFlightsToResponse(flights []domain.Flight) []response.Flight {
	result := make([]response.Flight, 0, len(flights))

	for _, f := range flights {
		result = append(result, mapFlightToResponse(f))
	}

	return result
}

func mapFlightToResponse(f domain.Flight) response.Flight {
	amenities := f.Amenities()
	if amenities == nil {
		amenities = []string{}
	}

	return response.Flight{
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
		Amenities:      amenities,
		Baggage: response.Baggage{
			CarryOn: f.Baggage().CarryOn(),
			Checked: f.Baggage().Checked(),
		},
	}
}

func formatDatetime(t time.Time) string {
	if t.IsZero() {
		return ""
	}

	return fmt.Sprintf("%sT%s", t.Format("2006-01-02"), t.Format("15:04:05Z07:00"))
}
