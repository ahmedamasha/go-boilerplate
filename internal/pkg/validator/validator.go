package validator

import (
	"github.com/go-playground/validator/v10"
)

var V = validator.New()

func Validate(s interface{}) error {
	return V.Struct(s)
}
