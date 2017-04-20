package agent

import (
	"runtime"

	log "github.com/Sirupsen/logrus"
	"github.com/concur/eve"
	"github.com/concur/eve/http"
	"github.com/concur/eve/pkg/terraform"
	"github.com/concur/eve/service"
	"github.com/spf13/cobra"
)

func DeleteCmd(stateServer *http.ApiServer) *cobra.Command {
	return &cobra.Command{
		Use:   "delete",
		Short: "To delete infrastructure",
		Long:  `To delete infrastructure based on user's credentials, quoin module and existing infrastructure state`,
		Run: func(cmd *cobra.Command, args []string) {
			infrastructureService := service.NewInfrastructureService(getAgentUser()) //&service.InfrastructureService{}
			if err := infrastructureService.SubscribeAsyncProc(eve.DELETE_INFRA, delete(infrastructureService, stateServer)); err != nil {
				log.Fatalln(err)
			}
			log.Printf("Listening on [%s]\n", eve.DELETE_INFRA)
			runtime.Goexit()
		},
	}
}

func delete(infraSvc eve.InfrastructureService, stateServer *http.ApiServer) eve.InfrastructureAsyncHandler {
	var toFailStatus = func(name string) {
		if statusErr := infraSvc.UpdateInfrastructureStatus(name, eve.FAILED); statusErr != nil {
			log.Println(statusErr)
		}
	}
	return func(infra *eve.Infrastructure) {
		if infra == nil {
			log.Println("Empty infrastructure object detected.")
			return
		}
		log.Printf("Start infrastructure deletion process for %s.\n", infra.Name)
		if err := infraSvc.UpdateInfrastructureStatus(infra.Name, eve.RUNNING); err != nil {
			log.Println(err)
		}
		quoinSvc := service.NewQuoinService(getAgentUser())
		id := quoinSvc.GetQuoinArchiveIdFromUri(infra.Quoin.ArchiveUri)
		quoinArchive, err := quoinSvc.GetQuoinArchive(id)
		if err != nil {
			toFailStatus(infra.Name)
			log.Println(err)
			return
		}
		if quoinArchive == nil {
			toFailStatus(infra.Name)
			log.Println("Invalid Quoin Archive Id: ", id)
			return
		}
		log.Println("Infrastructure", infra.Name, "gets Quoin Archive:", id, quoinArchive.QuoinName)
		varfile := createVarFile(infra.Variables)
		remoteState := stateEndpoint(stateServer, infra.Name)
		authenticator := createAuthenticator(infra.ProviderSlug)
		tf := terraform.NewTerraformWithAuthenticator(infra.Name, remoteState, quoinArchive.Modules, varfile, authenticator)
		if err := tf.DeleteQuoin(); err != nil {
			toFailStatus(infra.Name)
			log.Println(err)
			return
		}
		if err := infraSvc.UpdateInfrastructureStatus(infra.Name, eve.DESTROYED); err != nil {
			log.Println(err)
			return
		}
		log.Println("Deletion Done!")
	}
}
