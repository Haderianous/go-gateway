package gateway

import (
	errors "github.com/haderianous/go-error"
	"net/http"
)

func getStatusCodeByError(typ errors.Type) int {
	switch typ {
	case errors.TypeUnProcessable:
		return http.StatusUnprocessableEntity
	case errors.TypeNotFound:
		return http.StatusNotFound
	case errors.TypeUnAuthorized:
		return http.StatusUnauthorized
	case errors.TypeForbidden:
		return http.StatusForbidden
	case errors.TypeUnAvailable:
		return http.StatusServiceUnavailable
	case errors.TypeDuplicate:
		return http.StatusUnprocessableEntity
	case errors.TypeBadRequest:
		return http.StatusBadRequest
	case errors.TypeConflict:
		return http.StatusConflict
	case errors.TypeAccepted:
		return http.StatusAccepted
	}
	return http.StatusInternalServerError
}
