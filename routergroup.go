package gateway

import "net/http"

type RouterGroup interface {
	Group(path string) RouterGroup
	Get(path string, handlers ...Handler)
	Post(path string, handlers ...Handler)
	Put(path string, handlers ...Handler)
	Delete(path string, handlers ...Handler)
	ServeHttp(w http.ResponseWriter, req *http.Request)
	Middleware(handlers ...Handler)
}
