package httprouter

import (
	"encoding/json"
	"fmt"
	"github.com/concur/eve"
	eveHttp "github.com/concur/eve/http"
	"net/http"
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
	quoin.Status = eve.DEFAULT
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

func getUser(r *http.Request) (*eve.User, error) {
	user := r.Context().Value(eveHttp.CTX_USER).(*eve.User)
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
