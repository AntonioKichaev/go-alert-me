package handlers

import "net/http"

type ExecuteHandler interface {
	Register(server *http.ServeMux)
	GetParams(urlPath string) map[string]string
}
