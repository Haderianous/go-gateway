package gateway

import "github.com/haderianous/go-logger/logger"

type Controller interface {
	Responder
	Processor
}

type controller struct {
	Responder
	log logger.Logger
}

func NewController(responder Responder, log logger.Logger) Controller {
	return &controller{responder, log}
}

func (c *controller) Process(handler Handler, req Request, respond bool) bool {
	result, err := handler.Handle(req)
	if err != nil {
		c.log.With(logger.Field{
			"message": err.Message(),
			"error":   err.ErrorText(),
			"detail":  err.Detail(),
		}).ErrorF("error on handler func")
		c.RespondError(req, err)
		return false
	}

	if !req.IsResponded() && respond {
		c.Respond(req, result)
		return false
	}
	return true
}
