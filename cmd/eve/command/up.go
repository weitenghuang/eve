package command

import (
	"fmt"
	"github.com/concur/rohr/http"
	"github.com/concur/rohr/http/httprouter"
	"github.com/concur/rohr/pkg/config"
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
		// Router: httprouter.NewRouter(),
	}
}
