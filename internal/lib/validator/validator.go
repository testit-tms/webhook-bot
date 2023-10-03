package validator

import (
	"strings"

	"github.com/go-playground/validator/v10"
)

var (
	parseMode = []string{"html"} //"markdownv2", "markdown",
)

// ValidateParseMode checks if the provided parse mode is valid.
// Returns true if the parse mode is valid, false otherwise.
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
