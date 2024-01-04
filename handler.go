package gateway

import errors "github.com/haderianous/go-error"

type Handler interface {
	Handle(request Request) (any, errors.ErrorModel)
}
