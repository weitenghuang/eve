package httprouter

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/concur/eve"
	eveHttp "github.com/concur/eve/http"
	"net/http"
)

const (
	CTX_USER = "user"
)

func buildQuoin(r *http.Request, apiServer *eveHttp.ApiServer) (*eve.Quoin, error) {
	if r.Body == nil {
		return nil, fmt.Errorf("Empty request body is invalid for POST /quoin request")
	}
	var quoin eve.Quoin
	err := json.NewDecoder(r.Body).Decode(&quoin)
	if err != nil {
		return nil, err
	}
	quoin.ArchiveUri = fmt.Sprintf("%s://%s%s/quoin/%s/upload", apiServer.Scheme, apiServer.DNS, apiServer.Addr, quoin.Name)

	bindAuthorization(&quoin, r)
	return &quoin, nil
}

func buildInfrastructure(r *http.Request) (*eve.Infrastructure, error) {
	if r.Body == nil {
		return nil, fmt.Errorf("Empty request body is invalid for POST /infrastructure request")
	}
	var infrastructure eve.Infrastructure
	err := json.NewDecoder(r.Body).Decode(&infrastructure)
	if err != nil {
		return nil, err
	}

	bindAuthorization(&infrastructure, r)
	return &infrastructure, nil
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

func getUser(r *http.Request) (*eve.User, error) {
	user := r.Context().Value(CTX_USER).(*eve.User)
	if user == nil {
		return nil, fmt.Errorf("User information is missing.")
	}
	return user, nil
}

func bindAuthorization(resource eve.Authorizable, r *http.Request) (eve.Authorizable, error) {
	user, err := getUser(r)
	if err != nil {
		return nil, err
	}

	auth := eve.Authorization{
		Owner: user.Id,
		GroupAccess: map[eve.Group]eve.PolicyMode{
			eve.Group(user.Id):           eve.POLICY_ALL,
			eve.Group(user.Organization): eve.POLICY_READ,
			eve.Group("public"):          eve.POLICY_NONE,
		},
	}
	resource.BindAuthorization(auth)
	return resource, nil
}
