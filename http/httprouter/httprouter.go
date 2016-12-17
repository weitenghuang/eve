package httprouter

import (
	"fmt"
	"github.com/concur/rohr/http"
	"github.com/concur/rohr/service"
	"github.com/julienschmidt/httprouter"
	"log"
)

const (
	P_NAME      = "name"
	HEALTH_PATH = "/health"
	QUOIN_PATH  = "/quoin"
	INFRA_PATH  = "/infrastructure"
)

var (
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
	log.Println("Register route handlers for httprouter:")
	healthService := service.NewHealthService()
	quoinService := &service.QuoinService{}
	infrastructureService := &service.InfrastructureService{}
	r.httpRouter.GET(HEALTH_PATH, GetHealthHandler(healthService))
	log.Println("GET", HEALTH_PATH, "with GetHealthHandler")
	r.httpRouter.GET(QUOIN_NAME_PATH, GetQuoinHandler(quoinService))
	log.Println("GET", QUOIN_NAME_PATH, "with GetQuoinHandler")
	r.httpRouter.GET(INFRA_NAME_PATH, GetInfraHandler(infrastructureService))
	log.Println("GET", INFRA_NAME_PATH, "with GetInfraHandler")
	r.httpRouter.GET(INFRA_NAME_STATE_PATH, GetInfraStateHandler(infrastructureService))
	log.Println("GET", INFRA_NAME_STATE_PATH, "with GetInfraStateHandler")
	r.httpRouter.POST(QUOIN_PATH, PostQuoinHandler(quoinService, apiServer))
	log.Println("POST", QUOIN_PATH, "with PostQuoinHandler")
	r.httpRouter.POST(QUOIN_ARCHIVE_PATH, PostQuoinArchiveHandler(quoinService))
	log.Println("POST", QUOIN_ARCHIVE_PATH, "with PostQuoinArchiveHandler")
	r.httpRouter.POST(INFRA_PATH, PostInfraHandler(infrastructureService))
	log.Println("POST", INFRA_PATH, "with PostInfraHandler")
	r.httpRouter.POST(INFRA_NAME_STATE_PATH, PostInfraStateHandler(infrastructureService))
	log.Println("POST", INFRA_NAME_STATE_PATH, "with PostInfraStateHandler")
	r.httpRouter.DELETE(INFRA_NAME_PATH, DeleteInfraHandler(infrastructureService))
	log.Println("DELETE", INFRA_NAME_PATH, "with DeleteInfraHandler")
	r.httpRouter.DELETE(INFRA_NAME_STATE_PATH, DeleteInfraStateHandler(infrastructureService))
	log.Println("DELETE", INFRA_NAME_STATE_PATH, "with DeleteInfraStateHandler")
	return r.httpRouter
}
