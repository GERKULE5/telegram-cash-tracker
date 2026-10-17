package health

import (
	"encoding/json"
	"net/http"
	"time"
)

type Response struct {
	Status string    `json:"status"`
	Time   time.Time `json:"time"`
}

func Ping(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(Response{
		Status: "OK",
		Time:   time.Now(),
	})
}
