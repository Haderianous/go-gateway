package gateway

type Handler interface {
	Handle(request Request) (any, errors.ErrorModel)
}
