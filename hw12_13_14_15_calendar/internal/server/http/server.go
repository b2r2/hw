package httpserver

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"

	"github.com/b2r2/hw/hw12_13_14_15_calendar/internal/logger"

	"github.com/b2r2/hw/hw12_13_14_15_calendar/internal/app"
)

type Server struct {
	log    logger.Logger
	router chi.Router
	addr   string
	server *http.Server
}

type EventHandler struct {
	log logger.Logger
	app app.App
}

func NewHandler(log logger.Logger, app app.App) *EventHandler {
	return &EventHandler{app: app, log: log}
}

func NewRouter(log logger.Logger, h *EventHandler, v interface{}) *chi.Mux {
	r := chi.NewRouter()
	r.Use(cors.AllowAll().Handler)
	r.Use(middleware.Recoverer)
	r.Use(middleware.StripSlashes)
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Timeout(20 * time.Second))
	r.NotFound(notFoundHandler)
	r.Get("/hello", helloHandler)
	r.Get("/version", versionHandler(v))
	r.Route("/api", func(r chi.Router) {
		r.Group(func(r chi.Router) {
			r.Use(loggingMiddleware(log))
			r.Route("/v1", func(r chi.Router) {
				r.Post("/create", h.Create)
				r.Post("/update/{id}", h.Update)
				r.Get("/deleteAll", h.DeleteAll)
				r.Get("/delete/{id}", h.Delete)
				r.Get("/listAll", h.ListALl)
				r.Get("/listMonth", h.ListMonth)
				r.Get("/listWeek", h.ListWeek)
				r.Get("/listDay", h.ListDay)
				r.Get("/get/{id}", h.Get)
			})
		})
	})
	return r
}

func NewServer(log logger.Logger, r chi.Router, addr string) *Server {
	s := &Server{}
	s.log = log
	s.router = r
	s.addr = addr
	s.server = &http.Server{
		Addr:              s.addr,
		Handler:           s.router,
		ReadTimeout:       20 * time.Second,
		ReadHeaderTimeout: 20 * time.Second,
		WriteTimeout:      20 * time.Second,
	}
	return s
}

func (s *Server) Start() error {
	s.log.Infoln("server started")
	if err := s.server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return err
	}
	return nil
}

func (s *Server) Stop(ctx context.Context) error {
	s.log.Infoln("server stopped")
	return s.server.Shutdown(ctx)
}

func notFoundHandler(w http.ResponseWriter, _ *http.Request) {
	http.Error(w, "404 page not found,", http.StatusNotFound)
}

func helloHandler(w http.ResponseWriter, _ *http.Request) {
	_ = json.NewEncoder(w).Encode("Hello world")
}

func versionHandler(v interface{}) func(w http.ResponseWriter, _ *http.Request) {
	return func(w http.ResponseWriter, _ *http.Request) {
		_ = json.NewEncoder(w).Encode(v)
	}
}
