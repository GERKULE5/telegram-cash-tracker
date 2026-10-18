package health

import (
	"api/internal/httpx"
	"net/http"
	"time"
)

type Response struct {
	Status string    `json:"status"`
	Time   time.Time `json:"time"`
}

func Ping(w http.ResponseWriter, _ *http.Request) {
	httpx.WriteJSON(w, http.StatusOK, Response{
		Status: "OK",
		Time:   time.Now(),
	})
}
