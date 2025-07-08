package api

import (
	"ecom/server/handlers"
	"net/http"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v4"
)

type App struct {
	hs *handlers.Handlers
}

func NewApp(hs *handlers.Handlers) *App {
	return &App{
		hs: hs,
	}
}
func (app *App) Run(addr string) error {
	m := chi.NewMux()
	m.Use(middleware.Recoverer)
	m.Use(middleware.Logger)
	m.Get("/", app.hs.HandleHome)
	return http.ListenAndServe(addr, m)
}
