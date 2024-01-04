package gateway

import (
	"github.com/gin-gonic/gin"
	"github.com/haderianous/go-error"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"reflect"
	"time"
)

type Responder interface {
	Respond(Request, any)
	RespondError(Request, errors.ErrorModel)
	LanguageBundle() *i18n.Bundle
}

type responder struct {
	languageBundle *i18n.Bundle
}

type Response struct {
	Message       string `json:"message"`
	Error         string `json:"error,omitempty"`
	Version       string `json:"version"`
	RepresentedAt string `json:"represented_at"`
	Data          struct {
		Total   int `json:"total"`
		PerPage int `json:"per_page"`
		Result  any `json:"result"`
	} `json:"data"`
}

func NewResponder(languageBundle *i18n.Bundle) Responder {
	return &responder{languageBundle: languageBundle}
}

func (r *responder) LanguageBundle() *i18n.Bundle {
	return r.languageBundle
}

func (r *responder) Respond(req Request, result any) {
	req.SetIsResponded(true)
	response := Response{
		Message:       req.GetLanguage().Localize("SuccessMessage", "Operation has been done successfully"),
		Version:       "v1",
		RepresentedAt: time.Now().Format("2006-01-02 15:04:05"),
		Data: struct {
			Total   int `json:"total"`
			PerPage int `json:"per_page"`
			Result  any `json:"result"`
		}{
			Total:   req.Paginator().Total(),
			PerPage: req.Paginator().PerPage(),
			Result:  result,
		},
	}
	var status int
	switch req.GetMethod() {
	case http.MethodPost, http.MethodPut:
		if result == nil {
			status = http.StatusNoContent
			break
		}
		response.Message = req.GetLanguage().Localize("CreatedMessage", "Item has been created successfully")
		if req.GetMethod() == http.MethodPut {
			response.Message = req.GetLanguage().Localize("UpdatedMessage", "Item has been updated successfully")
		}
		response.Data.Result = []any{result}
		status = http.StatusCreated
	case http.MethodGet:
		status = http.StatusOK
		if result == nil || reflect.ValueOf(result).IsNil() {
			status = http.StatusNoContent
			break
		}
	case http.MethodDelete:
		if result == nil {
			status = http.StatusNoContent
			break
		}
		status = http.StatusOK
	}

	req.GetContext().(*gin.Context).JSON(status, response)
	return
}

func (r *responder) RespondError(req Request, err errors.ErrorModel) {
	ctx := req.GetContext().(*gin.Context)
	_ = ctx.Error(err)
	if err.ErrorId() != "" && (err.IsMsgDefault() || !err.IsIdDefault()) {
		err = err.WithMessage(req.GetLanguage().Localize(err.MessageId(), err.Message()))
		err = err.WithErrorText(req.GetLanguage().Localize(err.ErrorId(), err.ErrorText()))
	}
	ctx.JSON(getStatusCodeByError(err.Type()), err)
	ctx.Abort()
	return
}
