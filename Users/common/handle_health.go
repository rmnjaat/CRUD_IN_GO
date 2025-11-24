package common

import (
	"encoding/json"
	"log"
	"net/http"
)

func healthCheckhandler(w http.ResponseWriter, r *http.Request) {

	response_message := map[string]string{
		"status":  "ok",
		"message": "User service is healthy",
	}

	response_json, err := json.Marshal(response_message)

	if err != nil {
		http.Error(w, "Error encoding JSON", http.StatusInternalServerError)
		return
	}

	log.Println("Handled health check...")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(response_json)
}
