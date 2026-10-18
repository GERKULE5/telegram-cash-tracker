package httpx

import (
	"errors"
	"fmt"
	"net/http"
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
)

var validate = newValidator()

func newValidator() *validator.Validate {
	v := validator.New(validator.WithRequiredStructEnabled())
	v.RegisterTagNameFunc(func(field reflect.StructField) string {
		return strings.SplitN(field.Tag.Get("json"), ",", 2)[0]
	})
	return v
}

func Validate(target any) error {
	return validate.Struct(target)
}

func WriteValidationError(w http.ResponseWriter, err error) {
	message := "invalid request data"
	var fields validator.ValidationErrors
	if errors.As(err, &fields) && len(fields) > 0 {
		field := fields[0]
		message = fmt.Sprintf("%s: %s", field.Field(), field.Tag())
		if field.Param() != "" {
			message += "=" + field.Param()
		}
	}
	WriteError(w, ErrorResponse{Status: http.StatusBadRequest, Message: message})
}
