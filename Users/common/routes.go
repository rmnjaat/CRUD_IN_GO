package common

import "github.com/go-chi/chi"

func RegisterRoutes(r chi.Router) {
	r.Get("/healthz", healthCheckhandler)
}
