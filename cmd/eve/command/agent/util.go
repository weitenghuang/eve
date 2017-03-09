package agent

import (
	"fmt"
	"github.com/concur/rohr"
	"github.com/concur/rohr/http"
)

func createVarFile(quoinVars []rohr.QuoinVar) []byte {
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

func getAgentUser() *rohr.User {
	return &rohr.User{
		Id:           rohr.UserId(rohr.AGENT_USER),
		Organization: rohr.Organization(rohr.AGENT_USER),
		Teams:        []rohr.Team{rohr.Team(rohr.AGENT_USER)},
	}
}
