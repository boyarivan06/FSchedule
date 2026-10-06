package api

import (
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func NewRouter(h *Handler, logger *slog.Logger) http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.ClientIPFromRemoteAddr)
	r.Use(middleware.Logger)
	/// r.Use(recoverer(logger)) TODO: WHY?
	r.Route("/api/v1", func(r chi.Router) {
		// users
		r.Get("/users/{id}", h.GetUser)

		r.Post("/users", h.CreateUser)
		r.Put("/users/{id}", h.UpdateUser)
		r.Delete("/users/{id}", h.DeleteUser)
		// children
		/*r.Get("/children", h.GetChildren)
		r.Get("/children/{id}", h.GetChild)
		r.Post("/children", h.CreateChild)
		r.Put("/children/{id}", h.UpdateChild)
		r.Delete("/children/{id}", h.DeleteChild)
		// events
		r.Get("/events", h.GetEvents)
		r.Get("/events/{id}", h.GetEvent)
		r.Post("/events", h.CreateEvent)
		r.Put("/events/{id}", h.UpdateEvent)
		r.Delete("/events/{id}", h.DeleteEvent)*/
	})
	return r
}
