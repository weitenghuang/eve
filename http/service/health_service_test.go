package service_test

import (
	"github.com/concur/rohr"
	"github.com/concur/rohr/http/service"
	"os"
	"reflect"
	"testing"
)

var healthService *service.HealthService
var mock_healthInfo *rohr.HealthInfo

func init() {
	mock_healthInfo = &rohr.HealthInfo{
		Verion:      "v0.0.1",
		Environment: os.Getenv("ENVIRONMENT"),
	}
}

func TestHealthService_New(t *testing.T) {
	healthService = service.NewHealthService()
	if !reflect.DeepEqual(healthService, &service.HealthService{
		HealthInfo: mock_healthInfo,
	}) {
		t.Errorf("Unexpected HealthService value: %#v and HealthInfo %#v", healthService, healthService.HealthInfo)
	}
}

func TestHealthService_GetHealth(t *testing.T) {
	healthInfo, err := healthService.GetHealth()
	if err != nil {
		t.Error(err)
	}
	mock_healthInfo.Uptime = ""
	if reflect.DeepEqual(healthInfo, mock_healthInfo) {
		t.Errorf("HealthInfo Uptime should not be empty string")
	}
	t.Logf("HealthInfo value: %#v", healthInfo)
}
