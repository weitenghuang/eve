package http

import (
	"context"
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/scipian/eve"
	"github.com/scipian/eve/pkg/vault"
	"net/http"
)

const (
	CTX_USER               = "user"
	AUTHENTICATION_FAILURE = "User authentication failure"
)

func Authentication(w http.ResponseWriter, r *http.Request) (*http.Request, error) {
	username, password, ok := r.BasicAuth()
	if !ok {
		log.Infof(AUTHENTICATION_FAILURE)
		http.Error(w, AUTHENTICATION_FAILURE, http.StatusUnauthorized)
		return nil, fmt.Errorf("%s:%s", AUTHENTICATION_FAILURE, username)
	}
	secretPath := fmt.Sprintf("secret/user/%s", username)
	user, err := vault.GetLogicalData(secretPath)
	if err != nil {
		return nil, err
	}
	if user["name"] != username || user["password"] != password {
		log.Infof("%s Name: %s, Password: %s", AUTHENTICATION_FAILURE, username, password)
		http.Error(w, AUTHENTICATION_FAILURE, http.StatusUnauthorized)
		return nil, fmt.Errorf("%s:%s", AUTHENTICATION_FAILURE, username)
	}
	log.Infoln("User logins:", username)
	r = setUserContext(r)
	return r, nil
}

func setUserContext(r *http.Request) *http.Request {
	username, _, _ := r.BasicAuth()

	// TODO(weiteng.huang): User creation/retrieval will be from User service
	ctx := context.WithValue(r.Context(), CTX_USER, &eve.User{
		Id:           eve.UserId(username),
		Organization: eve.Organization("concur"),
		Teams:        []eve.Team{eve.Team(username)},
	})
	return r.WithContext(ctx)
}
