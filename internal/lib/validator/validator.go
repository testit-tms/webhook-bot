package validator

import (
	"strings"

	"github.com/go-playground/validator/v10"
)

var (
	parseMode = []string{"markdownv2", "markdown", "html"}
)

func ValidateParseMode(fl validator.FieldLevel) bool {
	value := fl.Field().String()

	if value == "" {
		return true
	}

	for _, v := range parseMode {
		if v == strings.ToLower(value) {
			return true
		}
	}

	return false
}
