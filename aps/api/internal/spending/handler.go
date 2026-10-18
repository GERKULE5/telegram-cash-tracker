package spending

import (
	"api/internal/httpx"
	"context"
	"errors"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
)

type Repository interface {
	Create(context.Context, CreateParams) (Spending, error)
	List(context.Context, ListParams) ([]Spending, error)
	Update(context.Context, int, int, UpdateParams) (Spending, error)
	Delete(context.Context, int, int) error
	ListAll(context.Context, ListParams) ([]Spending, error)
	GetByID(context.Context, int) (Spending, error)
	UpdateByID(context.Context, int, UpdateParams) (Spending, error)
	DeleteByID(context.Context, int) error
}

type Handler struct {
	spendings Repository
}

func NewHandler(spendings Repository) *Handler {
	return &Handler{spendings: spendings}
}

// Create adds a spending.
// @Summary Create spending
// @Description User_id is the internal users.id.
// @Tags Spendings
// @Accept json
// @Produce json
// @Param request body CreateParams true "Spending data"
// @Success 201 {object} Spending
// @Failure 400 {object} httpx.ErrorResponse
// @Failure 404 {object} httpx.ErrorResponse "User not found"
// @Failure 500 {object} httpx.ErrorResponse
// @Router /spendings [post]
func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	var params CreateParams
	if err := httpx.DecodeJSON(r, &params); err != nil {
		httpx.WriteError(w, errInvalidJSON)
		return
	}
	if err := httpx.Validate(params); err != nil {
		httpx.WriteValidationError(w, err)
		return
	}
	result, err := h.spendings.Create(r.Context(), params)
	if err != nil {
		handleRepositoryError(w, err)
		return
	}
	httpx.WriteJSON(w, http.StatusCreated, result)
}

// List returns the user's spendings.
// @Summary List spendings
// @Description Sorted by created_at DESC, id DESC. from is inclusive; to is exclusive.
// @Tags Spendings
// @Produce json
// @Param user_id query int true "Internal user ID" minimum(1)
// @Param category query string false "Category" Enums(food,transport,housing,entertainment,health,purchases,other)
// @Param from query string false "Period start (RFC3339)"
// @Param to query string false "Period end (RFC3339, exclusive)"
// @Param limit query int false "Page size" default(20) minimum(1) maximum(100)
// @Param offset query int false "Records to skip" default(0) minimum(0)
// @Success 200 {array} Spending
// @Failure 400 {object} httpx.ErrorResponse
// @Failure 500 {object} httpx.ErrorResponse
// @Router /spendings [get]
func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	h.list(w, r, false)
}

func (h *Handler) list(w http.ResponseWriter, r *http.Request, admin bool) {
	query := r.URL.Query()
	var userID int
	var err error
	if !admin || query.Has("user_id") {
		userID, err = strconv.Atoi(query.Get("user_id"))
		if err != nil {
			httpx.WriteError(w, errInvalidUserID)
			return
		}
		if err := httpx.Validate(UserIDParams{UserID: userID}); err != nil {
			httpx.WriteValidationError(w, err)
			return
		}
	}
	params := ListParams{UserID: userID, Category: Category(query.Get("category")), Limit: 20}
	if query.Has("limit") {
		params.Limit, err = strconv.Atoi(query.Get("limit"))
		if err != nil {
			httpx.WriteError(w, errInvalidPagination)
			return
		}
	}
	if query.Has("offset") {
		params.Offset, err = strconv.Atoi(query.Get("offset"))
		if err != nil {
			httpx.WriteError(w, errInvalidPagination)
			return
		}
	}
	if query.Has("from") {
		from, err := time.Parse(time.RFC3339, query.Get("from"))
		if err != nil {
			httpx.WriteError(w, errInvalidPeriod)
			return
		}
		params.From = &from
	}
	if query.Has("to") {
		to, err := time.Parse(time.RFC3339, query.Get("to"))
		if err != nil {
			httpx.WriteError(w, errInvalidPeriod)
			return
		}
		params.To = &to
	}
	if err := httpx.Validate(params); err != nil {
		httpx.WriteValidationError(w, err)
		return
	}
	if params.From != nil && params.To != nil {
		if err := httpx.Validate(PeriodParams{From: *params.From, To: *params.To}); err != nil {
			httpx.WriteValidationError(w, err)
			return
		}
	}
	var result []Spending
	if admin {
		result, err = h.spendings.ListAll(r.Context(), params)
	} else {
		result, err = h.spendings.List(r.Context(), params)
	}
	if err != nil {
		handleRepositoryError(w, err)
		return
	}
	httpx.WriteJSON(w, http.StatusOK, result)
}

// Update changes the provided fields of a spending.
// @Summary Update spending
// @Description Omitted or null fields remain unchanged. Use an empty string to clear the comment. user_id is the internal users.id.
// @Tags Spendings
// @Accept json
// @Produce json
// @Param id path int true "Spending ID" minimum(1)
// @Param user_id query int true "Internal owner ID" minimum(1)
// @Param request body UpdateParams true "Fields to update"
// @Success 200 {object} Spending
// @Failure 400 {object} httpx.ErrorResponse
// @Failure 404 {object} httpx.ErrorResponse
// @Failure 500 {object} httpx.ErrorResponse
// @Router /spendings/{id} [patch]
func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	h.update(w, r, false)
}

func (h *Handler) update(w http.ResponseWriter, r *http.Request, admin bool) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		httpx.WriteError(w, errInvalidID)
		return
	}
	if err := httpx.Validate(IDParams{ID: id}); err != nil {
		httpx.WriteValidationError(w, err)
		return
	}
	var userID int
	if !admin {
		userID, err = strconv.Atoi(r.URL.Query().Get("user_id"))
		if err != nil {
			httpx.WriteError(w, errInvalidUserID)
			return
		}
		if err := httpx.Validate(UserIDParams{UserID: userID}); err != nil {
			httpx.WriteValidationError(w, err)
			return
		}
	}
	var params UpdateParams
	if err := httpx.DecodeJSON(r, &params); err != nil {
		httpx.WriteError(w, errInvalidJSON)
		return
	}
	if err := httpx.Validate(params); err != nil {
		httpx.WriteValidationError(w, err)
		return
	}
	var result Spending
	if admin {
		result, err = h.spendings.UpdateByID(r.Context(), id, params)
	} else {
		result, err = h.spendings.Update(r.Context(), id, userID, params)
	}
	if err != nil {
		handleRepositoryError(w, err)
		return
	}
	httpx.WriteJSON(w, http.StatusOK, result)
}

// Delete removes a spending.
// @Summary Delete spending
// @Tags Spendings
// @Produce json
// @Param id path int true "Spending ID" minimum(1)
// @Param user_id query int true "Internal owner ID" minimum(1)
// @Success 204 "Spending deleted"
// @Failure 400 {object} httpx.ErrorResponse
// @Failure 404 {object} httpx.ErrorResponse
// @Failure 500 {object} httpx.ErrorResponse
// @Router /spendings/{id} [delete]
func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	h.delete(w, r, false)
}

func (h *Handler) delete(w http.ResponseWriter, r *http.Request, admin bool) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		httpx.WriteError(w, errInvalidID)
		return
	}
	if err := httpx.Validate(IDParams{ID: id}); err != nil {
		httpx.WriteValidationError(w, err)
		return
	}
	var userID int
	if !admin {
		userID, err = strconv.Atoi(r.URL.Query().Get("user_id"))
		if err != nil {
			httpx.WriteError(w, errInvalidUserID)
			return
		}
		if err := httpx.Validate(UserIDParams{UserID: userID}); err != nil {
			httpx.WriteValidationError(w, err)
			return
		}
	}
	if admin {
		err = h.spendings.DeleteByID(r.Context(), id)
	} else {
		err = h.spendings.Delete(r.Context(), id, userID)
	}
	if err != nil {
		handleRepositoryError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func handleRepositoryError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, ErrNotFound):
		httpx.WriteError(w, errNotFound)
	case errors.Is(err, ErrUserNotFound):
		httpx.WriteError(w, errUserNotFound)
	default:
		log.Printf("spending: %v", err)
		httpx.WriteError(w, errInternal)
	}
}
