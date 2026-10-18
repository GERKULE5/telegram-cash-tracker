package spending

import (
	"api/internal/httpx"
	"errors"
	"net/http"
)

var (
	ErrNotFound     = errors.New("spending not found")
	ErrUserNotFound = errors.New("user not found")
)

var (
	errInvalidJSON       = httpx.ErrorResponse{Status: http.StatusBadRequest, Message: "invalid JSON body"}
	errInvalidUserID     = httpx.ErrorResponse{Status: http.StatusBadRequest, Message: "user_id must be a positive internal user ID"}
	errInvalidID         = httpx.ErrorResponse{Status: http.StatusBadRequest, Message: "id must be a positive integer"}
	errInvalidPagination = httpx.ErrorResponse{Status: http.StatusBadRequest, Message: "limit must be between 1 and 100; offset must be non-negative"}
	errInvalidPeriod     = httpx.ErrorResponse{Status: http.StatusBadRequest, Message: "from and to must be RFC3339; from must be earlier than to"}
	errNotFound          = httpx.ErrorResponse{Status: http.StatusNotFound, Message: "spending not found"}
	errUserNotFound      = httpx.ErrorResponse{Status: http.StatusNotFound, Message: "user not found"}
	errInternal          = httpx.ErrorResponse{Status: http.StatusInternalServerError, Message: "internal server error"}
)
