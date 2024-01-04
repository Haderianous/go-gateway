package gateway

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type routerGroup struct {
	server     *gin.Engine
	group      *gin.RouterGroup
	controller Controller
}

func newRouterGroup(path string, s *gin.Engine, c Controller) RouterGroup {
	return &routerGroup{server: s, controller: c, group: s.Group(path)}
}

func (rg routerGroup) Group(path string) RouterGroup {
	rg.group = rg.group.Group(path)
	return rg
}

func (rg routerGroup) Get(path string, handlers ...Handler) {
	rg.group.GET(path, rg.matchRoute(handlers...)...)
}

func (rg routerGroup) Post(path string, handlers ...Handler) {
	rg.group.POST(path, rg.matchRoute(handlers...)...)
}

func (rg routerGroup) Put(path string, handlers ...Handler) {
	rg.group.PUT(path, rg.matchRoute(handlers...)...)
}

func (rg routerGroup) Delete(path string, handlers ...Handler) {
	rg.group.DELETE(path, rg.matchRoute(handlers...)...)
}

func (rg routerGroup) ServeHttp(w http.ResponseWriter, req *http.Request) {
	rg.server.ServeHTTP(w, req)
}

func (rg routerGroup) Middleware(handlers ...Handler) {
	hfs := rg.matchRoute(handlers...)
	rg.group.Use(hfs...)
}

func (rg routerGroup) matchRoute(handlers ...Handler) []gin.HandlerFunc {
	var hfs []gin.HandlerFunc
	for i, handler := range handlers {
		hfs = append(hfs, rg.getHandler(handler, i == len(handlers)-1))
	}
	return hfs
}

func (rg routerGroup) getHandler(handler Handler, shouldRespond bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req Request
		if _req, ok := c.Get("req"); ok {
			req = _req.(Request)
		} else {
			req = NewRequest(c, rg.controller.LanguageBundle())
			c.Set("req", req)
		}
		req.SetIsResponded(false)
		if rg.controller.Process(handler, req, shouldRespond) {
			c.Next()
		}
	}
}
