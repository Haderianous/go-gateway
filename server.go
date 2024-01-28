package gateway

import (
	"context"
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	gormsessions "github.com/gin-contrib/sessions/gorm"
	"github.com/gin-gonic/gin"
	"github.com/haderianous/go-logger/logger"
	"gorm.io/gorm"
	"net/http"
	"time"
)

type Server interface {
	NewRouterGroup(path string) RouterGroup
	Shutdown(timeout time.Duration) error
	LoadHTMLGlob(pattern string)
	NewSession(sessionName string, secretKey string)
	NewCorsMiddleware()
	NewGormSession(db *gorm.DB, sessionName string, secretKey string)
	Run(...string) error
}

type server struct {
	engine     *gin.Engine
	httpServer http.Server // gin engine is inside this server
	group      *gin.RouterGroup
	logger     logger.Logger
	controller Controller
}

func NewServer(c Controller) Server {
	return &server{
		engine:     gin.New(),
		logger:     logger.NewLogger(logger.InfoLevel, logger.JsonEncoding),
		controller: c,
	}
}

func (s *server) NewRouterGroup(path string) RouterGroup {
	return newRouterGroup(path, s.engine, s.controller)
}

func (s *server) Shutdown(timeout time.Duration) error {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	return s.httpServer.Shutdown(ctx)
}

func (s *server) NewSession(sessionName string, secretKey string) {
	store := cookie.NewStore([]byte(secretKey))
	s.engine.Use(sessions.Sessions(sessionName, store))
}

func (s *server) NewGormSession(db *gorm.DB, sessionName string, secretKey string) {
	store := gormsessions.NewStore(db, true, []byte("secretKey"))
	s.engine.Use(sessions.Sessions(sessionName, store))
}

func (s *server) NewCorsMiddleware() {
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"*"}                                        // Specify allowed origins
	config.AllowMethods = []string{"GET", "POST", "OPTIONS"}                   // Specify allowed HTTP methods
	config.AllowHeaders = []string{"Origin", "Content-Length", "Content-Type"} // Specify allowed headers
	s.engine.Use(cors.New(config))
}

func (s *server) LoadHTMLGlob(pattern string) {
	s.engine.LoadHTMLGlob(pattern)
}

func (s *server) Run(host ...string) error {
	s.httpServer = http.Server{
		Addr:    host[0],
		Handler: s.engine,
	}
	if gin.IsDebugging() {
		s.logger.InfoF("Listening and serving HTTP on %s", host[0])
	}
	err := s.httpServer.ListenAndServe()
	if err == http.ErrServerClosed {
		return nil
	}
	return err
}
