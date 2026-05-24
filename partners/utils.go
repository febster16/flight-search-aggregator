package partners

import (
	"fmt"
	"math"
	"regexp"
	"strconv"
	"time"
)

// ParseTime parses a time string using the given layout.
// Returns a wrapped error if parsing fails.
func ParseTime(raw string, layout string) (time.Time, error) {
	t, err := time.Parse(layout, raw)
	if err != nil {
		return time.Time{}, fmt.Errorf("failed to parse time %q with layout %q: %w", raw, layout, err)
	}

	return t, nil
}

// ParseTimeInLocation parses a time string using the given layout and
// associates it with the specified IANA timezone name (e.g. "Asia/Jakarta").
// Returns a wrapped error if the timezone cannot be loaded or parsing fails.
func ParseTimeInLocation(raw string, layout string, tz string) (time.Time, error) {
	loc, err := time.LoadLocation(tz)
	if err != nil {
		return time.Time{}, fmt.Errorf("failed to load timezone %q: %w", tz, err)
	}

	t, err := time.ParseInLocation(layout, raw, loc)
	if err != nil {
		return time.Time{}, fmt.Errorf("failed to parse time %q with layout %q in timezone %q: %w", raw, layout, tz, err)
	}

	return t, nil
}

// HoursToMinutes converts a fractional hour value to total minutes, rounded.
func HoursToMinutes(hours float64) int {
	return int(math.Round(hours * 60))
}

// ParseDurationString parses human-readable duration strings like "1h 45m",
// "2h", "30m" into total minutes.
func ParseDurationString(raw string) int {
	hourRe := regexp.MustCompile(`(\d+)h`)
	minRe := regexp.MustCompile(`(\d+)m`)

	var total int

	if matches := hourRe.FindStringSubmatch(raw); len(matches) == 2 {
		h, _ := strconv.Atoi(matches[1])
		total += h * 60
	}

	if matches := minRe.FindStringSubmatch(raw); len(matches) == 2 {
		m, _ := strconv.Atoi(matches[1])
		total += m
	}

	return total
}

// ExtractAirlineCode extracts the IATA airline code prefix from a flight code,
// which is the first 2-3 letters before any digits.
// For example, "QZ520" returns "QZ", "GA400" returns "GA".
// ref: https://en.wikipedia.org/wiki/List_of_airline_codes
func ExtractAirlineCode(flightCode string) string {
	for i, c := range flightCode {
		if c >= '0' && c <= '9' {
			return flightCode[:i]
		}
	}

	return flightCode
}

// TODO: For production use, consider replacing this with a full library such as:
// - github.com/mmcloughlin/openflights
// - github.com/gilby125/google-flights-api/iata
var airportCityMap = map[string]string{
	"CGK": "Jakarta",
	"DPS": "Denpasar",
	"SUB": "Surabaya",
	"UPG": "Makassar",
	"SOC": "Solo",
}

// AirportCity returns the city name for a given IATA airport code.
// Default to airport code if not found.
func AirportCity(code string) string {
	if city, ok := airportCityMap[code]; ok {
		return city
	}

	return code
}
