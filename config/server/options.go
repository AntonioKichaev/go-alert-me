package server

type Option func(server *Server)

func SetHTTPServerAdr(adr string) Option {
	return func(server *Server) {
		server.HTTPServerAdr = adr
	}
}
