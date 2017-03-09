package httprouter

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/concur/rohr"
	rohrHttp "github.com/concur/rohr/http"
	"net/http"
)

const (
	CTX_USER = "user"
)

func buildQuoin(r *http.Request, apiServer *rohrHttp.ApiServer) (*rohr.Quoin, error) {
	if r.Body == nil {
		return nil, fmt.Errorf("Empty request body is invalid for POST /quoin request")
	}
	var quoin rohr.Quoin
	err := json.NewDecoder(r.Body).Decode(&quoin)
	if err != nil {
		return nil, err
	}
	quoin.ArchiveUri = fmt.Sprintf("%s://%s%s/quoin/%s/upload", apiServer.Scheme, apiServer.DNS, apiServer.Addr, quoin.Name)

	bindAuthorization(&quoin, r)
	return &quoin, nil
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

	bindAuthorization(&infrastructure, r)
	return &infrastructure, nil
}

func setUserContext(r *http.Request) *http.Request {
	username, _, _ := r.BasicAuth()

	// TODO(weiteng.huang): User creation/retrieval will be from User service
	ctx := context.WithValue(r.Context(), CTX_USER, &rohr.User{
		Id:           rohr.UserId(username),
		Organization: rohr.Organization("concur"),
		Teams:        []rohr.Team{rohr.Team(username)},
	})
	return r.WithContext(ctx)
}

func getUser(r *http.Request) (*rohr.User, error) {
	user := r.Context().Value(CTX_USER).(*rohr.User)
	if user == nil {
		return nil, fmt.Errorf("User information is missing.")
	}
	return user, nil
}

func bindAuthorization(resource rohr.Authorizable, r *http.Request) (rohr.Authorizable, error) {
	user, err := getUser(r)
	if err != nil {
		return nil, err
	}

	auth := rohr.Authorization{
		Owner: user.Id,
		GroupAccess: map[rohr.Group]rohr.PolicyMode{
			rohr.Group(user.Id):           rohr.POLICY_ALL,
			rohr.Group(user.Organization): rohr.POLICY_READ,
			rohr.Group("public"):          rohr.POLICY_NONE,
		},
	}
	resource.BindAuthorization(auth)
	return resource, nil
}
