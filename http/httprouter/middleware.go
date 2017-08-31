package httprouter

import (
	log "github.com/Sirupsen/logrus"
	"github.com/julienschmidt/httprouter"
	eveHttp "github.com/scipian/eve/http"
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
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		var err error
		r, err = eveHttp.Authentication(w, r)
		if err != nil {
			log.Infoln("Unauthorized request: ", err)
		}
		routeHandler(w, r, p)
	}
}
