package gateway

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

type Request interface {
	GetContext() context.Context
	GetClientIp() string
	GetMethod() string
	GetFullPath() string
	GetHeader(key string) string
	Paginator() Paginator

	SetLanguage(lang Language)
	GetLanguage() Language

	GetBody() any
	SetBody(any)
	Request() *http.Request
	Writer() http.ResponseWriter
	Status() int
	Size() int
	BindRequest(req Validatable) (err errors.ErrorModel)
	GetQuery(key string) string
	GetParam(key string) string
	Filters() FilterParams
	BindFilters() errors.ErrorModel
	SetIsResponded(bool)
	IsResponded() bool
	Next()
	SetKey(key string, value any)
	GetKey(key string) (value any, exists bool)

	RespondHtml(status int, contentType string, body any)
}

type request struct {
	context     *gin.Context
	body        any
	language    Language
	paginator   Paginator
	filters     FilterParams
	isResponded bool
}

func NewRequest(ctx *gin.Context, languageBundle *i18n.Bundle) Request {
	req := &request{
		context: ctx,
	}
	if languageBundle != nil {
		acceptLang := ctx.Request.Header.Get("Accept-Language")
		req.SetLanguage(NewLanguage(languageBundle, acceptLang))
	}

	return req
}

func (r *request) GetContext() context.Context {
	return r.context
}

func (r *request) GetClientIp() string {
	return r.Request().RemoteAddr
}

func (r *request) GetMethod() string {
	return r.Request().Method
}

func (r *request) GetFullPath() string {
	return r.context.FullPath()
}

func (r *request) GetHeader(key string) string {
	return r.Request().Header.Get(key)
}

func (r *request) Paginator() Paginator {
	if r.paginator != nil {
		return r.paginator
	}
	p := NewPaginator()

	page, _ := strconv.Atoi(r.GetQuery("page"))
	p.SetPage(page)

	limit, _ := strconv.Atoi(r.GetQuery("limit"))
	p.SetLimit(limit)

	r.paginator = p

	return p
}

func (r *request) SetLanguage(lang Language) {
	r.language = lang
}

func (r *request) GetLanguage() Language {
	return r.language
}

func (r *request) GetBody() any {
	return r.body
}

func (r *request) SetBody(body any) {
	r.body = body
}

func (r *request) Request() *http.Request {
	return r.context.Request
}

func (r *request) Writer() http.ResponseWriter {
	return r.context.Writer
}

func (r *request) Status() int {
	return r.context.Writer.Status()
}

func (r *request) Size() int {
	return r.context.Writer.Size()
}

func (r *request) BindRequest(req Validatable) (err errors.ErrorModel) {
	e := r.context.ShouldBindUri(req)
	if e != nil {
		return errors.DefaultUnProcessable.WithError(e)
	}
	e = r.context.ShouldBind(req)
	if e != nil && e != io.EOF {
		return errors.DefaultUnProcessable.WithError(e)
	}
	e = r.context.ShouldBindHeader(req)
	if e != nil && e != io.EOF {
		return errors.DefaultUnProcessable.WithError(e)
	}
	body, e, errs := req.Validate(r.language)
	if e != nil {
		return errors.DefaultUnProcessable.WithError(e).WithErrors(errs)
	}
	if errs != nil {
		return errors.DefaultUnProcessable.WithErrors(errs)
	}
	r.body = body
	return nil
}

func (r *request) GetQuery(key string) string {
	return r.context.Query(key)
}

func (r *request) GetParam(key string) string {
	return r.context.Param(key)
}

func (r *request) Filters() FilterParams {
	return r.filters
}

func (r *request) SetIsResponded(responded bool) {
	r.isResponded = responded
}

func (r *request) IsResponded() bool {
	return r.isResponded
}

func (r *request) Next() {
	r.context.Next()
}

func (r *request) SetKey(key string, value any) {
	r.context.Set(key, value)
}

func (r *request) GetKey(key string) (value any, exists bool) {
	return r.context.Get(key)
}

func (r *request) RespondHtml(status int, name string, body any) {
	r.context.HTML(status, name, body)
	r.context.Abort()
}

func (r *request) BindFilters() errors.ErrorModel {
	qp, err := url.QueryUnescape(r.Request().URL.RawQuery)
	if err != nil {
		return errors.DefaultUnProcessable.WithError(err)
	}

	separatedQpCollection := strings.Split(qp, "&")
	sorts := make([]Sort, 0)
	filters := make([]Filter, 0)
	fpCollection, spCollection := make([]string, 0), make([]string, 0)

	for _, separatedQp := range separatedQpCollection {
		if strings.Contains(separatedQp, "filters[") {
			fpCollection = append(fpCollection, separatedQp)
		} else if strings.Contains(separatedQp, "sorts[") {
			spCollection = append(spCollection, separatedQp)
		}
	}

	if len(spCollection) > 0 {
		firstSpIndex := spCollection[0][0:8]
		sorts = fillSortsByEntity(spCollection, firstSpIndex)
	}

	if len(fpCollection) > 0 {
		firstFpIndex := fpCollection[0][0:10]
		filters = fillFiltersByEntity(fpCollection, firstFpIndex)
	}

	r.filters = FilterParams{
		Filters: filters,
		Sorts:   sorts,
		Page:    r.Paginator().Page(),
		Limit:   r.Paginator().PerPage(),
	}
	return nil
}

func fillSortsByEntity(spCollection []string, sStr string) (sorts []Sort) {
	var s Sort
	for index := 0; index < len(spCollection); index++ {
		if strings.Contains(spCollection[index], sStr) {
			keyValues := strings.Split(spCollection[index], "=")
			k := keyValues[0][9:10]

			if k == "k" {
				s.Key = keyValues[1]
			} else if k == "v" {
				s.Value = keyValues[1]
			}
			spCollection = append(spCollection[:index], spCollection[index+1:]...)
			index = -1
		}
	}
	sorts = append(sorts, s)
	if len(spCollection) > 0 {
		sort := fillSortsByEntity(spCollection, spCollection[0][0:8])
		sorts = append(sorts, sort...)
	}
	return
}

func fillFiltersByEntity(fpCollection []string, sStr string) (filters []Filter) {
	var v []interface{}
	var f Filter
	for index := 0; index < len(fpCollection); index++ {
		if strings.Contains(fpCollection[index], sStr) {
			keyValues := strings.Split(fpCollection[index], "=")
			k := keyValues[0][11:12]
			if k == "k" {
				f.Key = keyValues[1]
			} else if k == "v" {
				v = append(v, keyValues[1])
				f.Value = v
			} else if k == "o" {
				f.Op = Operation(keyValues[1])
			}
			fpCollection = append(fpCollection[:index], fpCollection[index+1:]...)
			index = -1
		}
	}
	filters = append(filters, f)
	if len(fpCollection) > 0 {
		fil := fillFiltersByEntity(fpCollection, fpCollection[0][0:10])
		filters = append(filters, fil...)
	}
	return
}
