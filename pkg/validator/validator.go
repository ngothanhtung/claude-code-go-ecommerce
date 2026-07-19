package validator

import (
	"strings"

	"github.com/go-playground/validator/v10"

	apperr "github.com/ngothanhtung/go-tutorials/internal/common/errors"
)

var validate = validator.New()

// Struct validates s and returns an AppError (validation) on failure.
func Struct(s interface{}) error {
	if err := validate.Struct(s); err != nil {
		if verrs, ok := err.(validator.ValidationErrors); ok {
			msgs := make([]string, 0, len(verrs))
			for _, fe := range verrs {
				msgs = append(msgs, fieldErrMsg(fe))
			}
			return apperr.NewValidation(strings.Join(msgs, "; "))
		}
		return apperr.NewValidation(err.Error())
	}
	return nil
}

// WrapBind converts a gin binding error into an AppError (validation).
func WrapBind(err error) error {
	return apperr.NewValidation("invalid request: " + err.Error())
}

func fieldErrMsg(fe validator.FieldError) string {
	field := fe.Field()
	switch fe.Tag() {
	case "required":
		return field + " is required"
	case "email":
		return field + " must be a valid email"
	case "min":
		return field + " must be at least " + fe.Param() + " characters"
	case "max":
		return field + " must be at most " + fe.Param() + " characters"
	case "gte":
		return field + " must be >= " + fe.Param()
	case "lte":
		return field + " must be <= " + fe.Param()
	default:
		return field + " is invalid (" + fe.Tag() + ")"
	}
}

// Var validates a single variable.
func Var(field interface{}, tag string) error {
	return validate.Var(field, tag)
}
