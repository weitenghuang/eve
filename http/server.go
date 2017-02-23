package http

import (
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/concur/rohr/pkg/vault"
	"net/http"
	"strings"
)

type ApiServer struct {
	Addr     string // TCP address to listen on, ":http" if empty
	DNS      string
	Scheme   string
	CertFile string
	KeyFile  string
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

	handler := s.Router.RegisterRoute(s)
	server := &http.Server{Addr: s.Addr, Handler: handler}

	switch strings.ToLower(s.Scheme) {
	case "http":
		log.Fatal(server.ListenAndServe())
	case "https":
		server.Handler = Auth(server.Handler)
		log.Fatal(server.ListenAndServeTLS(s.CertFile, s.KeyFile))
	default:
		log.Fatal("Invalid scheme vaule: %v", s.Scheme)
	}
}

func (s *ApiServer) GetServerAddr() string {
	return fmt.Sprintf("%s://%s%s/", s.Scheme, s.DNS, s.Addr)
}

func Auth(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		username, password, _ := r.BasicAuth()
		secretPath := fmt.Sprintf("secret/user/%s", username)
		user, err := vault.GetLogicalData(secretPath)
		if err != nil {
			log.Errorln(err)
		}
		if user["name"] != username || user["password"] != password {
			msg := "User authentication failure"
			log.Infof("%s Name: %s, Password: %s", msg, username, password)
			http.Error(w, msg, http.StatusUnauthorized)
			return
		}
		log.Infoln("User logins:", username)
		handler.ServeHTTP(w, r)
	})
}
