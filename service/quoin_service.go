package service

import (
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/concur/eve"
	"github.com/concur/eve/pkg/terraform"
	"github.com/concur/eve/service/rethinkdb"
	"strings"
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
	err := db.InsertQuoin(quoin)
	if err != nil {
		return quoin, err
	}
	return quoin, nil
}

// CreateQuoinArchive creates Quoin archive record on database
func (q QuoinService) CreateQuoinArchive(quoinArchive *eve.QuoinArchive) error {
	tf := terraform.NewTerraform(quoinArchive.QuoinName, "", quoinArchive.Modules, nil)
	if err := tf.PlanQuoin(); err != nil {
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
	// Can not delete the Quoin if it is still in used by an infratructure
	return nil
}

func (q QuoinService) DeleteQuoinArchive(id string) error {
	// Can not delete the Quoin if it is still in used by a quoin
	return nil
}
