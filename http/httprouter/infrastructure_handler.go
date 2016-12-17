package httprouter

import (
	"encoding/json"
	"fmt"
	"github.com/concur/rohr"
	"github.com/julienschmidt/httprouter"
	"log"
	"net/http"
)

func GetInfraHandler(infraSvc rohr.InfrastructureService) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		log.Printf("Invoke GetInfrastructure API")
		name := p.ByName(P_NAME)
		infrastructure, err := infraSvc.GetInfrastructure(name)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			log.Printf("GetInfrastructure API returns error: %#v", err)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if infrastructure == nil {
			log.Println("GetInfrastructure API returns: nil")
			return
		}
		if err := json.NewEncoder(w).Encode(infrastructure); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			log.Printf("Encoding infrastructure returns error: %#v", err)
			return
		}
		log.Printf("GetInfrastructure API returns: %#v", struct {
			Id   string
			Name string
			*rohr.Quoin
			rohr.Status
			Variables []rohr.QuoinVar
		}{
			infrastructure.Id,
			infrastructure.Name,
			infrastructure.Quoin,
			infrastructure.Status,
			infrastructure.Variables,
		})
	}
}

func GetInfraStateHandler(infraSvc rohr.InfrastructureService) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		log.Printf("Invoke GetInfrastructureState API")
		name := p.ByName(P_NAME)
		state, err := infraSvc.GetInfrastructureState(name)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			log.Printf("GetInfrastructureState API returns error: %#v", err)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if len(state) == 0 {
			return
		} else if err := json.NewEncoder(w).Encode(state); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			log.Printf("Encoding infrastructure state returns error: %#v", err)
			return
		}
		log.Printf("GetInfrastructureState API returns: %s\n", state)
	}
}

func PostInfraHandler(infraSvc rohr.InfrastructureService) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
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
}

func PostInfraStateHandler(infraSvc rohr.InfrastructureService) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
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
}

func DeleteInfraHandler(infraSvc rohr.InfrastructureService) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
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
}

func DeleteInfraStateHandler(infraSvc rohr.InfrastructureService) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
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
}

func buildInfrastructure(r *http.Request) (*rohr.Infrastructure, error) {
	if r.Body == nil {
		return nil, fmt.Errorf("Empty request body is invalid for POST /infrastructure request")
	}
	var infrastructure rohr.Infrastructure
	err := json.NewDecoder(r.Body).Decode(&infrastructure)
	if err != nil {
		return nil, err
	}
	return &infrastructure, nil
}
