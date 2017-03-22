package command

import (
	"fmt"
	"github.com/concur/eve/http"
	"github.com/concur/eve/http/httprouter"
	"github.com/concur/eve/pkg/config"
	"github.com/spf13/cobra"
)

var apiServer *http.ApiServer

var upCmd = &cobra.Command{
	Use:   "up",
	Short: "To start eve api server",
	Long:  `To start eve api server`,
	Run: func(cmd *cobra.Command, args []string) {
		apiServer.Router = httprouter.NewRouter()
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
