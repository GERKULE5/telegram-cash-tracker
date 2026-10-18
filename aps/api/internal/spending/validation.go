package spending

import (
	"api/internal/httpx"

	"github.com/go-playground/validator/v10"
)

func init() {
	httpx.RegisterStructValidation(validateUpdateParams, UpdateParams{})
}

func validateUpdateParams(level validator.StructLevel) {
	params := level.Current().Interface().(UpdateParams)
	if params.Amount == nil && params.Category == nil && params.SpentAt == nil && params.Comment == nil {
		level.ReportError(params, "", "", "notempty", "")
	}
}
