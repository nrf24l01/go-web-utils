package middleware

import (
	"regexp"

	validator "github.com/go-playground/validator/v10"
)

func RegisterRegexValidator(v *validator.Validate) {
	v.RegisterValidation("regex", func(fl validator.FieldLevel) bool {
		pattern := fl.Param()
		if pattern == "" {
			return false
		}
		re, err := regexp.Compile(pattern)
		if err != nil {
			return false
		}
		val := fl.Field().String()
		return re.MatchString(val)
	})
}
