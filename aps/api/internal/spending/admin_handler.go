package spending

import (
	"api/internal/httpx"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
)

type AdminHandler struct {
	spendings Repository
}

func NewAdminHandler(spendings Repository) *AdminHandler {
	return &AdminHandler{spendings: spendings}
}

// Create adds a spending for any user.
// @Summary Create spending for any user
// @Description user_id in the body is the internal users.id.
// @Tags Admin Spendings
// @Accept json
// @Produce json
// @Param request body CreateParams true "Spending data"
// @Success 201 {object} Spending
// @Failure 400 {object} httpx.ErrorResponse
// @Failure 404 {object} httpx.ErrorResponse
// @Failure 500 {object} httpx.ErrorResponse
// @Router /admin/spendings [post]
func (h *AdminHandler) Create(w http.ResponseWriter, r *http.Request) {
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

// List returns spendings across all users.
// @Summary List all family spendings
// @Description Sorted by created_at DESC, id DESC. from is inclusive; to is exclusive.
// @Tags Admin Spendings
// @Produce json
// @Param user_id query int false "Internal user ID filter" minimum(1)
// @Param category query string false "Category" Enums(food,transport,housing,entertainment,health,purchases,other)
// @Param from query string false "Period start (RFC3339)"
// @Param to query string false "Period end (RFC3339, exclusive)"
// @Param limit query int false "Page size" default(20) minimum(1) maximum(100)
// @Param offset query int false "Records to skip" default(0) minimum(0)
// @Success 200 {array} Spending
// @Failure 400 {object} httpx.ErrorResponse
// @Failure 500 {object} httpx.ErrorResponse
// @Router /admin/spendings [get]
func (h *AdminHandler) List(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	var userID int32
	var err error
	if query.Has("user_id") {
		value, parseErr := strconv.ParseInt(query.Get("user_id"), 10, 32)
		if parseErr != nil {
			httpx.WriteError(w, errInvalidUserID)
			return
		}
		userID = int32(value)
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
	result, err := h.spendings.ListAll(r.Context(), params)
	if err != nil {
		handleRepositoryError(w, err)
		return
	}
	httpx.WriteJSON(w, http.StatusOK, result)
}

// Get returns a spending by ID regardless of its owner.
// @Summary Get spending
// @Tags Admin Spendings
// @Produce json
// @Param id path int true "Spending ID" minimum(1)
// @Success 200 {object} Spending
// @Failure 400 {object} httpx.ErrorResponse
// @Failure 404 {object} httpx.ErrorResponse
// @Failure 500 {object} httpx.ErrorResponse
// @Router /admin/spendings/{id} [get]
func (h *AdminHandler) Get(w http.ResponseWriter, r *http.Request) {
	value, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 32)
	if err != nil {
		httpx.WriteError(w, errInvalidID)
		return
	}
	id := int32(value)
	if err := httpx.Validate(IDParams{ID: id}); err != nil {
		httpx.WriteValidationError(w, err)
		return
	}
	result, err := h.spendings.GetByID(r.Context(), id)
	if err != nil {
		handleRepositoryError(w, err)
		return
	}
	httpx.WriteJSON(w, http.StatusOK, result)
}

// Update edits a spending regardless of its owner.
// @Summary Update any spending
// @Description Omitted or null fields remain unchanged. Use an empty string to clear the comment.
// @Tags Admin Spendings
// @Accept json
// @Produce json
// @Param id path int true "Spending ID" minimum(1)
// @Param request body UpdateParams true "Fields to update"
// @Success 200 {object} Spending
// @Failure 400 {object} httpx.ErrorResponse
// @Failure 404 {object} httpx.ErrorResponse
// @Failure 500 {object} httpx.ErrorResponse
// @Router /admin/spendings/{id} [patch]
func (h *AdminHandler) Update(w http.ResponseWriter, r *http.Request) {
	value, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 32)
	if err != nil {
		httpx.WriteError(w, errInvalidID)
		return
	}
	id := int32(value)
	if err := httpx.Validate(IDParams{ID: id}); err != nil {
		httpx.WriteValidationError(w, err)
		return
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
	result, err := h.spendings.UpdateByID(r.Context(), id, params)
	if err != nil {
		handleRepositoryError(w, err)
		return
	}
	httpx.WriteJSON(w, http.StatusOK, result)
}

// Delete removes a spending regardless of its owner.
// @Summary Delete any spending
// @Tags Admin Spendings
// @Produce json
// @Param id path int true "Spending ID" minimum(1)
// @Success 204 "Spending deleted"
// @Failure 400 {object} httpx.ErrorResponse
// @Failure 404 {object} httpx.ErrorResponse
// @Failure 500 {object} httpx.ErrorResponse
// @Router /admin/spendings/{id} [delete]
func (h *AdminHandler) Delete(w http.ResponseWriter, r *http.Request) {
	value, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 32)
	if err != nil {
		httpx.WriteError(w, errInvalidID)
		return
	}
	id := int32(value)
	if err := httpx.Validate(IDParams{ID: id}); err != nil {
		httpx.WriteValidationError(w, err)
		return
	}
	err = h.spendings.DeleteByID(r.Context(), id)
	if err != nil {
		handleRepositoryError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
