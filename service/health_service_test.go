package service_test

import (
	"github.com/concur/rohr"
	"github.com/concur/rohr/pkg/config"
	"github.com/concur/rohr/service"
	"reflect"
	"testing"
)

var healthService *service.HealthService

func TestHealthService_New(t *testing.T) {
	healthService = service.NewHealthService()
	if !reflect.DeepEqual(healthService, &service.HealthService{}) {
		t.Errorf("Unexpected HealthService value: %#v", healthService)
	}
}

func TestHealthService_GetHealth(t *testing.T) {
	sysConfig := config.NewSystemConfig()
	mock_healthInfo := &rohr.HealthInfo{
		Hostname: sysConfig.Hostname,
		Metadata: map[string]string{
			"Version":     sysConfig.Version,
			"Environment": sysConfig.Environment,
		},
		Errors: make([]rohr.Error, 4),
		Uptime: "",
	}

	healthInfo := healthService.GetHealth()

	if len(healthInfo.Uptime) == len(mock_healthInfo.Uptime) {
		t.Errorf("HealthInfo Uptime should not be empty string")
	}
	if len(healthInfo.Errors) != len(mock_healthInfo.Errors) {
		t.Errorf("HealthInfo should report 4 errors when none of dependencies is available")
	}
	if !reflect.DeepEqual(healthInfo.Metadata, mock_healthInfo.Metadata) {
		t.Errorf("Unexpected healthInfo metadata value: %#v", healthInfo.Metadata)
	}
}
