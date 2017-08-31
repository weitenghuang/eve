package command

import (
	"fmt"
	"github.com/scipian/eve/http"
	"github.com/scipian/eve/http/httprouter"
	"github.com/scipian/eve/http/vault"
	"github.com/scipian/eve/pkg/config"
	"github.com/spf13/cobra"
)

var apiServer *http.ApiServer

var upCmd = &cobra.Command{
	Use:   "up",
	Short: "To start eve api server",
	Long:  `To start eve api server`,
	Run: func(cmd *cobra.Command, args []string) {
		apiServer.Routers = make(map[string]http.Router)
		apiServer.Routers[http.HTTPROUTER] = httprouter.NewRouter()
		apiServer.Routers[http.VAULT] = vault.NewRouter()
		apiServer.ListenAndServe()
	},
}

func init() {
	apiConfig := config.NewApiServerConfig()
	apiServer = &http.ApiServer{
		Addr:   fmt.Sprintf(":%s", apiConfig.Port),
		DNS:    apiConfig.DNS,
		Scheme: apiConfig.Scheme,
	}
	if apiConfig.Scheme == "https" {
		apiServer.CertFile = apiConfig.CertFile // Retrieve from Vault
		apiServer.KeyFile = apiConfig.KeyFile   // Retrieve from Vault
	}
}
