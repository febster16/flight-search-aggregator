package domain

import (
	"fmt"
	"time"
)

type Flight struct {
	id             string
	provider       string
	airline        Airline
	flightNumber   string
	departure      Endpoint
	arrival        Endpoint
	duration       Duration
	stops          int
	price          Price
	availableSeats int
	cabinClass     string
	aircraft       *string
	amenities      []string
	baggage        Baggage
}

type Airline struct {
	name string
	code string
}

type Endpoint struct {
	airport  string
	city     string
	datetime time.Time
}

type Duration struct {
	totalMinutes int
}

type Price struct {
	amount   int
	currency string
}

type Baggage struct {
	carryOn string
	checked string
}

func NewFlight(
	id string,
	provider string,
	airline Airline,
	flightNumber string,
	departure Endpoint,
	arrival Endpoint,
	duration Duration,
	stops int,
	price Price,
	availableSeats int,
	cabinClass string,
	aircraft *string,
	amenities []string,
	baggage Baggage,
) Flight {
	return Flight{
		id:             id,
		provider:       provider,
		airline:        airline,
		flightNumber:   flightNumber,
		departure:      departure,
		arrival:        arrival,
		duration:       duration,
		stops:          stops,
		price:          price,
		availableSeats: availableSeats,
		cabinClass:     cabinClass,
		aircraft:       aircraft,
		amenities:      amenities,
		baggage:        baggage,
	}
}

func (f Flight) ID() string           { return f.id }
func (f Flight) Provider() string     { return f.provider }
func (f Flight) Airline() Airline     { return f.airline }
func (f Flight) FlightNumber() string { return f.flightNumber }
func (f Flight) Departure() Endpoint  { return f.departure }
func (f Flight) Arrival() Endpoint    { return f.arrival }
func (f Flight) Duration() Duration   { return f.duration }
func (f Flight) Stops() int           { return f.stops }
func (f Flight) Price() Price         { return f.price }
func (f Flight) AvailableSeats() int  { return f.availableSeats }
func (f Flight) CabinClass() string   { return f.cabinClass }
func (f Flight) Aircraft() *string    { return f.aircraft }
func (f Flight) Amenities() []string  { return f.amenities }
func (f Flight) Baggage() Baggage     { return f.baggage }

func NewAirline(
	name, code string,
) Airline {
	return Airline{
		name: name,
		code: code,
	}
}

func (a Airline) Name() string { return a.name }
func (a Airline) Code() string { return a.code }

func NewEndpoint(
	airport, city string, datetime time.Time,
) Endpoint {
	return Endpoint{
		airport:  airport,
		city:     city,
		datetime: datetime,
	}
}

func (e Endpoint) Airport() string     { return e.airport }
func (e Endpoint) City() string        { return e.city }
func (e Endpoint) Datetime() time.Time { return e.datetime }
func (e Endpoint) Timestamp() int64    { return e.datetime.Unix() }

func NewDuration(
	totalMinutes int,
) Duration {
	return Duration{
		totalMinutes: totalMinutes,
	}
}

func (d Duration) TotalMinutes() int { return d.totalMinutes }

func (d Duration) Formatted() string {
	h := d.totalMinutes / 60
	m := d.totalMinutes % 60
	if h == 0 {
		return fmt.Sprintf("%dm", m)
	}
	if m == 0 {
		return fmt.Sprintf("%dh", h)
	}
	return fmt.Sprintf("%dh %dm", h, m)
}

func NewPrice(
	amount int,
	currency string,
) Price {
	return Price{
		amount:   amount,
		currency: currency,
	}
}

func (p Price) Amount() int      { return p.amount }
func (p Price) Currency() string { return p.currency }

func NewBaggage(
	carryOn, checked string,
) Baggage {
	return Baggage{
		carryOn: carryOn,
		checked: checked,
	}
}

func (b Baggage) CarryOn() string { return b.carryOn }
func (b Baggage) Checked() string { return b.checked }
