package handler

import (
	"github.com/Koshsky/erp/api/internal/service"
)

func isNotFoundError(err error) bool {
	return service.IsNotFoundError(err)
}

func isValidationError(err error) bool {
	return service.IsValidationError(err)
}
