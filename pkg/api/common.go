package api

import (
	"encoding/json"
	"log"
	"net/http"
)

func CheckAndWriteError(err error, statusCode int, w http.ResponseWriter) bool {
	if err == nil {
		return false
	}
	log.Println(err.Error())
	http.Error(w, err.Error(), statusCode)
	return true
}

func SendJSONResponse(data interface{}, w http.ResponseWriter) {
	jsonData, err := json.Marshal(data)
	if CheckAndWriteError(err, http.StatusInternalServerError, w) {
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonData)
}
