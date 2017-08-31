package service

import (
	"fmt"
	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/scipian/eve"
	"github.com/scipian/eve/pkg/terraform"
	"github.com/scipian/eve/service/rethinkdb"
)

type QuoinService struct {
	*eve.User
}

func NewQuoinService(user *eve.User) *QuoinService {
	return &QuoinService{
		User: user,
	}
}

// GetQuoin returns Quoin information from database
func (q QuoinService) GetQuoin(name string) (*eve.Quoin, error) {
	db := rethinkdb.DefaultSession()
	quoin, err := db.GetQuoinByName(name)
	if err != nil {
		return nil, err
	}

	if quoin == nil {
		return nil, nil
	}

	if !quoin.AuthorizedRead(q.User) {
		return nil, fmt.Errorf("User %s is not authorized to read Quoin %s", q.User.Id, quoin.Name)
	}

	return quoin, nil
}

// GetQuoinArchive returns Quoin archive module from database
func (q QuoinService) GetQuoinArchive(id string) (*eve.QuoinArchive, error) {
	db := rethinkdb.DefaultSession()
	quoinArchive, err := db.GetQuoinArchiveById(id)
	if err != nil {
		return nil, err
	}

	if quoinArchive == nil {
		return nil, nil
	}

	if !quoinArchive.AuthorizedRead(q.User) {
		return nil, fmt.Errorf("User %s is not authorized to read Quoin Archive %s", q.User.Id, quoinArchive.Id)
	}
	return quoinArchive, nil
}

// GetQuoinArchives returns list of quoin archive modules from database
func (q QuoinService) GetQuoinArchiveIds(quoinName string) ([]string, error) {
	return []string{}, nil
}

func (q QuoinService) GetQuoinArchiveIdFromUri(archiveUri string) string {
	uri := strings.SplitAfter(archiveUri, "/upload/")
	if len(uri) > 1 {
		return uri[1]
	} else {
		return ""
	}
}

// CreateQuoin creates Quoin record on database and calls CreateQuoinArchive
func (q QuoinService) CreateQuoin(quoin *eve.Quoin) (*eve.Quoin, error) {
	db := rethinkdb.DefaultSession()
	qu, err := db.GetQuoinByName(quoin.Name)
	if err != nil {
		return nil, err
	}

	if qu != nil {
		if qu.Status != eve.OBSOLETED {
			return qu, fmt.Errorf("Quoin %s already exists", quoin.Name)
		} else {
			if !qu.AuthorizedWrite(q.User) {
				return nil, fmt.Errorf("User %s is not authorized to re-create Quoin %s", q.User.Id, quoin.Name)
			}
			if err := db.UpdateQuoin(quoin.Name, quoin); err != nil {
				return nil, err
			}
			return quoin, nil
		}
	}

	if err := db.InsertQuoin(quoin); err != nil {
		return nil, err
	}
	return quoin, nil
}

// CreateQuoinArchive creates Quoin archive record on database
func (q QuoinService) CreateQuoinArchive(quoinArchive *eve.QuoinArchive) error {
	tf := terraform.NewTerraform(quoinArchive.QuoinName, "", quoinArchive.Modules, nil)
	if err := tf.ValidateQuoin(); err != nil {
		return err
	}
	log.Printf("Quoin Archive for %s is valid. Terraform plan has been generated.", quoinArchive.QuoinName)
	db := rethinkdb.DefaultSession()
	if err := db.InsertQuoinArchive(quoinArchive); err != nil {
		return err
	}
	log.Printf("Quoin Archive for %s is stored in eve db with id: %s", quoinArchive.QuoinName, quoinArchive.Id)
	return nil
}

func (q QuoinService) DeleteQuoin(name string) error {
	db := rethinkdb.DefaultSession()
	quoin, err := db.GetQuoinByName(name)
	if err != nil {
		return err
	}

	if quoin == nil {
		return fmt.Errorf("Quoin %s doesn't exist", name)
	}

	if !quoin.AuthorizedWrite(q.User) {
		return fmt.Errorf("User %s is not authorized to delete Quoin %s", q.User.Id, name)
	}

	infraSvc := NewInfrastructureService(q.User)
	count, err := infraSvc.CountInfrastructureByQuoin(name)
	if err != nil {
		return err
	}
	if count > 0 {
		return fmt.Errorf("Quoin is still used by infrastructure and cannot be deleted")
	}

	if quoin.Status != eve.OBSOLETED {
		quoin.Status = eve.OBSOLETED
		if err := db.UpdateQuoin(name, quoin); err != nil {
			return err
		}
	}

	return nil
}

func (q QuoinService) DeleteQuoinArchive(id string) error {
	// Can not delete the Quoin if it is still in used by a quoin
	return nil
}
