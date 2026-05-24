package application

import (
	"strings"
	"time"

	"flight-search-aggregator/flight/domain"
	"flight-search-aggregator/partners"
)

// filterFlights applies the search criteria to the aggregated flight results.
// It filters by origin, destination, departure date, passenger count (based on
// available seats), and cabin class.
func filterFlights(flights []domain.Flight, req *partners.SearchRequest) []domain.Flight {
	if len(flights) == 0 {
		return flights
	}

	var filtered []domain.Flight

	for _, f := range flights {
		if !matchesOrigin(f, req.Origin) {
			continue
		}

		if !matchesDestination(f, req.Destination) {
			continue
		}

		if !matchesDepartureDate(f, req.DepartureDate) {
			continue
		}

		if !matchesPassengers(f, req.Passengers) {
			continue
		}

		if !matchesCabinClass(f, req.CabinClass) {
			continue
		}

		filtered = append(filtered, f)
	}

	return filtered
}

// matchesOrigin checks if the flight departs from the requested origin airport.
func matchesOrigin(f domain.Flight, origin string) bool {
	if origin == "" {
		return true
	}

	return strings.EqualFold(f.Departure().Airport(), origin)
}

// matchesDestination checks if the flight arrives at the requested destination airport.
func matchesDestination(f domain.Flight, destination string) bool {
	if destination == "" {
		return true
	}

	return strings.EqualFold(f.Arrival().Airport(), destination)
}

// matchesDepartureDate checks if the flight departs on the requested date.
// Compares date portion only (YYYY-MM-DD), ignoring time-of-day.
func matchesDepartureDate(f domain.Flight, departureDate string) bool {
	if departureDate == "" {
		return true
	}

	requestedDate, err := time.Parse("2006-01-02", departureDate)
	if err != nil {
		// If the date can't be parsed, skip this filter (don't exclude results).
		return true
	}

	flightDate := f.Departure().Datetime()

	return flightDate.Year() == requestedDate.Year() &&
		flightDate.Month() == requestedDate.Month() &&
		flightDate.Day() == requestedDate.Day()
}

// matchesPassengers checks if the flight has enough available seats for the
// requested number of passengers.
func matchesPassengers(f domain.Flight, passengers int) bool {
	if passengers <= 0 {
		return true
	}

	return f.AvailableSeats() >= passengers
}

func matchesCabinClass(f domain.Flight, cabinClass string) bool {
	if cabinClass == "" {
		return true
	}

	return f.CabinClass() == cabinClass
}
