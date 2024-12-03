package main

import "net/http"

func handlerReadiness(w http.ResponseWriter, r *http.Request) {
	data := struct {
		Status string `json:"status"`
	}{
		Status: "available",
	}
	respondWithJSON(w, 200, data)
}
