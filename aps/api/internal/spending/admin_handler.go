package spending

import (
	"api/internal/httpx"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type AdminHandler struct {
	handler *Handler
}

func NewAdminHandler(spendings Repository) *AdminHandler {
	return &AdminHandler{handler: NewHandler(spendings)}
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
	h.handler.Create(w, r)
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
	h.handler.list(w, r, true)
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
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		httpx.WriteError(w, errInvalidID)
		return
	}
	if err := httpx.Validate(IDParams{ID: id}); err != nil {
		httpx.WriteValidationError(w, err)
		return
	}
	result, err := h.handler.spendings.GetByID(r.Context(), id)
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
	h.handler.update(w, r, true)
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
	h.handler.delete(w, r, true)
}
