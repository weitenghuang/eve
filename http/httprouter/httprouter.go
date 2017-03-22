package httprouter

import (
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/concur/eve/http"
	"github.com/concur/eve/service"
	"github.com/julienschmidt/httprouter"
)

const (
	P_NAME        = "name"
	HEALTH_PATH   = "/health"
	PROVIDER_PATH = "/provider"
	QUOIN_PATH    = "/quoin"
	INFRA_PATH    = "/infrastructure"
)

var (
	PROVIDER_NAME_PATH    string = fmt.Sprintf("%s/:%s", PROVIDER_PATH, P_NAME)
	QUOIN_NAME_PATH       string = fmt.Sprintf("%s/:%s", QUOIN_PATH, P_NAME)
	QUOIN_ARCHIVE_PATH    string = fmt.Sprintf("%s/upload", QUOIN_NAME_PATH)
	INFRA_NAME_PATH       string = fmt.Sprintf("%s/:%s", INFRA_PATH, P_NAME)
	INFRA_NAME_STATE_PATH string = fmt.Sprintf("%s/state", INFRA_NAME_PATH)
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

func (r *Router) RegisterRoute(apiServer *http.ApiServer) http.Handler {
	if r.httpRouter == nil {
		log.Panicln("API server's router misses httprouter field value.")
	}

	log.Infoln("Register route handlers for httprouter:")
	healthService := service.NewHealthService()
	r.httpRouter.GET(HEALTH_PATH, mChain(getHealthHandler(healthService)))
	log.Infoln("GET", HEALTH_PATH, "with getHealthHandler")
	r.httpRouter.GET(PROVIDER_NAME_PATH, mChain(getProviderHandler, authentication))
	log.Infoln("GET", PROVIDER_NAME_PATH, "with GetProviderHandler")
	r.httpRouter.GET(QUOIN_NAME_PATH, mChain(getQuoinHandler, logging, authentication))
	log.Infoln("GET", QUOIN_NAME_PATH, "with getQuoinHandler")
	r.httpRouter.GET(INFRA_NAME_PATH, mChain(getInfraHandler, authentication))
	log.Infoln("GET", INFRA_NAME_PATH, "with getInfraHandler")
	r.httpRouter.GET(INFRA_NAME_STATE_PATH, mChain(getInfraStateHandler, authentication))
	log.Infoln("GET", INFRA_NAME_STATE_PATH, "with getInfraStateHandler")
	r.httpRouter.POST(QUOIN_PATH, mChain(postQuoinHandler(apiServer), authentication))
	log.Infoln("POST", QUOIN_PATH, "with postQuoinHandler")
	r.httpRouter.POST(QUOIN_ARCHIVE_PATH, mChain(postQuoinArchiveHandler, authentication))
	log.Infoln("POST", QUOIN_ARCHIVE_PATH, "with postQuoinArchiveHandler")
	r.httpRouter.POST(INFRA_PATH, mChain(postInfraHandler, authentication))
	log.Infoln("POST", INFRA_PATH, "with postInfraHandler")
	r.httpRouter.POST(INFRA_NAME_STATE_PATH, mChain(postInfraStateHandler, authentication))
	log.Infoln("POST", INFRA_NAME_STATE_PATH, "with postInfraStateHandler")
	r.httpRouter.DELETE(INFRA_NAME_PATH, mChain(deleteInfraHandler, authentication))
	log.Infoln("DELETE", INFRA_NAME_PATH, "with deleteInfraHandler")
	r.httpRouter.DELETE(INFRA_NAME_STATE_PATH, mChain(deleteInfraStateHandler, authentication))
	log.Infoln("DELETE", INFRA_NAME_STATE_PATH, "with deleteInfraStateHandler")
	return r.httpRouter
}
