package http

import (
	"log"
	"net/http"
	"strings"
)

type ApiServer struct {
	Addr string // TCP address to listen on, ":http" if empty
	DNS  string
	Router
}

type Router interface {
	RegisterRoute() http.Handler
}

func (s *ApiServer) ListenAndServe() {
	if !strings.ContainsAny(s.Addr, ":") {
		s.Addr = ":http"
	}
	log.Printf("start listening on %s%s", s.DNS, s.Addr)
	server := &http.Server{Addr: s.Addr, Handler: s.Router.RegisterRoute()}
	log.Fatal(server.ListenAndServe())
}
