package gateway

import (
	"encoding/json"
	errors "github.com/haderianous/go-error"
	"github.com/haderianous/go-logger/logger"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"github.com/stretchr/testify/assert"
	"golang.org/x/text/language"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestServer_Run(t *testing.T) {
	respond := NewResponder(i18n.NewBundle(language.English))
	c := NewController(respond, logger.NewLogger(logger.WarnLevel, logger.JsonEncoding))
	s := NewServer(c)
	rg := s.NewRouterGroup("test")
	rg.Get("success", NewMiddleware(), NewHelloHandler())
	rg.Get("err", NewErrorHandler())

	req, _ := http.NewRequest(http.MethodGet, "/test/success", nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept-Language", "fa")
	w := httptest.NewRecorder()

	rg.ServeHttp(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	var res Response
	err := json.NewDecoder(w.Body).Decode(&res)
	assert.Nil(t, err)

	assert.Len(t, res.Data.Result, 2)

	req, _ = http.NewRequest(http.MethodGet, "/test/err", nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept-Language", "fa")
	w = httptest.NewRecorder()

	rg.ServeHttp(w, req)
	assert.Equal(t, http.StatusForbidden, w.Code)

}

type handler struct {
}

func NewHelloHandler() Handler {
	return &handler{}
}

func (h *handler) Handle(req Request) (any, errors.ErrorModel) {
	type user struct {
		Name     string `json:"name"`
		Age      int    `json:"age"`
		Location string `json:"location"`
	}
	users := []user{
		{
			Name:     "ali",
			Age:      25,
			Location: "turkey",
		},
		{
			Name:     "saeed",
			Age:      30,
			Location: "berlin",
		},
	}
	req.Paginator().SetTotal(2)
	return users, nil
}

type errHandler struct {
}

func NewErrorHandler() Handler {
	return &errHandler{}
}

func (h *errHandler) Handle(req Request) (any, errors.ErrorModel) {
	return nil, errors.DefaultForbiddenError
}

type middleware struct{}

func NewMiddleware() Handler {
	return &middleware{}
}

func (h *middleware) Handle(req Request) (any, errors.ErrorModel) {
	return nil, nil
}
