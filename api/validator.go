package api

import (
	"simple_bank/util"

	"github.com/go-playground/validator/v10"
)

// validCurrency is a validator.Func type, which used to validate currency based on scalability.
var validCurrency validator.Func = func(fieldLevel validator.FieldLevel) bool {
	if currency, ok := fieldLevel.Field().Interface().(string); ok {
		return util.IsSupportedCurrency(currency)
	}
	return false
}
