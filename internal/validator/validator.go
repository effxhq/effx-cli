package validator

import (
	"github.com/effxhq/effx-go/data"
	v "gopkg.in/go-playground/validator.v9"
)

func ValidateObject(object *data.Data) error {
	validate := v.New()
	return validate.Struct(object)
}
