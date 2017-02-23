package service

import (
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/concur/rohr"
	"github.com/concur/rohr/service/nats"
	"github.com/concur/rohr/service/rethinkdb"
)

type InfrastructureService struct {
	QuoinService
}

func (infraSvc InfrastructureService) GetInfrastructure(name string) (*rohr.Infrastructure, error) {
	db := rethinkdb.DefaultSession()
	infrastructure, err := db.GetInfrastructureByName(name)
	if err != nil {
		return nil, err
	}
	return infrastructure, nil
}

func (infraSvc InfrastructureService) GetInfrastructureState(name string) (map[string]interface{}, error) {
	db := rethinkdb.DefaultSession()
	state, err := db.GetInfrastructureStateByName(name)
	if err != nil {
		return nil, err
	}
	return state, nil
}

func (infraSvc InfrastructureService) CreateInfrastructure(infra *rohr.Infrastructure) error {
	searchResult, _ := infraSvc.GetInfrastructure(infra.Name)
	if searchResult != nil {
		log.Printf("Found existing infrastructure %s.\n", infra.Name)
		switch searchResult.Status {
		case rohr.RUNNING, rohr.DEPLOYED, rohr.OBSOLETED:
			return fmt.Errorf("Infrastructure %s cannot be created at this moment. Please check its current status first.", infra.Name)
		default:
			log.Printf("Re-create existing infrastructure %s.\n", infra.Name)
		}
	} else {
		archiveId := infraSvc.GetQuoinArchiveIdFromUri(infra.Quoin.ArchiveUri)
		if archiveId != "" {
			infra.Status = rohr.VALIDATED
		} else {
			return fmt.Errorf("Empty archive id error: To create an infrastructure, please provide valid quoin archive id.")
		}
		db := rethinkdb.DefaultSession()
		if err := db.InsertInfrastructure(infra); err != nil {
			return err
		}
		log.Printf("New infrastructure %s is stored in eve db.\n", infra.Name)
	}

	if err := infraSvc.PublishMessageToQueue(rohr.CREATE_INFRA, infra); err != nil {
		return err
	}

	return nil
}

func (infraSvc InfrastructureService) DeleteInfrastructure(name string) error {
	infra, err := infraSvc.GetInfrastructure(name)
	if err != nil {
		return err
	}
	if infra == nil {
		return fmt.Errorf("Infrastructure %s not found", name)
	}
	if len(infra.State) == 0 {
		return fmt.Errorf("Infrastructure %s's state is missing", name)
	}
	if infra.Status == rohr.RUNNING {
		return fmt.Errorf("Infrastructure %s's deletion is in process", name)
	}

	// Avoid NATS queue message size limit
	infra.State = nil

	if err := infraSvc.PublishMessageToQueue(rohr.DELETE_INFRA, infra); err != nil {
		return err
	}

	return nil
}

func (infraSvc InfrastructureService) DeleteInfrastructureState(name string) error {
	return nil
}

func (infraSvc InfrastructureService) UpdateInfrastructureState(name string, state map[string]interface{}) error {
	db := rethinkdb.DefaultSession()
	if err := db.UpdateInfrastructureState(name, state); err != nil {
		return err
	}
	return nil
}

func (infraSvc InfrastructureService) UpdateInfrastructureStatus(name string, status rohr.Status) error {
	db := rethinkdb.DefaultSession()
	if err := db.UpdateInfrastructureStatus(name, status); err != nil {
		return err
	}
	return nil
}

func (infraSvc InfrastructureService) SubscribeAsyncProc(subject rohr.Subject, handler rohr.InfrastructureAsyncHandler) error {
	c, err := nats.EncodedConn()
	if err != nil {
		log.Println(err)
		return err
	}
	// defer c.Close()
	// To close connection by runtime.Goexit()
	subject_s := string(subject)
	if _, err := c.QueueSubscribe(subject_s, subject_s, handler); err != nil {
		return err
	}
	return nil
}

func (infraSvc InfrastructureService) PublishMessageToQueue(subject rohr.Subject, infra *rohr.Infrastructure) error {
	c, err := nats.EncodedConn()
	if err != nil {
		log.Println(err)
		return err
	}
	defer c.Close()
	subject_s := string(subject)
	if err := c.Publish(subject_s, infra); err != nil {
		log.Println(err)
		return err
	}
	log.Printf("Publish Infrastructure %s to infra queue %s.\n", infra.Name, subject_s)
	return nil
}
