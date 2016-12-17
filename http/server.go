package http

import (
	"fmt"
	"log"
	"net/http"
	"strings"
)

type ApiServer struct {
	Addr   string // TCP address to listen on, ":http" if empty
	DNS    string
	Scheme string
	Router
}

type Router interface {
	RegisterRoute(apiServer *ApiServer) Handler
}

type Handler interface {
	http.Handler
}

func (s *ApiServer) ListenAndServe() {
	if !strings.ContainsAny(s.Addr, ":") {
		s.Addr = ":http"
	}
	log.Printf("start listening on %s%s", s.DNS, s.Addr)
	server := &http.Server{Addr: s.Addr, Handler: s.Router.RegisterRoute(s)}
	log.Fatal(server.ListenAndServe())
}

func (s *ApiServer) GetServerAddr() string {
	return fmt.Sprintf("%s://%s%s/", s.Scheme, s.DNS, s.Addr)
}
