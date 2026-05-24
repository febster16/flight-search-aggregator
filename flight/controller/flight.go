package controller

import (
	"encoding/json"
	"flight-search-aggregator/flight/application"
	"flight-search-aggregator/flight/request"
	"net/http"

	"github.com/gorilla/mux"
)

type flight struct {
	flightApp application.Flight
}

func (f *flight) RegisterRoutes(router *mux.Router) {
	subrouterFlight := router.PathPrefix("/flights").Subrouter()

	subrouterFlight.HandleFunc("/search", f.searchHandler).Methods(http.MethodPost)
}

func (f *flight) searchHandler(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()

	searchReq := &request.Search{}
	if err := json.NewDecoder(req.Body).Decode(searchReq); err != nil {
		writeErrorResponse(w, http.StatusBadRequest, "invalid request body: "+err.Error())

		return
	}

	// TODO: Implement search request validation using: github.com/go-playground/validator/v10 v10.15.4

	searchRes, err := f.flightApp.Search(ctx, searchReq)
	if err != nil {
		writeErrorResponse(w, http.StatusInternalServerError, "error searching flights: "+err.Error())

		return
	}

	writeJSONResponse(w, http.StatusOK, searchRes)
}

func writeJSONResponse(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(data)
}

func writeErrorResponse(w http.ResponseWriter, statusCode int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(map[string]string{
		"error": message,
	})
}
