package agent

import (
	"errors"
	"runtime"

	log "github.com/Sirupsen/logrus"
	"github.com/scipian/eve"
	"github.com/scipian/eve/http"
	"github.com/scipian/eve/pkg/terraform"
	"github.com/scipian/eve/service"
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
		log.Printf("Start infrastructure deletion process for %s.\n", infra.Name)
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
		if err := tf.DeleteQuoin(); err != nil {
			writeError(infra.Name, err)
			log.Println(err)
			return
		}
		if err := infraSvc.UpdateInfrastructureStatus(infra.Name, eve.DESTROYED); err != nil {
			writeError(infra.Name, err)
			log.Println(err)
			return
		}
		writeError(infra.Name, nil)
		log.Println("Deletion Done!")
	}
}
