package httprouter

import (
	"encoding/json"
	log "github.com/Sirupsen/logrus"
	"github.com/concur/rohr"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

func GetHealthHandler(healthService rohr.HealthService) httprouter.Handle {
	return func(w http.ResponseWriter, _ *http.Request, _ httprouter.Params) {
		log.Printf("Invoke GetHealth API")

		data := healthService.GetHealth()

		w.Header().Set("Content-Type", "application/json")
		if len(data.Errors) > 0 {
			w.WriteHeader(http.StatusServiceUnavailable)
		} else {
			w.WriteHeader(http.StatusOK)
		}
		if err := json.NewEncoder(w).Encode(data); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			log.Printf("Encoding health data returns error: %#v", err)
			return
		}

		log.Printf("GetHealth API returns: %#v", data)
	}
}
