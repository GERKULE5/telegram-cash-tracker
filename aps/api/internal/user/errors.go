package user

import (
	"api/internal/httpx"
	"net/http"
)

var (
	errInvalidJSON       = httpx.ErrorResponse{Status: http.StatusBadRequest, Message: "invalid JSON body"}
	errUserAlreadyExists = httpx.ErrorResponse{Status: http.StatusConflict, Message: "user with this telegram_user_id already exists"}
	errInternal          = httpx.ErrorResponse{Status: http.StatusInternalServerError, Message: "internal server error"}
)
