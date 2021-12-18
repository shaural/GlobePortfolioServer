package api

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/shaural/GlobePersonalWebsite/server/pkg/db"
)

func AddMapHandler(router *mux.Router) {
	mapRouter := router.PathPrefix("/api/map").Subrouter()
	mapRouter.HandleFunc("/country", handleGetCountry).Methods("GET")
	mapRouter.HandleFunc("/state", handleGetAllStates).Methods("GET")
	mapRouter.HandleFunc("/state/{country:[a-zA-Z]+}", handleGetState).Methods("GET")
}

func handleGetCountry(w http.ResponseWriter, r *http.Request) {
	ldb, err := db.NewDatabase(r.Context())
	if CheckAndWriteError(err, http.StatusInternalServerError, w) {
		return
	}
	defer ldb.Close()
	countries, err := ldb.GetCountries()
	if CheckAndWriteError(err, http.StatusInternalServerError, w) {
		return
	}
	SendJSONResponse(countries, w)
}

func handleGetAllStates(w http.ResponseWriter, r *http.Request) {
	ldb, err := db.NewDatabase(r.Context())
	if CheckAndWriteError(err, http.StatusInternalServerError, w) {
		return
	}
	defer ldb.Close()
	states, err := ldb.GetStates("")
	if CheckAndWriteError(err, http.StatusInternalServerError, w) {
		return
	}
	SendJSONResponse(states, w)
}

func handleGetState(w http.ResponseWriter, r *http.Request) {
	ldb, err := db.NewDatabase(r.Context())
	if CheckAndWriteError(err, http.StatusInternalServerError, w) {
		return
	}
	defer ldb.Close()
	country, ok := mux.Vars(r)["country"]
	if !ok {
		http.Error(w, "no country id provided", http.StatusBadRequest)
		return
	}
	states, err := ldb.GetStates(country)
	if CheckAndWriteError(err, http.StatusInternalServerError, w) {
		return
	}
	SendJSONResponse(states, w)
}
