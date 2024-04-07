package validation

import (
	"regexp"

	"github.com/go-playground/validator/v10"
)

func isHourTime(fl validator.FieldLevel) bool {
	pattern := "^([0-1]?[0-9]|2[0-3]):[0-5][0-9]$"
	regex := regexp.MustCompile(pattern)

	return regex.MatchString(fl.Field().String())
}

func Validate(data interface{}) error {
	v := validator.New()

	if err := v.RegisterValidation("hourtime", isHourTime); err != nil {
		return err
	}

	return v.Struct(data)
}
