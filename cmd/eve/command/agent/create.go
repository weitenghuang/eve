package agent

import (
	"errors"
	"runtime"

	log "github.com/Sirupsen/logrus"
	"github.com/concur/eve"
	"github.com/concur/eve/http"
	"github.com/concur/eve/pkg/terraform"
	"github.com/concur/eve/service"
	"github.com/spf13/cobra"
)

func CreateCmd(stateServer *http.ApiServer) *cobra.Command {
	return &cobra.Command{
		Use:   "create",
		Short: "To create infrastructure",
		Long:  `To create infrastructure based on user's credentials, quoin module and existing infrastructure state`,
		Run: func(cmd *cobra.Command, args []string) {
			infrastructureService := service.NewInfrastructureService(getAgentUser()) //&service.InfrastructureService{}
			if err := infrastructureService.SubscribeAsyncProc(eve.CREATE_INFRA, create(infrastructureService, stateServer)); err != nil {
				log.Fatalln(err)
			}
			log.Printf("Listening on [%s]\n", eve.CREATE_INFRA)
			runtime.Goexit()
		},
	}
}

func create(infraSvc eve.InfrastructureService, stateServer *http.ApiServer) eve.InfrastructureAsyncHandler {
	var writeError = func(name string, infraError error) {
		if err := infraSvc.UpdateInfrastructureError(name, infraError); err != nil {
			log.Println(err)
		}
	}
	return func(infra *eve.Infrastructure) {
		if infra == nil {
			err := errors.New("Empty infrastructure object detected")
			writeError(infra.Name, err)
			log.Println(err)
			return
		}
		log.Printf("Start infrastructure creation process for %s.\n", infra.Name)
		if err := infraSvc.UpdateInfrastructureStatus(infra.Name, eve.RUNNING); err != nil {
			writeError(infra.Name, err)
			log.Println(err)
		}
		quoinSvc := service.NewQuoinService(getAgentUser())
		id := quoinSvc.GetQuoinArchiveIdFromUri(infra.Quoin.ArchiveUri)
		quoinArchive, err := quoinSvc.GetQuoinArchive(id)
		if err != nil {
			writeError(infra.Name, err)
			log.Println(err)
			return
		}
		if quoinArchive == nil {
			err := errors.New("Invalid Quoin Archive Id: " + id)
			writeError(infra.Name, err)
			log.Println(err)
			return
		}
		log.Println("Infrastructure", infra.Name, "gets Quoin Archive:", id, quoinArchive.QuoinName)
		varfile := createVarFile(infra.Variables)
		remoteState := stateEndpoint(stateServer, infra.Name)
		authenticator := createAuthenticator(infra.ProviderSlug)
		tf := terraform.NewTerraformWithAuthenticator(infra.Name, remoteState, quoinArchive.Modules, varfile, authenticator)
		if err := tf.ApplyQuoin(); err != nil {
			writeError(infra.Name, err)
			log.Println(err)
			return
		}
		if err := infraSvc.UpdateInfrastructureStatus(infra.Name, eve.DEPLOYED); err != nil {
			writeError(infra.Name, err)
			log.Println(err)
			return
		}
		writeError(infra.Name, nil)
		log.Println("Creation Done!")
	}
}
