package http

import (
	"fmt"
	log "github.com/Sirupsen/logrus"
	"net/http"
	"strings"
)

const (
	VAULT      = "vault"
	HTTPROUTER = "httprouter"
)

type ApiServer struct {
	Addr     string // TCP address to listen on, ":http" if empty
	DNS      string
	Scheme   string
	CertFile string
	KeyFile  string
	Routers  map[string]Router
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

	routerSwitch := make(RouterSwitch)
	for key, route := range s.Routers {
		routerSwitch[key] = route.RegisterRoute(s)
	}
	server := &http.Server{Addr: s.Addr, Handler: routerSwitch}

	switch strings.ToLower(s.Scheme) {
	case "http":
		log.Fatal(server.ListenAndServe())
	case "https":
		log.Fatal(server.ListenAndServeTLS(s.CertFile, s.KeyFile))
	default:
		log.Fatal("Invalid scheme vaule: %v", s.Scheme)
	}
}

func (s *ApiServer) GetServerAddr() string {
	return fmt.Sprintf("%s://%s%s/", s.Scheme, s.DNS, s.Addr)
}

type RouterSwitch map[string]Handler

// Implement the ServerHTTP method on our new type
func (rs RouterSwitch) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Use Vault reverse proxy router for endpoint: `/vault`
	if len(r.URL.Path) > len(VAULT) && VAULT == r.URL.Path[1:len(VAULT)+1] {
		rs[VAULT].ServeHTTP(w, r)
	} else { // Use httprouter for all other endpoints
		rs[HTTPROUTER].ServeHTTP(w, r)
	}
}
