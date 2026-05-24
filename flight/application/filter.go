package application

import (
	"sort"
	"strings"
	"time"

	"flight-search-aggregator/flight/domain"
	"flight-search-aggregator/flight/request"
	"flight-search-aggregator/partners"
)

// filterFlights applies the core search criteria to the flight results.
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

// applyAdvancedFilters applies additional filters: price range, stops, airlines, duration.
func applyAdvancedFilters(flights []domain.Flight, req *request.Search) []domain.Flight {
	if len(flights) == 0 {
		return flights
	}

	var filtered []domain.Flight

	for _, f := range flights {
		if !matchesPriceRange(f, req.MinPrice, req.MaxPrice) {
			continue
		}

		if !matchesMaxStops(f, req.MaxStops) {
			continue
		}

		if !matchesAirlines(f, req.Airlines) {
			continue
		}

		if !matchesMaxDuration(f, req.MaxDuration) {
			continue
		}

		filtered = append(filtered, f)
	}

	return filtered
}

// Supported values: "price_asc", "price_desc", "duration_asc", "duration_desc",
// "departure_asc", "departure_desc", "arrival_asc", "arrival_desc".
// Defaults to price_asc.
func sortFlights(flights []domain.Flight, sortBy string) []domain.Flight {
	if len(flights) <= 1 {
		return flights
	}

	sort.SliceStable(flights, func(i, j int) bool {
		switch sortBy {
		case "price_desc":
			return flights[i].Price().Amount() > flights[j].Price().Amount()
		case "duration_asc":
			return flights[i].Duration().TotalMinutes() < flights[j].Duration().TotalMinutes()
		case "duration_desc":
			return flights[i].Duration().TotalMinutes() > flights[j].Duration().TotalMinutes()
		case "departure_asc":
			return flights[i].Departure().Datetime().Before(flights[j].Departure().Datetime())
		case "departure_desc":
			return flights[i].Departure().Datetime().After(flights[j].Departure().Datetime())
		case "arrival_asc":
			return flights[i].Arrival().Datetime().Before(flights[j].Arrival().Datetime())
		case "arrival_desc":
			return flights[i].Arrival().Datetime().After(flights[j].Arrival().Datetime())
		default:
			return flights[i].Price().Amount() < flights[j].Price().Amount()
		}
	})

	return flights
}

func matchesOrigin(f domain.Flight, origin string) bool {
	if origin == "" {
		return true
	}

	return strings.EqualFold(f.Departure().Airport(), origin)
}

func matchesDestination(f domain.Flight, destination string) bool {
	if destination == "" {
		return true
	}

	return strings.EqualFold(f.Arrival().Airport(), destination)
}

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

func matchesPriceRange(f domain.Flight, minPrice, maxPrice *int) bool {
	if minPrice != nil && f.Price().Amount() < *minPrice {
		return false
	}

	if maxPrice != nil && f.Price().Amount() > *maxPrice {
		return false
	}

	return true
}

func matchesMaxStops(f domain.Flight, maxStops *int) bool {
	if maxStops == nil {
		return true
	}

	return f.Stops() <= *maxStops
}

// Matches against airline code (e.g., "GA", "QZ") case-insensitively.
func matchesAirlines(f domain.Flight, airlines []string) bool {
	if len(airlines) == 0 {
		return true
	}

	for _, code := range airlines {
		if strings.EqualFold(f.Airline().Code(), code) {
			return true
		}
	}

	return false
}

func matchesMaxDuration(f domain.Flight, maxDuration *int) bool {
	if maxDuration == nil {
		return true
	}

	return f.Duration().TotalMinutes() <= *maxDuration
}
