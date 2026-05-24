package controller

import (
	"flight-search-aggregator/flight/application"

	"github.com/gorilla/mux"
)

type Controller interface {
	RegisterRoutes(router *mux.Router)
}

func GetController(router *mux.Router) Controller {
	flightApp := application.NewFlightApplication()

	return NewFlightController(router, flightApp)
}

func NewFlightController(
	router *mux.Router,
	flightApp application.Flight,
) Controller {
	ctrl := &flight{
		flightApp: flightApp,
	}

	ctrl.RegisterRoutes(router)

	return ctrl
}
