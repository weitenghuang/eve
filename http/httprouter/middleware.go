package httprouter

import (
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/concur/eve/pkg/vault"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

type middleware func(next httprouter.Handle) httprouter.Handle

func mChain(route httprouter.Handle, chain ...middleware) httprouter.Handle {
	var h httprouter.Handle
	last := len(chain) - 1
	for i := last; i >= 0; i-- {
		m := chain[i]
		if i == last {
			h = m(route)
		} else {
			h = m(h)
		}
	}
	if h != nil {
		return h
	}
	return route
}

func logging(routeHandler httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		log.Infoln("logging is called")
		routeHandler(w, r, p)
	}
}

func authentication(routeHandler httprouter.Handle) httprouter.Handle {
	msg := "User authentication failure"
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		username, password, ok := r.BasicAuth()
		if !ok {
			log.Infof(msg)
			http.Error(w, msg, http.StatusUnauthorized)
			return
		}
		secretPath := fmt.Sprintf("secret/user/%s", username)
		user, err := vault.GetLogicalData(secretPath)
		if err != nil {
			log.Errorln(err)
		}
		if user["name"] != username || user["password"] != password {
			log.Infof("%s Name: %s, Password: %s", msg, username, password)
			http.Error(w, msg, http.StatusUnauthorized)
			return
		}
		log.Infoln("User logins:", username)
		r = setUserContext(r)
		routeHandler(w, r, p)
	}
}
