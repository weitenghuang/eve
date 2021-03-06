package httprouter

import (
	"encoding/json"
	log "github.com/Sirupsen/logrus"
	"github.com/concur/eve"
	eveHttp "github.com/concur/eve/http"
	"github.com/concur/eve/service"
	"github.com/julienschmidt/httprouter"
	"io/ioutil"
	"net/http"
)

func getQuoinHandler(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	user, err := getUser(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	quoinService := service.NewQuoinService(user)

	log.Printf("Invoke GetQuoin API")
	name := p.ByName(P_NAME)
	quoin, err := quoinService.GetQuoin(name)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Printf("GetQuoin API returns error: %#v", err)
		return
	}

	if quoin == nil {
		http.Error(w, RESOURCE_NOT_EXIST, http.StatusNotFound)
		log.Println("GetQuoin API returns: nil")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(quoin); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Printf("Encoding quoin returns error: %#v", err)
		return
	}
	log.Printf("GetQuoin API returns: %#v", quoin)
}

// postQuoinHandler returns the httprouter.Handle func for POST /quoin request
func postQuoinHandler(apiServer *eveHttp.ApiServer) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		user, err := getUser(r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}
		quoinService := service.NewQuoinService(user)

		log.Println("Invoke CreateQuoin API")
		quoinInput, err := buildQuoin(r, apiServer)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			log.Printf("buildQuoin returns error: %#v", err)
			return
		}

		quoin, err := quoinService.CreateQuoin(quoinInput)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			log.Printf("CreateQuoin API returns error: %#v", err)
			return
		}

		// API designs resumable upload:
		// https://developers.google.com/drive/v3/web/manage-uploads#resumable
		if quoin.ArchiveUri != "" {
			w.Header().Set("Location", quoin.ArchiveUri)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusAccepted)
		if err := json.NewEncoder(w).Encode(quoin); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			log.Printf("Encoding quoin returns error: %#v", err)
			return
		}
		log.Printf("CreateQuoin API returns: %#v\n", quoin)
	}
}

func postQuoinArchiveHandler(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	user, err := getUser(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	quoinService := service.NewQuoinService(user)

	log.Printf("Invoke CreateQuoinArchive API")
	name := p.ByName(P_NAME)
	content, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		log.Printf("Bad Request with invalid body. Error: %#v", err)
	}
	quoinArchive := &eve.QuoinArchive{
		QuoinName: name,
		Modules:   content,
	}
	bindAuthorization(quoinArchive, r)
	if err := quoinService.CreateQuoinArchive(quoinArchive); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Printf("CreateQuoinArchive API returns error: %#v", err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)
	json.NewEncoder(w).Encode(struct {
		Id        string
		QuoinName string
	}{
		quoinArchive.Id,
		quoinArchive.QuoinName,
	})
	log.Printf("CreateQuoinArchive API returns: %#v", quoinArchive.Id)
}

func deleteQuoinHandler(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	user, err := getUser(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	quoinService := service.NewQuoinService(user)

	log.Printf("Invoke DeleteQuoin API")
	name := p.ByName(P_NAME)
	if err := quoinService.DeleteQuoin(name); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Printf("DeleteQuoin API returns error: %#v", err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	log.Printf("DeleteQuoin API completed: %v", name)
}
