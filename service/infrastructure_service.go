package service

import (
	"fmt"

	log "github.com/Sirupsen/logrus"
	"github.com/scipian/eve"
	"github.com/scipian/eve/service/nats"
	"github.com/scipian/eve/service/rethinkdb"
)

type InfrastructureService struct {
	*eve.User
}

const (
	INVALID_QUOIN_ERROR = "To create an infrastructure, please use a valid quoin"
)

func NewInfrastructureService(user *eve.User) *InfrastructureService {
	return &InfrastructureService{
		User: user,
	}
}

func (infraSvc InfrastructureService) GetInfrastructure(name string) (*eve.Infrastructure, error) {
	log.Infoln("Get Infrastructure for user:", infraSvc.User)
	db := rethinkdb.DefaultSession()
	infrastructure, err := db.GetInfrastructureByName(name)
	if err != nil {
		return nil, err
	}

	if infrastructure == nil {
		return nil, nil
	}

	if !infrastructure.AuthorizedRead(infraSvc.User) {
		return nil, fmt.Errorf("User %s is not authorized to read infrastructure %s", infraSvc.User.Id, infrastructure.Name)
	}

	return infrastructure, nil
}

func (infraSvc InfrastructureService) GetInfrastructuresByQuoin(quoinName string) ([]eve.Infrastructure, error) {
	infras, err := infraSvc.getInfrastructuresByQuoin(quoinName)
	if err != nil {
		return nil, err
	}

	// If user is not authorized to read the infrastructure, we will only return the infrastructure's name
	for i, infra := range infras {
		if !infra.AuthorizedRead(infraSvc.User) {
			infras[i] = eve.Infrastructure{Name: infra.Name}
		}
	}
	return infras, err
}

func (infraSvc InfrastructureService) CountInfrastructureByQuoin(quoinName string) (int, error) {
	infras, err := infraSvc.getInfrastructuresByQuoin(quoinName)
	if err != nil {
		return 0, err
	}
	return len(infras), nil
}

func (infraSvc InfrastructureService) getInfrastructuresByQuoin(quoinName string) ([]eve.Infrastructure, error) {
	db := rethinkdb.DefaultSession()
	infrastructures, err := db.GetInfrastructuresByQuoin(quoinName)
	if err != nil {
		return nil, err
	}
	return infrastructures, nil
}

func (infraSvc InfrastructureService) GetInfrastructureState(name string) (map[string]interface{}, error) {
	infra, err := infraSvc.GetInfrastructure(name)
	if err != nil {
		return nil, err
	}

	if infra == nil {
		return nil, nil
	}

	return infra.State, nil
}

func (infraSvc InfrastructureService) CreateInfrastructure(infra *eve.Infrastructure) error {
	searchResult, err := infraSvc.GetInfrastructure(infra.Name)
	if err != nil {
		return err
	}

	if searchResult != nil {
		log.Printf("Found existing infrastructure %s.\n", infra.Name)
		switch searchResult.Status {
		case eve.RUNNING, eve.DEPLOYED, eve.OBSOLETED:
			return fmt.Errorf("Infrastructure %s cannot be created at this moment. Please check its current status first.", infra.Name)
		default:
			log.Printf("Re-create existing infrastructure %s.\n", infra.Name)
		}
	} else {
		// Validate infrastructure's quoin reference
		quoinSvc := NewQuoinService(infraSvc.User)
		quoin, err := quoinSvc.GetQuoin(infra.Quoin.Name)
		if err != nil {
			return err
		}

		if quoin == nil {
			return fmt.Errorf(INVALID_QUOIN_ERROR)
		}

		if quoin.Status != eve.VALIDATED {
			return fmt.Errorf(INVALID_QUOIN_ERROR)
		}

		// We allow user to use older version of quoin's archive
		if infra.Quoin.ArchiveUri != quoin.ArchiveUri {
			log.Infof("Current infrastructure request's %#v quoin %#v doesn't have the latest quoin archive reference %#v.", infra, infra.Quoin, quoin.ArchiveUri)
		}

		archiveId := quoinSvc.GetQuoinArchiveIdFromUri(infra.Quoin.ArchiveUri)
		if archiveId == "" {
			return fmt.Errorf(INVALID_QUOIN_ERROR)
		}

		infra.Status = eve.VALIDATED

		db := rethinkdb.DefaultSession()
		if err := db.InsertInfrastructure(infra); err != nil {
			return err
		}
		log.Printf("New infrastructure %s is stored in eve db.\n", infra.Name)
	}

	if err := infraSvc.PublishMessageToQueue(eve.CREATE_INFRA, infra); err != nil {
		return err
	}

	return nil
}

func (infraSvc InfrastructureService) DeleteInfrastructure(name string) error {
	if err := infraSvc.checkWritePermission(name); err != nil {
		return err
	}

	db := rethinkdb.DefaultSession()
	infra, err := db.GetInfrastructureByName(name)
	if err != nil {
		return err
	}

	if len(infra.State) == 0 {
		return fmt.Errorf("Infrastructure %s's state is missing", name)
	}

	if infra.Status == eve.RUNNING {
		return fmt.Errorf("Infrastructure %s's deletion is in process", name)
	}

	// Avoid NATS queue message size limit
	infra.State = nil

	if err := infraSvc.PublishMessageToQueue(eve.DELETE_INFRA, infra); err != nil {
		return err
	}

	return nil
}

func (infraSvc InfrastructureService) DeleteInfrastructureState(name string) error {
	return nil
}

func (infraSvc InfrastructureService) UpdateInfrastructureState(name string, state map[string]interface{}) error {
	if err := infraSvc.checkWritePermission(name); err != nil {
		return err
	}

	db := rethinkdb.DefaultSession()
	if err := db.UpdateInfrastructureState(name, state); err != nil {
		return err
	}

	return nil
}

func (infraSvc InfrastructureService) UpdateInfrastructureStatus(name string, status eve.Status) error {
	if err := infraSvc.checkWritePermission(name); err != nil {
		return err
	}

	db := rethinkdb.DefaultSession()
	if err := db.UpdateInfrastructureStatus(name, status); err != nil {
		return err
	}

	return nil
}

func (infraSvc InfrastructureService) UpdateInfrastructureError(name string, infraError error) error {
	if err := infraSvc.checkWritePermission(name); err != nil {
		return err
	}

	db := rethinkdb.DefaultSession()
	if err := db.UpdateInfrastructureError(name, infraError); err != nil {
		return err
	}

	return nil
}

func (infraSvc InfrastructureService) SubscribeAsyncProc(subject eve.Subject, handler eve.InfrastructureAsyncHandler) error {
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

func (infraSvc InfrastructureService) PublishMessageToQueue(subject eve.Subject, infra *eve.Infrastructure) error {
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

func (infraSvc InfrastructureService) checkWritePermission(name string) error {
	db := rethinkdb.DefaultSession()
	infra, err := db.GetInfrastructureByName(name)
	if err != nil {
		return err
	}

	if infra == nil {
		return fmt.Errorf("Infrastructure %s not found", name)
	}

	if !infra.AuthorizedWrite(infraSvc.User) {
		return fmt.Errorf("User %s is not authorized to modify infrastructure %s", infraSvc.User.Id, infra.Name)
	}
	return nil
}

func (infraSvc InfrastructureService) GetUser() *eve.User {
	return infraSvc.User
}
