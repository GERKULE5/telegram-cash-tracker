package httpx

import "net/http"

type ErrorResponse struct {
	Status  int    `json:"-"`
	Message string `json:"error"`
}

func WriteError(w http.ResponseWriter, response ErrorResponse) {
	WriteJSON(w, response.Status, response)
}
