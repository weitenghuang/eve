package httprouter

import (
	"encoding/json"
	"net/http"

	log "github.com/Sirupsen/logrus"
	"github.com/concur/eve/service"
	"github.com/julienschmidt/httprouter"
)

// GetProviderHandler returns the provider with the given name
func getProviderHandler(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	user, err := getUser(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	providerService := service.NewProviderService(user)

	log.Printf("Invoke GetProvider API")
	name := p.ByName(P_NAME)
	provider, err := providerService.GetProvider(name)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Printf("GetProvider API returns error: %#v", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if provider == nil {
		log.Println("GetProvider API returns: nil")
		return
	}
	if err := json.NewEncoder(w).Encode(provider); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Printf("Encoding provider returns error: %#v", err)
		return
	}
	log.Printf("GetProvider API returns: %#v", provider)
}
