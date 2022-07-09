package validator

import (
	"github.com/go-playground/validator/v10"
	"l.hilmy.dev/backend/helpers/errorhandler"
)

var validate = validator.New()

func Struct[T comparable](field *T) *error {
	if err := validate.Struct(field); err != nil {
		errorhandler.LogErrorThenContinue(&err)
		return &err
	}
	return nil
}
