package api

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

func AddMainHandler(router *mux.Router) {
	router.HandleFunc("/", handleStatus)
	router.HandleFunc("/statusCheck", handleStatus)
}

func handleStatus(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "SUCCESS")
}
