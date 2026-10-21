package user

import (
	"api/internal/httpx"
	"net/http"
)

var (
	errInvalidJSON       = httpx.ErrorResponse{Status: http.StatusBadRequest, Message: "invalid JSON body"}
	errUserAlreadyExists = httpx.ErrorResponse{Status: http.StatusConflict, Message: "user with this telegram_user_id already exists"}
	errNotFound          = httpx.ErrorResponse{Status: http.StatusNotFound, Message: "user not found"}
	errInvalidTelegramID = httpx.ErrorResponse{Status: http.StatusBadRequest, Message: "telegram_user_id must be an integer"}
	errInternal          = httpx.ErrorResponse{Status: http.StatusInternalServerError, Message: "internal server error"}
)
