package httprouter

import (
	"encoding/json"
	"github.com/concur/rohr"
	"github.com/concur/rohr/http/service"
	"github.com/julienschmidt/httprouter"
	"log"
	"net/http"
)

const (
	HEALTH_PATH = "/health"
)

type Router struct {
	httpRouter *httprouter.Router
}

func NewRouter() *Router {
	router := httprouter.New()
	return &Router{
		httpRouter: router,
	}
}

func (r *Router) RegisterRoute() http.Handler {
	r.httpRouter.GET(HEALTH_PATH, HealthHandler(service.NewHealthService()))
	return r.httpRouter
}

func HealthHandler(healthService rohr.HealthService) httprouter.Handle {
	return func(w http.ResponseWriter, _ *http.Request, _ httprouter.Params) {
		log.Printf("Invoke GetHealth API")

		data, err := healthService.GetHealth()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			log.Printf("GetHealth API returns error: %#v", err)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		enc := json.NewEncoder(w)
		enc.Encode(data)
		log.Printf("GetHealth API returns: %#v", data)
	}
}
