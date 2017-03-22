package agent

import (
	"fmt"
	"github.com/concur/eve"
	"github.com/concur/eve/http"
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
