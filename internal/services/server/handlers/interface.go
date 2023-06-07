package handlers

import (
	"github.com/go-chi/chi/v5"
)

type ExecuteHandler interface {
	Register(server *chi.Mux)
}
