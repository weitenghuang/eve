package httprouter

import (
	"encoding/json"
	log "github.com/Sirupsen/logrus"
	"github.com/julienschmidt/httprouter"
	"github.com/scipian/eve"
	"github.com/scipian/eve/service"
	"net/http"
)

func getInfraHandler(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	user, err := getUser(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	infraSvc := service.NewInfrastructureService(user)

	log.Printf("Invoke GetInfrastructure API")
	name := p.ByName(P_NAME)
	infrastructure, err := infraSvc.GetInfrastructure(name)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Printf("GetInfrastructure API returns error: %#v", err)
		return
	}
	if infrastructure == nil {
		http.Error(w, RESOURCE_NOT_EXIST, http.StatusNotFound)
		log.Println("GetInfrastructure API returns: nil")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(infrastructure); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Printf("Encoding infrastructure returns error: %#v", err)
		return
	}
	log.Printf("GetInfrastructure API returns: %#v", struct {
		Id   string
		Name string
		*eve.Quoin
		eve.Status
		Variables []eve.QuoinVar
	}{
		infrastructure.Id,
		infrastructure.Name,
		infrastructure.Quoin,
		infrastructure.Status,
		infrastructure.Variables,
	})
}

func getInfraStateHandler(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	user, err := getUser(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	infraSvc := service.NewInfrastructureService(user)

	log.Printf("Invoke GetInfrastructureState API")
	name := p.ByName(P_NAME)
	state, err := infraSvc.GetInfrastructureState(name)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Printf("GetInfrastructureState API returns error: %#v", err)
		return
	}

	if len(state) == 0 {
		http.Error(w, RESOURCE_NOT_EXIST, http.StatusNotFound)
		log.Println("GetInfrastructureState API returns: nil")
		return
	} else {
		// Terraform store user credentials on remote state server. We should propose the change to terraform
		username, _, _ := r.BasicAuth()
		if username != "terraform" && state["remote"] != nil {
			delete(state, "remote")
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(state); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Printf("Encoding infrastructure state returns error: %#v", err)
		return
	}
	log.Printf("GetInfrastructureState API returns: %s\n", state)
}

func postInfraHandler(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	user, err := getUser(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	infraSvc := service.NewInfrastructureService(user)

	log.Println("Invoke CreateInfrastructure API")
	infrastructure, err := buildInfrastructure(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		log.Printf("buildInfrastructure returns error: %#v\n", err)
		return
	}
	if err := infraSvc.CreateInfrastructure(infrastructure); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Printf("CreateInfrastructure API returns error: %#v", err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)
	log.Printf("CreateInfrastructure API accepted request for %#v\n", infrastructure)
}

func postInfraStateHandler(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	user, err := getUser(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	infraSvc := service.NewInfrastructureService(user)

	log.Printf("Invoke UpdateInfrastructureState API")
	name := p.ByName(P_NAME)
	var state map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&state); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		log.Printf("Decode infrastructure state creation request returns error: %#v\n", err)
		return
	}
	if err := infraSvc.UpdateInfrastructureState(name, state); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Printf("UpdateInfrastructureState API returns error: %#v", err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	log.Printf("UpdateInfrastructureState API accepted request for %v\n", name)
}

func deleteInfraHandler(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	user, err := getUser(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	infraSvc := service.NewInfrastructureService(user)

	log.Printf("Invoke DeleteInfrastructure API")
	name := p.ByName(P_NAME)

	if err := infraSvc.DeleteInfrastructure(name); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Printf("DeleteInfrastructure API returns error: %#v", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)
	log.Printf("DeleteInfrastructure API accepted request for %#v\n", name)
}

func deleteInfraStateHandler(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	user, err := getUser(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	infraSvc := service.NewInfrastructureService(user)

	log.Printf("Invoke DeleteInfrastructureState API")
	name := p.ByName(P_NAME)
	if err := infraSvc.DeleteInfrastructureState(name); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Printf("DeleteInfrastructureState API returns error: %#v", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	log.Printf("DeleteInfrastructureState API accepted request for %#v\n", name)
}
