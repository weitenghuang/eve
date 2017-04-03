package agent

import (
	"fmt"
	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/concur/eve"
	"github.com/concur/eve/client"
	"github.com/concur/eve/http"
	"github.com/concur/eve/pkg/vault"
	"github.com/concur/eve/provider/aws"
)

func createVarFile(quoinVars []eve.QuoinVar) []byte {
	var varfile []byte
	if varLen := len(quoinVars); varLen > 0 {
		for _, infraVar := range quoinVars {
			varfile = append(varfile, fmt.Sprint(infraVar.Key, "=\"", infraVar.Value, "\"\n")...)
		}
	}
	return varfile
}

func stateEndpoint(stateServer *http.ApiServer, name string) string {
	return fmt.Sprint(stateServer.GetServerAddr(), "infrastructure/", name, "/state")
}

func getAgentUser() *eve.User {
	return &eve.User{
		Id:           eve.UserId(eve.AGENT_USER),
		Organization: eve.Organization(eve.AGENT_USER),
		Teams:        []eve.Team{eve.Team(eve.AGENT_USER)},
	}
}

func getAccount(providerSlug string) *aws.Account {
	s := strings.Split(providerSlug, ":")
	if len(s) == 0 {
		// TODO(nl): return error
	}
	log.Info("providerSlug: %s", s)
	c := client.NewDefaultClient()
	p := c.GetProvider(s[0])
	a := aws.NewProvider(p)

	return a.GetAccount(s[1])
}

func getRole() string {
	key := "secret/quoin/providers/aws/meta"
	data, err := vault.GetLogicalData(key)
	if err != nil {
		// TODO(nl): return error
	}
	role := data["role"].(string)
	return role
}

func createAuthenticator(providerSlug string) *aws.Authenticator {
	a := getAccount(providerSlug)
	i := getRole()
	sessionName := "quoin"
	auth := aws.NewAuthenticator(a, i, sessionName)
	return auth
}
